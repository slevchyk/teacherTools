package models

import (
	"time"
)

//Users is a struct of table users
type Users struct {
	ID        int
	Email     string
	Password  []byte
	FirstName string
	LastName  string
	Type      string
	Userpic   string
}

//Sessions is a struct of table sessions
type Sessions struct {
	ID           int
	UUID         string
	UserID       int
	LastActivity time.Time
	IP           string
	UserAgent    string
}

//Levels is a struct of table levels
type Levels struct {
	ID    int
	Name  string
	Score int
}

//Questions is a struct of table questions
type Questions struct {
	ID          int
	Question    string
	Type        string
	Score       float32
	DateCteated time.Time
	LevelID     int
	TeacherID   int
}

//Answers is a struct of table answers
type Answers struct {
	ID          int
	Answer      string
	Correct     bool
	DateCteated time.Time
	QuestionsID int
	TeacherID   int
}

//Hometasks is a struct of table hometasks
type Hometasks struct {
	ID            int
	Score         float32
	DateStarted   time.Time
	DateCompleted time.Time
	LevelID       int
	StudentID     int
	TeacherID     int
}

//HometaskSpecs is a struct of table hometaskSpecs
type HometaskSpecs struct {
	ID         int
	Answer     int
	Date       time.Time
	HometaskID int
	QuestionID int
}

//Teachers is a struct of table teachers
type Teachers struct {
	ID      int
	LevelID int
	UserID  int
	Active  bool
}

//Students is a struct of table students
type Students struct {
	ID      int
	LevelID int
	UserID  int
}

//Parents is a struct of table parents
type Parents struct {
	ID        int
	UserID    int
	StudentID int
}
