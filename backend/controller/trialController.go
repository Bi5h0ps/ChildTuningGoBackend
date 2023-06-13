package controller

import (
	"ChildTuningGoBackend/backend/model"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
)

type TrialController struct{}

func (c *TrialController) PostAsk(ctx *gin.Context) {
	var questionBody *model.AskBody
	if err := ctx.ShouldBindJSON(&questionBody); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}

	url := "http://18.163.79.71:5000/api"
	// Create a new HTTP request with the POST method
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}

	// Create a new multipart/form-data payload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form field "input_msg"
	err = writer.WriteField("input_msg", questionBody.Question)
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

	// Set the Content-Type header with the boundary from the multipart writer
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Set the request body with the multipart payload
	req.Body = io.NopCloser(body)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request and retrieve the response
	resp, err := client.Do(req)
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
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
	var randomQuizBody *model.RandomQuizRequest
	if err := ctx.ShouldBindJSON(&randomQuizBody); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	url := "http://18.163.79.71:5000/generate"
	// Create a new HTTP request with the POST method
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}

	// Create a new multipart/form-data payload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form field "input_msg"
	err = writer.WriteField("difficulty", randomQuizBody.Difficulty)
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}
	err = writer.WriteField("type", randomQuizBody.QuestionType)
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}
	//TODO hard coded to 1
	err = writer.WriteField("questionCount", "1")
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

	// Set the Content-Type header with the boundary from the multipart writer
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Set the request body with the multipart payload
	req.Body = io.NopCloser(body)

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request and retrieve the response
	resp, err := client.Do(req)
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
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
	if resp.StatusCode != 200 {
		errorCode, _ := strconv.Atoi(resp.Status)
		errorHandling(errorCode, response.Message, ctx)
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
