package model

import "time"

type ChatHistory struct {
	ID         int64     `json:"id" form:"ID" gorm:"primaryKey;autoIncrement;column:ID"`
	Username   string    `json:"username" gorm:"primaryKey;column:username"`
	Name       string    `json:"name" gorm:"column:name"`
	Message    string    `json:"msg" gorm:"column:msg"`
	IsSelf     bool      `json:"isSelf" gorm:"column:isSelf"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
}
