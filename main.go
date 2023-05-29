package main

import (
	"ChildTuningGoBackend/backend/controller"
	"ChildTuningGoBackend/backend/repository"
	"ChildTuningGoBackend/backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	ginServer := gin.Default()
	//logger middleware
	ginServer.Use(gin.Logger())
	gin.SetMode(gin.DebugMode)
	ginServer.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error", gin.H{
			"message": "Requested routing not exist",
		})
	})
	ginServer.Use(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.HTML(http.StatusInternalServerError, "error", gin.H{
					"message": err,
				})
			}
		}()
		c.Next()
	})

	repoUser := repository.NewUserRepository(nil)
	serviceUser := service.NewUserService(repoUser)
	controllerUser := controller.UserController{UserService: serviceUser}
	ginServer.POST("/signUp", controllerUser.PostSignUp)
	ginServer.POST("/signIn", controllerUser.PostSignIn)
	ginServer.GET("/signOut", controllerUser.GetSignOut)

	controllerTrial := controller.TrialController{}
	groupTrial := ginServer.Group("trial")
	{
		groupTrial.POST("ask", controllerTrial.PostAsk)
	}

	ginServer.Run(":9990")
}
