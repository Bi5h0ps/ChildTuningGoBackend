package model

type RandomQuizRequest struct {
	Difficulty    string `json:"difficulty"`
	QuestionType  string `json:"type"`
	QuestionCount string `json:"questionCount"`
}
