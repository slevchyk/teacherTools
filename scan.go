package main

import (
	"database/sql"
	"./models"
)

func scanUser(rows *sql.Rows, u *models.Users) {

	rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Type, &u.Userpic)
}

func scanTeacher(rows *sql.Rows, t *models.Teachers, u *models.Users, l *models.Levels) {

	rows.Scan(&t.ID, &t.Active, &t.LevelID, &l.Name, &u.Email, &u.FirstName, &u.LastName, &u.Userpic)
}