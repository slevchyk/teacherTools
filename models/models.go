package models

import (
	"time"
)

type Users struct {
	ID        int
	Email     string
	Password  []byte
	FirstName string
	LastName  string
	Type      string
	Userpic   string
}

type Sessions struct {
	ID           int
	UUID         string
	UserID       int
	LastActivity time.Time
	IP           string
	UserAgent    string
}

type Levels struct {
	ID    int
	Score int8
	Name  string
}

type Questions struct {
	ID          int
	Question    string
	Type        string
	Score       float32
	DateCteated time.Time
	LevelID     int
	TeacherID   int
}

type Answers struct {
	ID          int
	Answer      string
	Correct     bool
	DateCteated time.Time
	QuestionsID int
	TeacherID   int
}

type Hometasks struct {
	ID            int
	Score         float32
	DateStarted   time.Time
	DateCompleted time.Time
	LevelID       int
	StudentID     int
	TeacherID     int
}

type HometaskSpecs struct {
	ID         int
	Answer     int
	Date       time.Time
	HometaskID int
	QuestionID int
}

type Teachers struct {
	ID      int
	LevelID int
	UserID  int
}

type Students struct {
	ID      int
	LevelID int
	UserID  int
}

type Parents struct {
	ID        int
	UserID    int
	StudentID int
}
