package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hugohenrick/gtasks/controllers"
	"github.com/hugohenrick/gtasks/middlewares"
)

// AddUserRoutes adds users routes to gin router
func AddUserRoutes(router *gin.Engine) {
	router.GET("/user", middlewares.Authenticate(), controllers.GetUsers)
	router.POST("/user", controllers.CreateUser)
	router.POST("/user/login", controllers.LoginUser)
}
