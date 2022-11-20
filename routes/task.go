package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hugohenrick/gtasks/controllers"
)

// AddTaskRoutes adds tasks routes to gin router
func AddTaskRoutes(router *gin.Engine) {
	router.GET("/task", controllers.GetTasks)
	router.GET("/task/:id", controllers.GetTaskById)
	router.POST("/task", controllers.CreateTask)
	router.POST("/task/execute/:id", controllers.ExecuteTask)
	router.PATCH("/task/:id", controllers.UpdateTask)
	router.DELETE("/task/:id", controllers.DeleteTask)
}
