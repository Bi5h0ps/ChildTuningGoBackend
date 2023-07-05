package model

type ChatHistory struct {
	ID         int64  `gorm:"primaryKey;autoIncrement;column:ID"`
	QuestionId string `gorm:"column:question_id"`
	Username   string `gorm:"column:username"`
	Name       string `gorm:"column:name"`
	Message    string `gorm:"column:msg"`
	IsSelf     bool   `gorm:"column:isSelf"`
	CreateTime string `gorm:"column:create_time;type:timestamp"`
	IsFavorite bool   `gorm:"column:is_favorite"`
}

type ChatHistoryResponse struct {
	Name    string `json:"name"`
	Message string `json:"msg"`
	IsSelf  bool   `json:"isSelf"`
}
