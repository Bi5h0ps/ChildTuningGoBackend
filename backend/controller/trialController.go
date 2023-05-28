package controller

import (
	"ChildTuningGoBackend/backend/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TrialController struct{}

func (c *TrialController) GetAsk(ctx *gin.Context) {
	var questionBody *model.AskBody
	if err := ctx.ShouldBindJSON(&questionBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	response := fmt.Sprintf("Question: %v Answer: %v",
		questionBody.Question, "python backend is currently under construction")
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "",
		"data": response,
	})
	return
}
