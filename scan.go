package main

import (
	"./models"
	"database/sql"
)

func scanUser(rows *sql.Rows, u *models.Users) error {

	err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Type, &u.Userpic)
	return err
}

func scanTeacher(rows *sql.Rows, t *models.Teachers, u *models.Users, l *models.Levels) error {

	err := rows.Scan(&t.ID, &t.Active, &t.LevelID, &l.Name, &u.Email, &u.FirstName, &u.LastName, &u.Userpic)
	return err

}
