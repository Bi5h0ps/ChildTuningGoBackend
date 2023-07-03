package model

type ResponseBody struct {
	Content string `json:"content"`
	Message string `json:"message"`
	Status  int    `json:"status"`
	UserId  string `json:"userId"`
}
