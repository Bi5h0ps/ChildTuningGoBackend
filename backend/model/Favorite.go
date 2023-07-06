package model

type Favorite struct {
	ID            int64  `gorm:"primaryKey;autoIncrement;column:id"`
	Username      string `gorm:"column:username"`
	Origin        string `gorm:"column:origin"` //"normal" or "asking"
	OriginId      string `gorm:"column:origin_id"`
	Question      string `gorm:"column:question"`
	Choices       string `gorm:"column:choices"`
	Answer        string `gorm:"column:answer"`
	AnswerIndex   int    `gorm:"column:answer_index"`
	Analysis      string `gorm:"column:analysis"`
	HasDerivation bool   `gorm:"column:has_derivation"`
	CreateTime    string `gorm:"column:create_time;type:timestamp"`
	IsDeleted     bool   `gorm:"column:is_deleted"`
}
