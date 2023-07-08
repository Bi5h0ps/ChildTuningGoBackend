package model

type DerivedExercise struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement;column:ID"`
	Username    string `json:"-" gorm:"column:username"`
	FavoriteId  int    `json:"favorite_id" gorm:"column:favorite_id"`
	Question    string `json:"question" gorm:"column:question"`
	Choices     string `json:"choices" gorm:"column:choices"`
	Answer      string `json:"answer" gorm:"column:answer"`
	AnswerIndex int    `json:"answer_index" gorm:"column:answer_index"`
	Analysis    string `json:"analysis" gorm:"column:analysis"`
	CreateTime  string `json:"-" gorm:"column:create_time"`
	IsDone      bool   `json:"isDone" gorm:"column:is_done"`
	UserChoice  int    `json:"userChoice" gorm:"column:user_choice"`
	IsDoneRight bool   `json:"isDoneRight" gorm:"column:is_done_right"`
}
