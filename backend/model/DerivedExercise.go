package model

import "time"

type DerivedExercise struct {
	ID          int64     `gorm:"primaryKey;autoIncrement;column:ID"`
	Username    string    `gorm:"column:username"`
	FavoriteId  int64     `gorm:"column:favorite_id"`
	Question    string    `gorm:"column:question"`
	Choices     string    `gorm:"column:choices"`
	Answer      string    `gorm:"column:answer"`
	AnswerIndex int       `gorm:"column:answer_index"`
	Analysis    string    `gorm:"column:analysis"`
	CreateTime  time.Time `gorm:"column:create_time"`
	IsDone      bool      `gorm:"column:is_done"`
	UserChoice  int       `gorm:"column:user_choice"`
	IsDoneRight bool      `gorm:"column:is_done_right"`
}
