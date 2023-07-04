package model

type RandomQuizResponse struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    []Quiz `json:"data"`
	UserId  string `json:"userId"`
}

type Quiz struct {
	ExerciseId  string   `json:"exerciseId"`
	Question    string   `json:"question"`
	Choices     []string `json:"choices"`
	Answer      string   `json:"answer"`
	AnswerIndex int      `json:"answer_index"`
	Analysis    string   `json:"analysis"`
}
