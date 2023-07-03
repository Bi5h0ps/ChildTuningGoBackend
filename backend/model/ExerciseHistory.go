package model

import "time"

type ExerciseHistory struct {
	ID           int64     `gorm:"primaryKey;autoIncrement;column:ID"`
	Username     string    `gorm:"column:username"`
	Origin       string    `gorm:"column:origin"`
	ExerciseId   int64     `gorm:"column:exercise_id"`
	IsFavorite   bool      `gorm:"column:is_favorite"`
	DerivationId int64     `gorm:"column:derivation_id"`
	Question     string    `gorm:"column:question"`
	Choices      string    `gorm:"column:string"`
	Answer       string    `gorm:"column:answer"`
	AnswerIndex  int       `gorm:"column:answer_index"`
	Analysis     string    `gorm:"column:analysis"`
	CreateTime   time.Time `gorm:"column:create_time"`
	IsDone       bool      `gorm:"column:is_done"`
	UserChoice   int       `gorm:"column:user_choice"`
	IsDoneRight  bool      `gorm:"column:is_done_right"`
}
