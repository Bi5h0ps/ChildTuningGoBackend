package controller

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/provider"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
)

type TrialController struct{}

func (c *TrialController) PostAsk(ctx *gin.Context) {
	questionBody := map[string]string{
		"question": "",
	}
	if err := ctx.ShouldBindJSON(&questionBody); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}

	url := "http://18.163.79.71:5000/api"
	// Create a new multipart/form-data payload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form field "input_msg"
	err := writer.WriteField("input_msg", questionBody["question"])
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
		errorHandling(response.Status, response.Message, ctx)
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  response.Message,
			"data": response.Content,
		})
		return
	}
}

func (c *TrialController) PostRandomQuiz(ctx *gin.Context) {
	randomQuizBody := map[string]string{
		"difficulty":    "",
		"type":          "",
		"questionCount": "",
	}
	if err := ctx.ShouldBindJSON(&randomQuizBody); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	url := "http://18.163.79.71:5000/agenerate"

	// Create a new multipart/form-data payload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form field "input_msg"
	err := writer.WriteField("difficulty", randomQuizBody["difficulty"])
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}
	err = writer.WriteField("type", randomQuizBody["type"])
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}
	err = writer.WriteField("questionCount", randomQuizBody["questionCount"])
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

	// Read the response body
	respBody, err := provider.HttpClientProvider.Post(url, body, writer, ctx)
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}

	// Unmarshal the JSON response
	var response model.RandomQuizResponse
	err = json.Unmarshal(respBody, &response)

	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}

	//unsuccessful
	if response.Status != 200 {
		errorHandling(response.Status, response.Message, ctx)
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  response.Message,
			"data": response.Data,
		})
		return
	}
}

func errorHandling(code int, msg string, ctx *gin.Context) {
	ctx.JSON(code, gin.H{
		"msg":  "",
		"data": msg,
	})
	return
}
