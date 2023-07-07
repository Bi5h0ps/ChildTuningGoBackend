package controller

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type ExHistoryController struct {
	ExerciseService service.IExHistoryService
	FavoriteSerivce service.IFavoriteService
}

const (
	TAG_NORMAL  = "normal"
	TAG_DERIVED = "derived"
)

func (c *ExHistoryController) PostExerciseDo(ctx *gin.Context) {
	var requestPayload model.ExerciseHistoryRequest
	//user info should be stored in the context by the auth middleware
	username := ctx.GetString("username")
	if username == "" {
		errorHandling(http.StatusBadRequest, "User not signed in, middleware uncaught error", ctx)
		return
	}
	if err := ctx.ShouldBindJSON(&requestPayload); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	exerciseHistory := model.ExerciseHistory{
		ExerciseId:   requestPayload.ExerciseID,
		Username:     username,
		Origin:       TAG_NORMAL,
		IsFavorite:   false,
		DerivationId: -1,
		Question:     requestPayload.Question,
		Choices:      strings.Join(requestPayload.Choices, "/"),
		Answer:       requestPayload.Answer,
		AnswerIndex:  requestPayload.AnswerIndex,
		Analysis:     requestPayload.Analysis,
		CreateTime:   time.Now().Format("2006-01-02 15:04:05"),
		IsDoneRight:  requestPayload.IsDoneRight,
	}
	err := c.ExerciseService.SaveExHistory(&exerciseHistory)
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "",
			"data": nil,
		})
	}
}

func (c *ExHistoryController) GetExerciseHistory(ctx *gin.Context) {
	//user info should be stored in the context by the auth middleware
	username := ctx.GetString("username")
	if username == "" {
		errorHandling(http.StatusBadRequest, "User not signed in, middleware uncaught error", ctx)
		return
	}
	historyList, err := c.ExerciseService.GetExeHistoriesByUsername(username)
	if err != nil {
		errorHandling(http.StatusInternalServerError, err.Error(), ctx)
		return
	} else {
		data := []map[string]interface{}{}
		for _, v := range historyList {
			data = append(data, map[string]interface{}{
				//TODO add logic here to determine which id for using
				"id":           v.ExerciseId,
				"origin":       v.Origin,
				"is_favorite":  v.IsFavorite,
				"question":     v.Question,
				"choices":      strings.Split(v.Choices, "/"),
				"answer":       v.Answer,
				"answer_index": v.AnswerIndex,
				"analysis":     v.Analysis,
				"isDoneRight":  v.IsDoneRight,
				"createTime":   v.CreateTime,
			})
		}
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "",
			"data": data,
		})
	}
}

func (c *ExHistoryController) PostFavoriteExercise(ctx *gin.Context) {
	//user info should be stored in the context by the auth middleware
	username := ctx.GetString("username")
	if username == "" {
		errorHandling(http.StatusBadRequest, "User not signed in, middleware uncaught error", ctx)
		return
	}
	payload := map[string]string{
		"exerciseId": "",
	}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	exHistory, err := c.ExerciseService.GetExHistoryById(payload["exerciseId"])
	if err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	exHistory.IsFavorite = true
	c.ExerciseService.SaveExHistoryUpdate(exHistory)
	exFavorite := &model.Favorite{
		ID:            0,
		Username:      username,
		Origin:        TAG_NORMAL,
		OriginId:      payload["exerciseId"],
		Question:      exHistory.Question,
		Choices:       exHistory.Choices,
		Answer:        exHistory.Answer,
		AnswerIndex:   exHistory.AnswerIndex,
		Analysis:      exHistory.Analysis,
		HasDerivation: false,
		CreateTime:    time.Now().Format("2006-01-02 15:04:05"),
		IsDeleted:     false,
	}
	err = c.FavoriteSerivce.FavoriteExercise(exFavorite)
	if err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": nil,
	})
	return
}

func (c *ExHistoryController) PostUnFavoriteExercise(ctx *gin.Context) {
	//user info should be stored in the context by the auth middleware
	username := ctx.GetString("username")
	if username == "" {
		errorHandling(http.StatusBadRequest, "User not signed in, middleware uncaught error", ctx)
		return
	}
	payload := map[string]string{
		"exerciseId": "",
	}
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	exHistory, err := c.ExerciseService.GetExHistoryById(payload["exerciseId"])
	if err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	exHistory.IsFavorite = false
	c.ExerciseService.SaveExHistoryUpdate(exHistory)
	err = c.FavoriteSerivce.RemoveExerciseFavorite(payload["exerciseId"])
	if err != nil {
		errorHandling(http.StatusBadRequest, err.Error(), ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": nil,
	})
	return
}
