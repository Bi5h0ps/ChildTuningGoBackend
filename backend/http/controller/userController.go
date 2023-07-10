package controller

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/provider"
	"ChildTuningGoBackend/backend/service"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UserController struct {
	ChatService      service.IChatService
	FavoriteService  service.IFavoriteService
	DerivedService   service.IDerivedService
	ExHistoryService service.IExHistoryService
}

const (
	TAG_ME     = "me"
	TAG_GPT    = "VT"
	TAG_ASKING = "asking"
)

func (c *UserController) PostAsk(ctx *gin.Context) {
	questionBody := map[string]string{
		"questionId": "",
		"question":   "",
	}
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
	saveChatHistory(questionBody["questionId"], true, username, questionBody["question"], c.ChatService)

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
		//todo: figure out server side error message
		errorHandling(response.Status, response.Message, ctx)
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": response.Message,
			"data": map[string]string{
				"questionId": questionBody["questionId"],
				"answer":     response.Content,
			},
		})
		saveChatHistory(questionBody["questionId"], false, username, response.Content, c.ChatService)
		return
	}
}

func saveChatHistory(questionId string, isClient bool, username, msg string, s service.IChatService) {
	//write chat history to database
	var tag string
	if isClient {
		tag = TAG_ME
	} else {
		tag = TAG_GPT
	}
	history := model.ChatHistory{
		QuestionId: questionId,
		Username:   username,
		Name:       tag,
		Message:    msg,
		IsSelf:     isClient,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		IsFavorite: false,
	}
	err := s.WriteChatHistory(&history)
	if err != nil {
		if err != nil {
			log.Default().Println(err.Error())
		}
	}
}

func (c *UserController) GetChatHistory(ctx *gin.Context) {
	//user info should be stored in the context by the auth middleware
	username := ctx.GetString("username")
	if username == "" {
		errorHandling(http.StatusBadRequest, "User not signed in, middleware uncaught error", ctx)
		return
	}
	//username ok, get user chat history from database
	data, err := c.ChatService.GetChatHistoryByUsername(username)
	if err != nil {
		errorHandling(http.StatusBadRequest, "Failed to retrieve chat history", ctx)
		return
	}
	dataTrimed := make([]map[string]interface{}, 0)
	for _, v := range data {
		dataTrimed = append(dataTrimed, map[string]interface{}{
			"questionId": v.QuestionId,
			"name":       v.Name,
			"msg":        v.Message,
			"isSelf":     v.IsSelf,
			"isFavorite": v.IsFavorite,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": dataTrimed,
	})
}

func (c *UserController) PostUserRandomQuiz(ctx *gin.Context) {
	randomQuizBody := map[string]string{
		"difficulty":    "",
		"type":          "",
		"questionCount": "",
	}
	if err := ctx.ShouldBindJSON(&randomQuizBody); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	url := "http://18.163.79.71:5000/generate"

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
		for i, _ := range response.Data {
			id := uuid.New().String()
			response.Data[i].ExerciseId = id
		}
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  response.Message,
			"data": response.Data,
		})
		return
	}
}

func (c *UserController) PostFavoriteQuestion(ctx *gin.Context) {
	payload := map[string]string{
		"questionId": "",
	}
	username := ctx.GetString("username")
	if username == "" {
		errorHandling(http.StatusBadRequest, "User not signed in, middleware uncaught error", ctx)
		return
	}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	pair, err := c.ChatService.GetChatHistoryByQuestionId(payload["questionId"])
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	} else if len(pair) == 0 {
		errorHandling(http.StatusInternalServerError, "QuestionId does not exists", ctx)
		return
	}
	var question string
	var analysis string
	if pair[0].Name == TAG_ME {
		question = pair[0].Message
		analysis = pair[1].Message
	} else {
		question = pair[1].Message
		analysis = pair[0].Message
	}
	err = c.ChatService.FavoriteChatHistory(payload["questionId"])
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}
	favoriteBase := &model.Favorite{
		Username:      username,
		Origin:        TAG_ASKING,
		OriginId:      payload["questionId"],
		Question:      question,
		Choices:       "",
		Answer:        "",
		AnswerIndex:   0,
		Analysis:      analysis,
		HasDerivation: false,
		CreateTime:    time.Now().Format("2006-01-02 15:04:05"),
		IsDeleted:     false,
	}
	err = c.FavoriteService.FavoriteAsking(favoriteBase)
	if err != nil {
		//TODO rollback the changes made on chat history
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": nil,
	})
	return
}

func (c *UserController) PostUnFavoriteQuestion(ctx *gin.Context) {
	payload := map[string]string{
		"questionId": "",
	}
	username := ctx.GetString("username")
	if username == "" {
		errorHandling(http.StatusBadRequest, "User not signed in, middleware uncaught error", ctx)
		return
	}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	err := c.ChatService.UnFavoriteChatHistory(payload["questionId"])
	if err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	err = c.FavoriteService.RemoveAskingFavorite(payload["questionId"])
	if err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": nil,
	})
}

func (c *UserController) GetFavorite(ctx *gin.Context) {
	username := ctx.GetString("username")
	if username == "" {
		errorHandling(http.StatusBadRequest, "User not signed in, middleware uncaught error", ctx)
		return
	}
	favoriteList, err := c.FavoriteService.GetFavoriteList(username)
	if err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": favoriteList,
	})
}

func (c *UserController) PostGenerateQuestion(ctx *gin.Context) {
	payload := map[string]string{
		"id": "",
	}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	username := ctx.GetString("username")
	if username == "" {
		errorHandling(http.StatusBadRequest, "User not signed in, middleware uncaught error", ctx)
		return
	}
	id, err := strconv.Atoi(payload["id"])
	if err != nil {
		errorHandling(http.StatusBadRequest, "id int conversion failed", ctx)
		return
	}
	exercise, err := c.FavoriteService.GetFavoriteExercise(id)
	if err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	question := exercise.Question
	url := "http://18.163.79.71:5000/regenerate"
	// Create a new multipart/form-data payload
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form field "input_msg"
	err = writer.WriteField("question", question)
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
	var response model.GenerateResponse
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	}

	//unsuccessful
	if response.Status != 200 {
		errorHandling(response.Status, response.Msg, ctx)
		return
	} else {
		favorite, err := c.FavoriteService.GetFavoriteExercise(id)
		if err != nil {
			errorHandling(response.Status, response.Msg, ctx)
			return
		}
		favorite.HasDerivation = true
		c.FavoriteService.FavoriteUpdate(favorite)
		base := model.DerivedExercise{
			ID:          0,
			Username:    username,
			FavoriteId:  id,
			Question:    response.Data[0].Question,
			Choices:     strings.Join(response.Data[0].Choices, "|"),
			Answer:      response.Data[0].Answer,
			AnswerIndex: response.Data[0].AnswerIndex,
			Analysis:    response.Data[0].Analysis,
			CreateTime:  time.Now().Format("2006-01-02 15:04:05"),
			IsDone:      false,
			UserChoice:  0,
			IsDoneRight: false,
		}
		recordId, err := c.DerivedService.SaveNewDerived(&base)
		if err != nil {
			errorHandling(http.StatusInternalServerError, err.Error(), ctx)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "",
			"data": map[string]interface{}{
				"id":           recordId,
				"favoriteId":   base.FavoriteId,
				"question":     base.Question,
				"choices":      response.Data[0].Choices,
				"answer":       base.Answer,
				"answer_index": base.AnswerIndex,
				"analysis":     base.Analysis,
				"isDone":       false,
				"userChoice":   nil,
				"isDoneRight":  nil,
			},
		})
	}
}

func (c *UserController) PostDerivedQuestionDo(ctx *gin.Context) {
	payload := struct {
		Id          string `json:"id"`
		UserChoice  int    `json:"userChoice"`
		IsDoneRight bool   `json:"isDoneRight"`
	}{}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	username := ctx.GetString("username")
	if username == "" {
		errorHandling(http.StatusBadRequest, "User not signed in, middleware uncaught error", ctx)
		return
	}
	id, err := strconv.Atoi(payload.Id)
	if err != nil {
		errorHandling(http.StatusBadRequest, "id int conversion failed", ctx)
		return
	}
	derivedRecord, err := c.DerivedService.GetDerivedById(id)
	if err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	derivedRecord.IsDone = true
	derivedRecord.IsDoneRight = payload.IsDoneRight
	derivedRecord.UserChoice = payload.UserChoice
	c.DerivedService.SaveDerivedUpdate(derivedRecord)
	err = c.ExHistoryService.SaveExHistory(&model.ExerciseHistory{
		ExerciseId:   "",
		Username:     username,
		Origin:       TAG_DERIVED,
		IsFavorite:   false,
		DerivationId: derivedRecord.FavoriteId,
		Question:     derivedRecord.Question,
		Choices:      derivedRecord.Choices,
		Answer:       derivedRecord.Answer,
		AnswerIndex:  derivedRecord.AnswerIndex,
		Analysis:     derivedRecord.Analysis,
		CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
		IsDoneRight:  derivedRecord.IsDoneRight,
	})
	if err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "",
			"data": nil,
		})
	}
}

func (c *UserController) PostGetDerivation(ctx *gin.Context) {
	payload := map[string]string{
		"id": "",
	}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	username := ctx.GetString("username")
	if username == "" {
		errorHandling(http.StatusBadRequest, "User not signed in, middleware uncaught error", ctx)
		return
	}
	id, err := strconv.Atoi(payload["id"])
	if err != nil {
		errorHandling(http.StatusBadRequest, "id int conversion failed", ctx)
		return
	}
	favorite, err := c.FavoriteService.GetFavoriteExercise(id)
	if err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	if !favorite.HasDerivation {
		errorHandling(http.StatusBadRequest, "This favorite question has no derivations", ctx)
		return
	}
	derivedList, err := c.DerivedService.GetAllDerived(username, id)
	if err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": derivedList,
	})
}
