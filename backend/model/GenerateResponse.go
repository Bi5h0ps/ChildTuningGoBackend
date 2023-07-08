package model

type GenerateResponse struct {
	Data   []Derived `json:"data"`
	Msg    string    `json:"msg"`
	Status int       `json:"status"`
	UserId string    `json:"userId"`
}

type Derived struct {
	Analysis    string   `json:"analysis"`
	Answer      string   `json:"answer"`
	AnswerIndex int      `json:"answer_index"`
	Choices     []string `json:"choices"`
	Question    string   `json:"question"`
}
