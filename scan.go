package main

import (
	"database/sql"

	"github.com/slevchyk/teacherTools/models"
)

func scanUser(rows *sql.Rows, u *models.Users) error {

	err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.FirstName, &u.LastName, &u.Type, &u.Userpic)
	return err
}

func scanTeacher(rows *sql.Rows, t *models.Teachers, u *models.Users, l *models.Levels) error {

	err := rows.Scan(&t.ID, &t.UserID, &t.LevelID, &t.DeletedAt, &l.Name, &u.Email, &u.FirstName, &u.LastName, &u.Userpic)
	return err
}

func scanQuestion(rows *sql.Rows, q *models.Questions, l *models.Levels) error {

	err := rows.Scan(&q.ID, &q.Name, &q.Type, &q.Score, &q.CreatedAt, &q.LevelID, &l.Name)
	return err
}

func scanAnswer(rows *sql.Rows, a *models.Answers) error {

	err := rows.Scan(&a.ID, &a.Name, &a.Correct, &a.CreatedAt, &a.QuestionsID, &a.DeletedAt)
	return err
}

func scanSession(rows *sql.Rows, s *models.Sessions) error {

	err := rows.Scan(&s.ID, &s.UUID, &s.UserID, &s.LastActivity, &s.ID, &s.UserAgent, &s.StartedAt)
	return err
}
