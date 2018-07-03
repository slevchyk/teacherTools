package dbase

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/slevchyk/teacherTools/models"
	"golang.org/x/crypto/bcrypt"
)

//InitDB function for first database initialization
//	creating all needed tables
//	created user with admins permissions
//			Email: admin@domain.com
//			Password: password`
func InitDB(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	var err error
	var msg string

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
				id serial PRIMARY KEY,
				email text NOT NULL UNIQUE,
				password bytea NOT NULL,
				first_name text,
				last_name text,
				type text,
				userpic text,
				deleted_at timestamp with time zone)`)

	if err != nil {
		msg = fmt.Sprintln("creating users table\n", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS sessions (
				id serial PRIMARY KEY,
				uuid text,
				user_id int references users(id),
				started_at timestamp with time zone,
				last_activity timestamp with time zone,
				ip text,
				user_agent text)`)

	if err != nil {
		msg = fmt.Sprintln("creating sessions table\n", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS levels (
				id serial PRIMARY KEY,
				name text NOT NULL UNIQUE,
				score text)`)

	if err != nil {
		msg = fmt.Sprintln("creating levels table\n", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS teachers (
				id serial PRIMARY KEY,							
				user_id int references users(id) NOT NULL,
				level_id int references levels(id) NOT NULL,
				deleted_at timestamp with time zone)`)

	if err != nil {
		msg = fmt.Sprintln("creating teachers table\n", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
				id serial PRIMARY KEY,				
				user_id int references users(id) NOT NULL,
				teacher_id int references teachers(id) NOT NULL,
				level_id int references levels(id) NOT NULL,
				deleted_at timestamp with time zone)`)

	if err != nil {
		msg = fmt.Sprintln("creating students table\n", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS parents (
				id serial PRIMARY KEY,
				user_id int references users(id),
				student_id int references students(id),
				deleted_at timestamp with time zone)`)

	if err != nil {
		msg = fmt.Sprintln("creating parents table\n", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS questions (
				id serial PRIMARY KEY,
				question text,
				type text,
				score real,
				created_at timestamp with time zone,
				deleted_at timestamp with time zone,
	  			level_id int references levels(id),
				teacher_id int references teachers(id))`)

	if err != nil {
		msg = fmt.Sprintln("creating questions table\n", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS answers (
				id serial PRIMARY KEY,
				name text,
				correct boolean,
				created_at timestamp with time zone,
				deleted_at timestamp with time zone,
				question_id int references questions(id))`)

	if err != nil {
		msg = fmt.Sprintln("creating answers table\n", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS homeworks (
				id serial PRIMARY KEY,
				score real,
				started_at timestamp with time zone,
				completed_at timestamp with time zone,
				level_id int references levels(id),
				student_id int references students(id),
	  			teacher_id int references teachers(id))`)

	if err != nil {
		msg = fmt.Sprintln("creating homeworks table\n", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS homework_specs (
				id serial PRIMARY KEY,
				answer text,				
				date timestamp with time zone,
				question_id int references questions(id),
				hometask_id int references homeworks(id),				
	  			teacher_id int references teachers(id))`)

	if err != nil {
		msg = fmt.Sprintln("creating homework_specs table\n", err.Error())
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(SelectUserByEmail(), "admin@domain.com")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
		_, err = db.Query(InsertUser(), "admin@domain.com", encryptedPassword, "Root", "User", models.UserTypeAdmin, nil)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
