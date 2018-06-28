package main

import (
	"./dbase"
	"./models"
	"crypto/sha1"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"encoding/json"
	"encoding/base64"
)

type tplErr struct {
	Name  string
	Value string
}

type tplData struct {
	Err    []tplErr
	SysMsg string
	User   models.Users
}

var db *sql.DB
var tpl *template.Template

var td tplData
var lastSessionCleaned time.Time

var questionTypes = map[string]string{
	"qTypeChooseCorrectTranslation": "Choose the correct translation",
	"qTypeSelectMissingWord": "Select the missing word",
}

const sessionLenght = 300

func init() {
	var err error

	db, err = sql.Open("postgres", "postgres://postgres:sql@localhost/coach?sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))

	lastSessionCleaned = time.Now()
}

func main() {

	defer db.Close()

	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/levels", levelsHandler)
	http.HandleFunc("/teachers", teachersHandler)
	http.HandleFunc("/teacher", teacherHandler)
	http.HandleFunc("/questions", questionsHandler)
	http.HandleFunc("/question", questionHandler)
	http.HandleFunc("/answers", answersHandler)
	http.HandleFunc("/admin/db", adminDbHandler)
	http.HandleFunc("/admin/sessions", adminSessionsHandler)
	http.HandleFunc("/admin/editteacher", editteacherHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	err := tpl.ExecuteTemplate(w, "index.gohtml", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if r.Method == http.MethodPost {

		email := r.FormValue("email")
		password := r.FormValue("password")

		//check user
		rows, err := db.Query(dbase.GetQuery(dbase.SUserByEmail), email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer rows.Close()

		var u models.Users

		if rows.Next() {
			err = scanUser(rows, &u)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Usarname do not patch", http.StatusForbidden)
			return
		}

		//check password
		err = bcrypt.CompareHashAndPassword(u.Password, []byte(password))
		if err != nil {
			http.Error(w, "Password do not match", http.StatusForbidden)
			return
		}

		//create session
		sessionID, _ := uuid.NewV4()

		c := &http.Cookie{
			Name:   "session",
			Value:  sessionID.String(),
			MaxAge: sessionLenght,
		}
		http.SetCookie(w, c)

		var s models.Sessions

		s.UUID = sessionID.String()
		s.UserID = u.ID
		s.LastActivity = time.Now()
		s.IP = r.RemoteAddr
		s.UserAgent = r.Header.Get("User-Agent")

		_, err = db.Query(dbase.GetQuery(dbase.ISession), s.UUID, u.ID, s.LastActivity, s.IP, s.UserAgent)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	err := tpl.ExecuteTemplate(w, "login.gohtml", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	if !alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	c, err := r.Cookie("session")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	sessionID := c.Value
	_, err = db.Query(dbase.GetQuery(dbase.DSessionByUUID), sessionID)
	if err != nil {
		panic(err)
	}

	c.MaxAge = -1
	http.SetCookie(w, c)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {

	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if r.Method == http.MethodPost {

		var u models.Users

		u.Email = r.FormValue("email")

		rows, err := db.Query(dbase.GetQuery(dbase.SUserByEmail), u.Email)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		if rows.Next() {
			http.Error(w, "Usarname already taken", http.StatusForbidden)
			return
		}

		password := r.FormValue("password")
		sessionID, _ := uuid.NewV4()

		c := &http.Cookie{
			Name:  "session",
			Value: sessionID.String(),
		}
		http.SetCookie(w, c)

		u.Password, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Can't encrypt password", http.StatusInternalServerError)
			return
		}

		u.FirstName = r.FormValue("firstName")
		u.LastName = r.FormValue("lastName")

		mf, fh, err := r.FormFile("userpic")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		ext := strings.Split(fh.Filename, ".")[1]
		h := sha1.New()

		_, err = io.Copy(h, mf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		u.Userpic = fmt.Sprintf("%x", h.Sum(nil)) + "." + ext

		wd, err := os.Getwd()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		path := filepath.Join(wd, "public", "userpics", u.Userpic)
		nf, err := os.Create(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer nf.Close()

		_, err = mf.Seek(0, 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = io.Copy(nf, mf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = db.Query(dbase.GetQuery(dbase.IUser), u.Email, u.Password, u.FirstName, u.LastName, u.Type, u.Userpic)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		rows, err = db.Query(dbase.GetQuery(dbase.SUserByEmail), u.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer rows.Close()

		if rows.Next() {
			err = scanUser(rows, &u)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		if u.ID != 0 {
			_, err = db.Query(dbase.GetQuery(dbase.ISession), sessionID.String(), u.ID, time.Now(), r.Header.Get("X-Forwarded-For"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := tpl.ExecuteTemplate(w, "signup.gohtml", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func adminDbHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		action := r.FormValue("action")

		switch action {
		case "init":
			dbase.InitDB(db, w, r)
		}
	}

	err := tpl.ExecuteTemplate(w, "admin.gohtml", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func adminSessionsHandler(w http.ResponseWriter, r *http.Request) {

	var action string
	action = r.FormValue("action")

	if action == "delete" {
		var id int

		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = db.Query(dbase.GetQuery(dbase.DSessionByID), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/admin/sessions", http.StatusSeeOther)
	}

	rows, err := db.Query(dbase.GetQuery(dbase.SSessions))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var s models.Sessions

	type row struct {
		Number  int
		Session models.Sessions
	}

	var ts []row
	var i int

	for rows.Next() {

		err = rows.Scan(&s.ID, &s.UUID, &s.UserID, &s.LastActivity, &s.IP, &s.UserAgent)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		i++
		ts = append(ts, row{Number: i, Session: s})
	}

	err = tpl.ExecuteTemplate(w, "sessions.gohtml", ts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func editteacherHandler(w http.ResponseWriter, r *http.Request) {

	var sl = make(map[int]string)
	var l models.Levels
	var u models.Users
	var t models.Teachers

	rows, err := db.Query(dbase.GetQuery(dbase.SLevels))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&l.ID, &l.Name, &l.Score)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		sl[l.ID] = l.Name
	}

	if r.Method == http.MethodPost {

		l.ID, err = strconv.Atoi(r.FormValue("level"))
		if err != nil {
			panic(err)
		}

		u.Email = r.FormValue("email")
		u.FirstName = r.FormValue("firstName")
		u.LastName = r.FormValue("lastName")

		rows, err := db.Query(dbase.GetQuery(dbase.SUserByEmail), u.Email)
		if err != nil {
			panic(err)
		}
		defer rows.Close()

		if rows.Next() {
			http.Error(w, "Usarname already taken", http.StatusInternalServerError)
			return
		}

		password := r.FormValue("password")
		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Can't encrypt password", http.StatusInternalServerError)
			return
		}
		u.Password = encryptedPassword

		_, err = db.Query(dbase.GetQuery(dbase.IUser), u.Email, u.Password, u.FirstName, u.LastName, "teacher")
		if err != nil {
			http.Error(w, "Can't create user", http.StatusInternalServerError)
			return
		}

		rows, err = db.Query(dbase.GetQuery(dbase.SUserByEmail), u.Email)
		if err != nil {
			http.Error(w, "Can't select user", http.StatusInternalServerError)
		}
		defer rows.Close()

		if rows.Next() {
			err = scanUser(rows, &u)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		if u.ID != 0 {

			t.UserID = u.ID
			t.LevelID = l.ID

			_, err = db.Query(dbase.GetQuery(dbase.ITeacher), t.UserID, t.LevelID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	err = tpl.ExecuteTemplate(w, "editteacher.gohtml", sl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {

	type tplData struct {
		View bool
		User models.Users
	}

	var action string
	var td tplData

	td.View = false

	action = r.FormValue("action")
	if action == "view" {

		var u models.Users

		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		rows, err := db.Query(dbase.GetQuery(dbase.SUserByID), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer rows.Close()

		if rows.Next() {
			err = scanUser(rows, &u)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		td.View = true
		td.User = u
	}

	err := tpl.ExecuteTemplate(w, "user.gohtml", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {

	type tplData struct {
		Number  int
		Teacher models.Teachers
		User    models.Users
		Level   models.Levels
	}

	var std []tplData
	var i int

	var t models.Teachers
	var u models.Users
	var l models.Levels

	rows, err := db.Query(dbase.GetQuery(dbase.STeachers))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		err = scanTeacher(rows, &t, &u, &l)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		i++
		std = append(std, tplData{Number: i, Teacher: t, User: u, Level: l})
	}

	err = tpl.ExecuteTemplate(w, "teacherslist.gohtml", std)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func teacherHandler(w http.ResponseWriter, r *http.Request) {

	type tplData struct {
		View    bool
		Teacher models.Teachers
		User    models.Users
		Level   models.Levels
		Levels  []models.Levels
	}

	var td tplData
	var l models.Levels
	var action string

	td.View = false

	rows, err := db.Query(dbase.GetQuery(dbase.SLevels))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&l.ID, &l.Name, &l.Score)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		td.Levels = append(td.Levels, l)
	}

	action = r.FormValue("action")
	if action == "view" {

		var t models.Teachers
		var u models.Users

		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		rows, err = db.Query(dbase.GetQuery(dbase.STeacherByID), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer rows.Close()

		if rows.Next() {
			err = scanTeacher(rows, &t, &u, &l)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		td.View = true
		td.Teacher = t
		td.User = u
		td.Level = l
	}

	err = tpl.ExecuteTemplate(w, "teacher.gohtml", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func levelsHandler(w http.ResponseWriter, r *http.Request) {

	var td models.TplLevels
	var sr []models.LevelRow
	var l models.Levels
	var i int

	rows, err := db.Query(dbase.GetQuery(dbase.SLevels))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&l.ID, &l.Name, &l.Score)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		i++
		sr = append(sr, models.LevelRow{Number: i, Levels: l})
	}
	td.Rows = sr

	var action string
	action = r.FormValue("action")

	switch action {
	case "add":
		l.Score, err = strconv.Atoi(r.FormValue("score"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		l.Name = r.FormValue("name")

		_, err = db.Query(dbase.GetQuery(dbase.ILevel), l.Name, l.Score)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/levels", http.StatusSeeOther)
		return

	case "edit":

		l.ID, err = strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		td.ID = l.ID
	case "update":

		l.ID, err = strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		l.Score, err = strconv.Atoi(r.FormValue("score"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		l.Name = r.FormValue("name")

		_, err := db.Query(dbase.GetQuery(dbase.ULevel), l.ID, l.Name, l.Score)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/levels", http.StatusSeeOther)
		return
	}

	err = tpl.ExecuteTemplate(w, "levels.gohtml", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func questionsHandler(w http.ResponseWriter, r *http.Request)  {

	var td models.TplQuestions
	var sr []models.QuestionsRow
	var q models.Questions
	var l models.Levels
	var i int

	var cvLevel string
	var cvType string
	var cvScore string
	var cvDateCreated string

	c, err := r.Cookie("questionColumnsVisibility")
	if err != nil {
		td.ColumnsVisibility.Level = true
		td.ColumnsVisibility.QType = true
		td.ColumnsVisibility.Score = true
		td.ColumnsVisibility.DateCreated = false
	} else {

		dsb, err := base64.StdEncoding.DecodeString(c.Value)
		if err != nil {
			c.MaxAge = -1
			http.SetCookie(w, c)

			http.Redirect(w, r, "/questions", http.StatusSeeOther)
		}

		err = json.Unmarshal(dsb, &td.ColumnsVisibility)
		if err != nil {
			c.MaxAge = -1
			http.SetCookie(w, c)

			http.Redirect(w, r, "/questions", http.StatusSeeOther)		}
	}

	if r.Method == http.MethodPost {

		cvLevel = r.FormValue("cvLevel")
		if cvLevel == "1" {
			td.ColumnsVisibility.Level = true
		} else {
			td.ColumnsVisibility.Level = false
		}

		cvType = r.FormValue("cvType")
		if cvType == "1" {
			td.ColumnsVisibility.QType = true
		} else {
			td.ColumnsVisibility.QType = false
		}

		cvScore = r.FormValue("cvScore")
		if cvScore == "1" {
			td.ColumnsVisibility.Score = true
		} else {
			td.ColumnsVisibility.Score = false
		}

		cvDateCreated = r.FormValue("cvDateCreated")
		if cvDateCreated == "1" {
			td.ColumnsVisibility.DateCreated = true
		} else {
			td.ColumnsVisibility.DateCreated = false
		}

		fmt.Println(td.ColumnsVisibility)

		bs, err := json.Marshal(td.ColumnsVisibility)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		s64 := base64.StdEncoding.EncodeToString(bs)

		c := &http.Cookie{
			Name:   "questionColumnsVisibility",
			Value:  s64,
		}
		http.SetCookie(w, c)

		http.Redirect(w, r, "/questions", http.StatusSeeOther)

	}

	rows, err := db.Query(dbase.GetQuery(dbase.SelectQuestions))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		err = scanQuestion(rows, &q, &l)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		//change question type key to alias
		q.Type = questionTypes[q.Type]

		i++
		sr = append(sr, models.QuestionsRow{Number: i, Question: q, Level: l})
	}

	td.Rows = sr



	err = tpl.ExecuteTemplate(w, "questions.gohtml", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func questionHandler(w http.ResponseWriter, r *http.Request) {

	var td models.TplQuestion
	var action string
	var q models.Questions
	var l models.Levels
	var a models.Answers
	var i int

	td.Edit = false
	td.QuestionTypes = questionTypes

	rows, err := db.Query(dbase.GetQuery(dbase.SLevels))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&l.ID, &l.Name, &l.Score)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		td.Levels = append(td.Levels, l)
	}

	action = r.FormValue("action")
	if action == "edit" {

		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		rows, err := db.Query(dbase.GetQuery(dbase.SelectQuestionByID), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		defer rows.Close()

		if rows.Next() {
			err = scanQuestion(rows, &q, &l)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		td.Edit = true
		td.Question = q

		rows, err = db.Query(dbase.GetQuery(dbase.SelectAnswersByQuestionID), q.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		for rows.Next() {
			err = scanAnswers(rows, &a)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			i++
			td.AnswerRows = append(td.AnswerRows, models.AnswerRow{Number: i, Answer: a})
		}
	}

	err = tpl.ExecuteTemplate(w, "question.gohtml", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func answersHandler(w http.ResponseWriter, r *http.Request)  {

	var td models.TplAnswers
	var a models.Answers
	var q models.Questions
	var l models.Levels
	var action string
	var err error
	var i int

	q.ID, err = strconv.Atoi(r.FormValue("qid"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	action = r.FormValue("action")

	switch action {
	case "edit":

		td.AnswerID, err = strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "update":

		a.ID, err = strconv.Atoi(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		a.Name = r.FormValue("name")

		correct := r.FormValue("correct")
		if correct == "true" {
			a.Correct = true
		} else {
			a.Correct = false
		}

		_, err := db.Query(dbase.GetQuery(dbase.UpdateAnswer), a.ID, a.Name, a.Correct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		url := "answers?qid=" + strconv.Itoa(q.ID)
		http.Redirect(w, r, url, http.StatusSeeOther)

	case "add":

		a.Name = r.FormValue("name")

		correct := r.FormValue("correct")
		if correct == "true" {
			a.Correct = true
		} else {
			a.Correct = false
		}

		a.DateCreated = time.Now()

		_, err = db.Query(dbase.GetQuery(dbase.InsertAnswer), a.Name, a.Correct, a.DateCreated, q.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		url := "answers?qid=" + strconv.Itoa(q.ID)
		http.Redirect(w, r, url, http.StatusSeeOther)
	}

	rows, err := db.Query(dbase.GetQuery(dbase.SelectQuestionByID), q.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	if rows.Next() {
		err = scanQuestion(rows, &q, &l)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	td.Question = q
	td.Level = l


	rows, err = db.Query(dbase.GetQuery(dbase.SelectAnswersByQuestionID), q.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()


	for rows.Next() {
		err = scanAnswers(rows, &a)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		i++
		td.Rows = append(td.Rows, models.AnswerRow{i, a})
	}

	err = tpl.ExecuteTemplate(w, "answers.gohtml", td)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}