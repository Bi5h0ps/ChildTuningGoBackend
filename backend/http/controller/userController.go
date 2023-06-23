package controller

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/provider"
	"ChildTuningGoBackend/backend/service"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

type UserController struct {
	ChatService service.IChatService
}

const (
	TAG_ME  = "me"
	TAG_GPT = "VT"
)

func (c *UserController) PostAsk(ctx *gin.Context) {
	var questionBody *model.AskBody
	//user info should be stored in the context by the auth middleware
	username := ctx.GetString("username")
	if username == "" {
		errorHandling(http.StatusBadRequest, "User not signed in, middleware uncaught error", ctx)
		return
	}
	if err := ctx.ShouldBindJSON(&questionBody); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	saveChatHistory(true, username, questionBody.Question, c.ChatService)

	url := "http://18.163.79.71:5000/api"
	// Create a new multipart/form-data payload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form field "input_msg"
	err := writer.WriteField("input_msg", questionBody.Question)
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}

	// Close the multipart writer
	err = writer.Close()
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}

	respBody, err := provider.HttpClientProvider.Post(url, body, writer, ctx)
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}
	// Unmarshal the JSON response
	var response model.ResponseBody
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}

	//unsuccessful
	if response.Status != 200 {
		//todo: figure out server side error message
		errorHandling(response.Status, response.Message, ctx)
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  response.Message,
			"data": response.Content,
		})
		saveChatHistory(false, username, response.Content, c.ChatService)
		return
	}
}

func saveChatHistory(isClient bool, username, msg string, s service.IChatService) {
	//write chat history to database
	var tag string
	if isClient {
		tag = TAG_ME
	} else {
		tag = TAG_GPT
	}
	history := model.ChatHistory{
		ID:         0,
		Username:   username,
		Name:       tag,
		Message:    msg,
		IsSelf:     isClient,
		CreateTime: time.Now(),
	}
	_, err := s.WriteChatHistory(&history)
	if err != nil {
		if err != nil {
			log.Default().Println(err.Error())
		}
	}
}
