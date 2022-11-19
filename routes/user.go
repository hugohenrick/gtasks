package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hugohenrick/gtasks/controllers"
)

// AddUserRoutes adds users routes to gin router
func AddUserRoutes(router *gin.Engine) {
	router.GET("/", controllers.GetUsers)
	router.POST("/", controllers.CreateUser)
}
