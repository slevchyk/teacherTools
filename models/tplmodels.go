package models

type TplNavBar struct {
	LoggedIn bool
	User     Users
	MainMenu []struct {
		Link  string
		Alias string
	}
}

type TplError struct {
	Error   bool
	Title   string
	Message string
}

//LevelRow is a part of TplLevels struct for levels.gohtml
type LevelRow struct {
	Number int
	Levels Levels
}

//TplLevels data type for levels.gohtml
type TplLevels struct {
	ID   int
	Rows []LevelRow
}

//QuestionsRow is a part of TplQuestions struct for questions.gohtml
type QuestionsRow struct {
	Number   int
	Question Questions
	Level    Levels
}

//questionsColumnsVisibility is a part of TplQuestions struct for questions.gohtml
//set columns visibility on questions web-page
type questionsColumnsVisibility struct {
	Level       bool `json:"level"`
	QType       bool `json:"type"`
	Score       bool `json:"score"`
	DateCreated bool `json:"dateCreated"`
}

//TplQuestions data type for questions.gohtml
type TplQuestions struct {
	ColumnsVisibility questionsColumnsVisibility
	Rows              []QuestionsRow
}

type TplQuestion struct {
	Edit          bool
	Question      Questions
	Level         Levels
	Levels        []Levels
	AnswerRows    []AnswerRow
	QuestionTypes map[string]string
}

type AnswerRow struct {
	Number  int
	Deleted bool
	Answer  Answers
}

type TplAnswers struct {
	Question   Questions
	Level      Levels
	AnswerID   int
	AnswerRows []AnswerRow
}

type TeachersRow struct {
	Number  int
	Deleted bool
	Teacher Teachers
	User    Users
	Level   Levels
}

type TplTeachers struct {
	NavBar TplNavBar
	Error  TplError
	Rows   []TeachersRow
}

type TplTeacher struct {
	NavBar  TplNavBar
	Error   TplError
	Edit    bool
	Deleted bool
	Teacher Teachers
	User    Users
	Level   Levels
	Levels  []Levels
}

type TplIndex struct {
	NavBar TplNavBar
	Error  TplError
}

type TplAdmin struct {
	NavBar TplNavBar
	Error  TplError
}

type TplSessionRow struct {
	Number  int
	Session Sessions
}

type TplSessions struct {
	NavBar TplNavBar
	Error  TplError
	Rows   []TplSessionRow
}

type TplUser struct {
	NavBar TplNavBar
	Error  TplError
	View   bool
	User   Users
}
