package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/hugohenrick/gtasks/database"
	"github.com/hugohenrick/gtasks/models"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		interceptorValue := formatInterceptorValues(c)

		if interceptorValue != "POST/user/login" {
			tokenString := c.GetHeader("authorization")
			if tokenString == "" {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": "token not provided",
				})
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}

				return []byte(os.Getenv("SECRET")), nil
			})

			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
				return
			}

			if token == nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": fmt.Errorf("invalid token"),
				})
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if float64(time.Now().Unix()) > claims["exp"].(float64) {
					c.AbortWithStatus(http.StatusUnauthorized)
				}

				var user models.User
				database.DB.First(&user, claims["user"])

				if user.ID == 0 {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
						"error": "unauthorized user",
					})
					return
				}

				c.Set("isManager", user.IsManager)
				c.Set("userId", user.ID)

				c.Next()

			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized user",
				})
				return
			}
		}
	}
}
