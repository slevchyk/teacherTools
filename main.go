package main

import (
	"./cfg"
	"./models"
	"./dbase"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
	"strconv"
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

var td tplData
var lastSessionCleaned time.Time
var userTypes = map[string]string{
	"a": "admin",
	"t": "teacher",
	"s": "student",
	"p": "parents",
}

const SessionLenght int = 300

func init() {
	lastSessionCleaned = time.Now()
}

func main() {

	defer cfg.DB.Close()

	http.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/user", userHandler)
	http.HandleFunc("/admin/db", adminDbHandler)
	http.HandleFunc("/admin/sessions", adminSessionsHandler)
	http.HandleFunc("/admin/editteacher", editteacherHandler)
	http.ListenAndServe(":8080", nil)

}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	u := getUser(w, r)

	if u.Type == userTypes["t"] {

		teacherHandler(w, r, u)
	}

	cfg.Tpl.ExecuteTemplate(w, "index.gohtml", td)

}

func teacherHandler(w http.ResponseWriter, r *http.Request, u models.Users) {

	td := struct {
		U models.Users
	}{
		u,
	}

	cfg.Tpl.ExecuteTemplate(w, "teacher.gohtml", td)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	if alreadyLoggedIn(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	if r.Method == http.MethodPost {

		email := r.FormValue("email")
		password := r.FormValue("password")

		//check user
		rows, err := cfg.DB.Query(dbase.GetQuery(dbase.S_UserByEmail), email)
		if err != nil {
			panic(err)
		}

		var u models.Users

		if rows.Next() {
			rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Type)
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
			MaxAge: SessionLenght,
		}
		http.SetCookie(w, c)

		var s models.Sessions

		s.UUID = sessionID.String()
		s.UserID = u.ID
		s.LastActivity = time.Now()
		s.IP = r.RemoteAddr
		s.UserAgent = r.Header.Get("User-Agent")

		_, err = cfg.DB.Query(dbase.GetQuery(dbase.I_Session), s.UUID, u.ID, s.LastActivity, s.IP, s.UserAgent)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	cfg.Tpl.ExecuteTemplate(w, "login.gohtml", nil)
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
	_, err = cfg.DB.Query(dbase.GetQuery(dbase.D_SessionByUUID), sessionID)
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

		userName := r.FormValue("email")
		password := r.FormValue("password")
		firstName := r.FormValue("firstName")
		lastName := r.FormValue("lastName")

		rows, err := cfg.DB.Query(dbase.GetQuery(dbase.S_UserByEmail), userName)
		if err != nil {
			panic(err)
		}

		if rows.Next() {
			http.Error(w, "Usarname already taken", http.StatusForbidden)
			return
		}

		sessionID, _ := uuid.NewV4()

		c := &http.Cookie{
			Name:  "session",
			Value: sessionID.String(),
		}
		http.SetCookie(w, c)

		encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Can't encrypt password", http.StatusInternalServerError)
			return
		}

		_, err = cfg.DB.Query(dbase.GetQuery(dbase.I_User), userName, encryptedPassword, firstName, lastName, false)
		if err != nil {
			panic(err)
		}

		rows, err = cfg.DB.Query(dbase.GetQuery(dbase.S_UserByEmail), userName)
		if err != nil {
			panic(err)
		}

		var userID int
		if rows.Next() {
			rows.Scan(&userID)
		}

		if userID != 0 {
			cfg.DB.Query(dbase.GetQuery(dbase.I_Session), sessionID.String(), userID, time.Now(), r.Header.Get("X-Forwarded-For"))
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cfg.Tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func adminDbHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		action := r.FormValue("action")

		switch action {
		case "init":
			initDb(w, r)
		}
	}

	cfg.Tpl.ExecuteTemplate(w, "admin.gohtml", td)
}

func initDb(w http.ResponseWriter, r *http.Request) {

	clearTplData()

	var err error

	_, err = cfg.DB.Exec(`CREATE TABLE IF NOT EXISTS users (
				id serial PRIMARY KEY,
				email text NOT NULL UNIQUE,
				password bytea NOT NULL,
				firstName text,
				lastName text,
				type char(100))`)

	if err != nil {
		initDbErrors(w, r, "create users table", err.Error())
		return
	}

	_, err = cfg.DB.Exec(`CREATE TABLE IF NOT EXISTS sessions (
				id serial PRIMARY KEY,
				uuid text,
				userId int references users(ID),
				lastActivity timestamp with time zone,
				ip text,
				userAgent text)`)

	if err != nil {
		initDbErrors(w, r, "create sessions table", err.Error())
		return
	}

	_, err = cfg.DB.Exec(`CREATE TABLE IF NOT EXISTS levels (
				id serial PRIMARY KEY,
				name text NOT NULL UNIQUE,
				score text)`)

	if err != nil {
		initDbErrors(w, r, "create levels table", err.Error())
		return
	}

	_, err = cfg.DB.Exec(`CREATE TABLE IF NOT EXISTS teachers (
				id serial PRIMARY KEY,
				active boolean DAFAULT TRUE NOT NULL,
				userId int references users(ID) NOT NULL,
				levelId int references levels(ID) NOT NULL)`)

	if err != nil {
		initDbErrors(w, r, "create teachers table", err.Error())
		return
	}
	_, err = cfg.DB.Exec(`CREATE TABLE IF NOT EXISTS students (
				id serial PRIMARY KEY,
				active boolean DAFAULT TRUE NOT NULL,
				userId int references users(ID) NOT NULL,
				teacherId int references teachers(ID) NOT NULL,
				levelId int references levels(ID) NOT NULL)`)

	if err != nil {
		initDbErrors(w, r, "create students table", err.Error())
		return
	}

	_, err = cfg.DB.Exec(`CREATE TABLE IF NOT EXISTS parents (
				id serial PRIMARY KEY,
				userId int references users(ID),
				studentId int references students(ID))`)

	if err != nil {
		initDbErrors(w, r, "create parents table", err.Error())
		return
	}
	_, err = cfg.DB.Exec(`CREATE TABLE IF NOT EXISTS questions (
				id serial PRIMARY KEY,
				question text,
				type text,
				score real,
				datecreated timestamp,
	  			levelId int references levels(ID),
				teacherId int references teachers(ID))`)

	if err != nil {
		initDbErrors(w, r, "create questions table", err.Error())
		return
	}

	_, err = cfg.DB.Exec(`CREATE TABLE IF NOT EXISTS answers (
				id serial PRIMARY KEY,
				answer text,
				correct boolean,
				datecreated timestamp,
				questionId int references questions(ID),
	  			teacherId int references teachers(ID))`)

	if err != nil {
		initDbErrors(w, r, "create answers table", err.Error())
		return
	}

	_, err = cfg.DB.Exec(`CREATE TABLE IF NOT EXISTS hometasks (
				id serial PRIMARY KEY,
				score real,
				dateStarted text,
				dateCompleted timestamp,
				levelId int references levels(ID),
				studentId int references students(ID),
	  			teacherId int references teachers(ID))`)

	if err != nil {
		initDbErrors(w, r, "create hometasks table", err.Error())
		return
	}

	_, err = cfg.DB.Exec(`CREATE TABLE IF NOT EXISTS hometaskspecs (
				id serial PRIMARY KEY,
				answer text,				
				date timestamp,
				questionId int references questions(ID),
				hometaskId int references hometasks(ID),				
	  			teacherId int references teachers(ID))`)

	if err != nil {
		initDbErrors(w, r, "create hometaskspecs table", err.Error())
		return
	}

	rows, err := cfg.DB.Query(dbase.GetQuery(dbase.S_UserByEmail), "admin")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if !rows.Next() {
		encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
		_, err = cfg.DB.Exec(dbase.GetQuery(dbase.I_User), "admin@domain.com", encryptedPassword, "Root", "User", "admin")
		fmt.Println(err)
	}

	td.SysMsg = `Init DB completed.
				Created user with admins permissions
				Email: admin
				Password: password`

	http.Redirect(w, r, "/admin/db", http.StatusSeeOther)
}

func initDbErrors(w http.ResponseWriter, r *http.Request, key string, err string) {

	te := tplErr{key, err}

	td.Err = nil
	td.Err = append(td.Err, te)

	http.Redirect(w, r, "/admin/db", http.StatusSeeOther)
}

func clearTplData() {

	td.SysMsg = ""
	td.Err = nil
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

		_, err = cfg.DB.Query(dbase.GetQuery(dbase.D_SessionByID), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		http.Redirect(w, r, "/admin/sessions", http.StatusSeeOther)
	}

	rows, err := cfg.DB.Query(dbase.GetQuery(dbase.S_Sessions))
	if err != nil {
		panic(err)
	}

	var s models.Sessions

	type row struct {
		Number  int
		Session models.Sessions
	}

	var ts []row
	var i int

	for rows.Next() {

		rows.Scan(&s.ID, &s.UUID, &s.UserID, &s.LastActivity, &s.IP, &s.UserAgent)
		i++
		ts = append(ts, row{i, s})
	}

	cfg.Tpl.ExecuteTemplate(w, "sessions.gohtml", ts)
}

func editteacherHandler(w http.ResponseWriter, r *http.Request) {

	var sl = make(map[int]string)
	var l models.Levels
	var u models.Users
	var t models.Teachers

	rows, err := cfg.DB.Query(dbase.GetQuery(dbase.S_Levels))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&l.ID, &l.Name, &l.Score)
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

		rows, err := cfg.DB.Query(dbase.GetQuery(dbase.S_UserByEmail), u.Email)
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

		_, err = cfg.DB.Query(dbase.GetQuery(dbase.I_User), u.Email, u.Password, u.FirstName, u.LastName, "teacher")
		if err != nil {
			http.Error(w, "Can't create user", http.StatusInternalServerError)
			return
		}

		rows, err = cfg.DB.Query(dbase.GetQuery(dbase.S_UserByEmail), u.Email)
		if err != nil {
			http.Error(w, "Can't select user", http.StatusInternalServerError)
		}

		if rows.Next() {
			rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Type)
		}

		if u.ID != 0 {

			t.UserID = u.ID
			t.LevelID = l.ID

			_, err = cfg.DB.Query(dbase.GetQuery(dbase.I_Teacher), t.UserID, t.LevelID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	cfg.Tpl.ExecuteTemplate(w, "editteacher.gohtml", sl)
}

func userHandler(w http.ResponseWriter, r *http.Request)  {

	type tplData struct{
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

		rows, err := cfg.DB.Query(dbase.GetQuery(dbase.S_UserByID), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if rows.Next() {
			rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Type)
		}

		td.View = true
		td.User = u
	}

	cfg.Tpl.ExecuteTemplate(w, "user.gohtml", td)
}