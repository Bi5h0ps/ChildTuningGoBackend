package controller

import (
	"ChildTuningGoBackend/backend/model"
	"ChildTuningGoBackend/backend/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

type UserController struct {
	UserService service.IUserService
}

func (c *UserController) PostSignUp(ctx *gin.Context) {
	var user *model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}
	_, err := c.UserService.AddUser(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"data":  nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Sign up successfully!",
		"data": nil,
	})
	return
}

func (c *UserController) PostSignIn(ctx *gin.Context) {
	var userCredentials *model.UserLogin
	if err := ctx.ShouldBindJSON(&userCredentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	user, isOk := c.UserService.IsPwdSuccess(userCredentials.Username, userCredentials.Password)
	if !isOk {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "Incorrect username or password, please try again!",
			"data": nil,
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	secretKeyByte := []byte(os.Getenv("SECRET"))
	tokenString, err := token.SignedString(secretKeyByte)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Failed to generate a signed token",
			"data": nil,
		})
		return
	}
	//4. Write jwt token to cookie
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", tokenString, 3600, "", "", false, true)
	data := gin.H{
		"username": user.Username,
		"nickname": user.Nickname,
		"email":    user.Email,
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Sign in successfully!",
		"data": data,
	})
	return
}

func (c *UserController) GetSignOut(ctx *gin.Context) {
	cookie, err := ctx.Cookie("Authorization")
	if err != nil || cookie == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg":  "You have to sign in first!",
			"data": nil,
		})
		return
	}
	ctx.SetCookie("Authorization", "", -1, "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Sign out successfully!",
		"data": nil,
	})
	return
}
