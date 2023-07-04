package model

import "time"

type ExerciseHistory struct {
	ExerciseId   string    `gorm:"primaryKey;column:exercise_id"`
	Username     string    `gorm:"column:username"`
	Origin       string    `gorm:"column:origin"`
	IsFavorite   bool      `gorm:"column:is_favorite"`
	DerivationId int64     `gorm:"column:derivation_id"`
	Question     string    `gorm:"column:question"`
	Choices      string    `gorm:"column:string"`
	Answer       string    `gorm:"column:answer"`
	AnswerIndex  int       `gorm:"column:answer_index"`
	Analysis     string    `gorm:"column:analysis"`
	CreateTime   time.Time `gorm:"column:create_time"`
	IsDoneRight  bool      `gorm:"column:is_done_right"`
}
