package main

import (
	"database/sql"
	"./models"
)

func scanUser(rows *sql.Rows, u *models.Users) {

	rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Type, &u.Userpic)

}