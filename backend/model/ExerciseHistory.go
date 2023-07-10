package model

type ExerciseHistory struct {
	ID           int    `gorm:"primaryKey;autoIncrement;column:ID"`
	ExerciseId   string `gorm:"column:exercise_id"`
	Username     string `gorm:"column:username"`
	Origin       string `gorm:"column:origin"`
	IsFavorite   bool   `gorm:"column:is_favorite"`
	DerivationId int    `gorm:"column:derivation_id"`
	Question     string `gorm:"column:question"`
	Choices      string `gorm:"column:choices"`
	Answer       string `gorm:"column:answer"`
	AnswerIndex  int    `gorm:"column:answer_index"`
	Analysis     string `gorm:"column:analysis"`
	CreateTime   string `gorm:"column:create_time;type:timestamp"`
	IsDoneRight  bool   `gorm:"column:is_done_right"`
}
