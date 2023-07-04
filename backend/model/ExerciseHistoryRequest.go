package model

type ExerciseHistoryRequest struct {
	ExerciseID  string   `json:"exerciseId"`
	Question    string   `json:"question"`
	Choices     []string `json:"choices"`
	Answer      string   `json:"answer"`
	AnswerIndex int      `json:"answer_index"`
	Analysis    string   `json:"analysis"`
	IsDoneRight bool     `json:"isDoneRight"`
}
