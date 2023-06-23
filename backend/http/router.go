package http

import (
	"ChildTuningGoBackend/backend/http/controller"
	"ChildTuningGoBackend/backend/http/middleware"
	"ChildTuningGoBackend/backend/repository"
	"ChildTuningGoBackend/backend/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Router struct {
	ginServer *gin.Engine
}

func NewRouter() Router {
	return Router{ginServer: gin.Default()}
}

func (r *Router) StartServer() {
	//logger middleware
	r.ginServer.Use(gin.Logger())
	gin.SetMode(gin.DebugMode)
	r.ginServer.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error", gin.H{
			"message": "Requested routing not exist",
		})
	})
	r.ginServer.Use(func(c *gin.Context) {
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
	controllerOnBoarding := controller.OnBoardingController{UserService: serviceUser}
	r.ginServer.POST("/signUp", controllerOnBoarding.PostSignUp)
	r.ginServer.POST("/signIn", controllerOnBoarding.PostSignIn)
	r.ginServer.GET("/signOut", controllerOnBoarding.GetSignOut)

	controllerTrial := controller.TrialController{}
	groupTrial := r.ginServer.Group("trial")
	{
		groupTrial.POST("ask", controllerTrial.PostAsk)
		groupTrial.POST("exercise", controllerTrial.PostRandomQuiz)
	}

	repoChat := repository.NewChatRepository(nil)
	serviceChat := service.NewChatService(repoChat)
	controllerChat := controller.UserController{ChatService: serviceChat}
	groupUser := r.ginServer.Group("user")
	groupUser.Use(middleware.RequireAuth(serviceUser))
	{
		groupUser.GET("/askingHistory", controllerChat.GetChatHistory)
		groupUser.POST("/ask", controllerChat.PostAsk)
	}

	r.ginServer.Run(":9990")
}
