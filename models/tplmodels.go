package models

//LevelRow is a part of TplLevels struct for levels.gohtml
type LevelRow struct {
	Number int
	Levels Levels
}

//TplLevels data type for levels.gohtml
type TplLevels struct {
	ID int
	Rows []LevelRow
}