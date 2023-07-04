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
		CreateTime:   time.Now(),
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
				"exerciseId":   v.ExerciseId,
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
