package main

import (
	"./dbase"
	"./models"
	"net/http"
	"time"
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

	rows, err := db.Query(dbase.GetQuery(dbase.SUserBySessionID), sessionID)
	if err != nil {
		panic(err)
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

	rows, err := db.Query(dbase.GetQuery(dbase.SUserBySessionID), sessionID)
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

	rows, err := db.Query(dbase.GetQuery(dbase.SSessions))
	if err != nil {
		panic(err)
	}

	var s models.Sessions

	for rows.Next() {
		err = rows.Scan(&s.ID, &s.UUID, &s.UserID, &s.LastActivity, &s.IP, &s.UserAgent)
		if err == nil {
			if time.Now().Sub(s.LastActivity) > (time.Duration(sessionLenght) * time.Second) {
				_, err = db.Query(dbase.GetQuery(dbase.DSessionByUUID), s.ID)
				if err != nil {
					panic(err)
				}
			}
		}
	}

	lastSessionCleaned = time.Now()
}
