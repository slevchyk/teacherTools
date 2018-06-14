package models

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

type QuestionsRow struct {
	Number int
	Question Questions
	Level Levels
}

type QuestionsColumnsVisibility struct {
	Level bool `json:"level"`
	QType bool `json:"type"`
	Score bool `json:"score"`
	DateCreated bool `json:"dateCreated"`
}

type TplQuestions struct {
	ColumnsVisibility QuestionsColumnsVisibility
	Rows              []QuestionsRow
}