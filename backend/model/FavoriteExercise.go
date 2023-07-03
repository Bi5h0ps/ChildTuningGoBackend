package model

import "time"

type FavoriteExercise struct {
	ID            int64     `gorm:"primaryKey;autoIncrement;column:ID"`
	Username      string    `gorm:"column:username"`
	Origin        string    `gorm:"column:origin"` //"normal" or "asking"
	OriginId      string    `gorm:"column:origin_id"`
	Question      string    `gorm:"column:question"`
	Choices       string    `gorm:"column:choices"`
	Answer        string    `gorm:"column:answer"`
	AnswerIndex   int       `gorm:"column:answer_index"`
	Analysis      string    `gorm:"column:analysis"`
	HasDerivation bool      `gorm:"column:has_derivation"`
	CreateTime    time.Time `gorm:"column:create_time"`
	IsDeleted     bool      `gorm:"column:is_deletedd"`
}
