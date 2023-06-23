package middleware

import (
	"ChildTuningGoBackend/backend/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
	"time"
)

func RequireAuth(service service.IUserService) gin.HandlerFunc {
	//TODO redirect to which page?
	location := "/user/login"
	return func(ctx *gin.Context) {
		tokenString, err := ctx.Cookie("Authorization")
		if err != nil {
			AuthFailedAndRedirect(ctx, location)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Method)
			}
			return []byte(os.Getenv("SECRET")), nil
		})
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			//check expiration
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				//token expired
				AuthFailedAndRedirect(ctx, location)
				return
			}
			user, err := service.GetUserByUsername(claims["sub"].(string))
			if err != nil || user == nil {
				//failed to get user
				AuthFailedAndRedirect(ctx, location)
				return
			}
			ctx.Set("username", user.Username)
			ctx.Next()
		} else {
			AuthFailedAndRedirect(ctx, location)
			return
		}
	}
}

func AuthFailedAndRedirect(ctx *gin.Context, location string) {
	ctx.Redirect(http.StatusTemporaryRedirect, location)
	ctx.Abort()
}
