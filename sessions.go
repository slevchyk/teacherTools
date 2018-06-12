package main

import (
	"net/http"
	"time"
	"./models"
	"./dbase"
)

func getUser(w http.ResponseWriter, r *http.Request) models.Users {

	var u models.Users

	c, err := r.Cookie("session")
	if err != nil {
		return u
	}
	c.MaxAge = SessionLenght
	http.SetCookie(w, c)

	sessionID := c.Value

	rows, err := DB.Query(dbase.GetQuery(dbase.S_UserBySessionID), sessionID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		scanUser(rows, &u)
	}

	return u
}

func alreadyLoggedIn(w http.ResponseWriter, r *http.Request) bool {

	if time.Now().Sub(lastSessionCleaned) > (time.Duration(SessionLenght) * time.Second) {
		go cleanSession()
	}

	c, err := r.Cookie("session")
	if err != nil {
		return false
	}

	sessionID := c.Value

	rows, err := DB.Query(dbase.GetQuery(dbase.S_UserBySessionID), sessionID)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	ok := false

	if rows.Next() {
		ok = true
	}

	c.MaxAge = SessionLenght
	http.SetCookie(w, c)

	return ok
}

func cleanSession() {

	rows, err := DB.Query(dbase.GetQuery(dbase.S_Sessions))
	if err != nil {
		panic(err)
	}

	var s models.Sessions

	for rows.Next()	{
		rows.Scan(&s.ID, &s.UUID, &s.UserID, &s.LastActivity, &s.IP, &s.UserAgent)

		if time.Now().Sub(s.LastActivity) > (time.Duration(SessionLenght) * time.Second) {
			_, err = DB.Query(dbase.GetQuery(dbase.D_SessionByUUID), s.ID)
			if err != nil {
				panic(err)
			}
		}
	}

	lastSessionCleaned = time.Now()
}


