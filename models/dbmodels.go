package models

import (
	"time"

	"github.com/lib/pq"
)

const (
	UserTypeAdmin   = "admin"
	UserTypeTeacher = "teacher"
	UserTypeStudent = "student"
	UserTypeParent  = "parent"
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
	DeletedAt pq.NullTime
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

//Level is a struct of table levels
type Levels struct {
	ID    int
	Name  string
	Score int
}

//Questions is a struct of table questions
type Questions struct {
	ID        int
	Name      string
	Type      string
	Score     float32
	CreatedAt time.Time
	DeletedAt time.Time
	LevelID   int
	TeacherID int
}

//AnswerRows is a struct of table answers
type Answers struct {
	ID          int
	Name        string
	Correct     bool
	CreatedAt   time.Time
	QuestionsID int
	DeletedAt   pq.NullTime
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
	ID        int
	LevelID   int
	UserID    int
	DeletedAt pq.NullTime
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
