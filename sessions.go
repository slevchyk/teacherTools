package main

import (
	"net/http"
	"time"

	"github.com/slevchyk/teacherTools/dbase"
	"github.com/slevchyk/teacherTools/models"
)

func getUser(w http.ResponseWriter, r *http.Request) models.Users {

	var u models.Users

	c, err := r.Cookie("session")
	if err != nil {
		return u
	}
	c.MaxAge = sessionLenght
	http.SetCookie(w, c)

	sessionID := c.Value

	_, err = db.Query(dbase.UpdateSessionLastActivityByUuid(), sessionID, time.Now().UTC())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	rows, err := db.Query(dbase.SelectUserBySessionID(), sessionID)
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

	return u
}

func alreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {

	if time.Now().Sub(lastSessionCleaned) > (time.Duration(sessionLenght) * time.Second) {
		go cleanSession()
	}

	c, err := r.Cookie("session")
	if err != nil {
		return false
	}

	sessionID := c.Value

	_, err = db.Query(dbase.UpdateSessionLastActivityByUuid(), sessionID, time.Now().UTC())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	rows, err := db.Query(dbase.SelectUserBySessionID(), sessionID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	ok := false

	if rows.Next() {
		ok = true
	}

	c.MaxAge = sessionLenght
	http.SetCookie(w, c)

	return ok
}

func cleanSession() {

	rows, err := db.Query(dbase.SelectSessions())
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var s models.Sessions

	for rows.Next() {
		err = rows.Scan(&s.ID, &s.UUID, &s.UserID, &s.LastActivity, &s.IP, &s.UserAgent)
		if err == nil {
			if time.Now().Sub(s.LastActivity) > (time.Duration(sessionLenght) * time.Second) {
				_, err = db.Query(dbase.DeleteSessionByUUID(), s.ID)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	lastSessionCleaned = time.Now()
}
