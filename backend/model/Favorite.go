package model

type Favorite struct {
	ID            int64  `json:"-" gorm:"primaryKey;autoIncrement;column:id"`
	Username      string `json:"-" gorm:"column:username"`
	Origin        string `json:"origin" gorm:"column:origin"` //"normal" or "asking"
	OriginId      string `json:"id" gorm:"column:origin_id"`
	Question      string `json:"question" gorm:"column:question"`
	Choices       string `json:"choices" gorm:"column:choices"`
	Answer        string `json:"answer" gorm:"column:answer"`
	AnswerIndex   int    `json:"answer_index" gorm:"column:answer_index"`
	Analysis      string `json:"analysis" gorm:"column:analysis"`
	HasDerivation bool   `json:"has_derivation" gorm:"column:has_derivation"`
	CreateTime    string `json:"create_time" gorm:"column:create_time;type:timestamp"`
	IsDeleted     bool   `gorm:"column:is_deleted"`
}
