package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hugohenrick/gtasks/database"
	"github.com/hugohenrick/gtasks/models"
	"github.com/hugohenrick/gtasks/utils"
)

func GetTasks(c *gin.Context) {
	var tasks []models.Task
	var task models.Task

	isManagerRaw, ok := c.Get("isManager")
	if !ok || isManagerRaw == nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserNotFound))
		return
	}
	isManager := isManagerRaw.(bool)

	if !isManager {
		userIdRaw, ok := c.Get("userId")
		if !ok || userIdRaw == nil {
			utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserNotFound))
			return
		}

		userId := userIdRaw.(uint32)
		task.UserId = userId
	}

	database.DB.Preload("User").Find(&tasks, task)
	utils.SendJSONResponse(c, http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBind(&task); err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v: %v", utils.InvalidJsonProvided, err))
		return
	}

	if task.Title == "" {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.TaskTitleRequired))
		return
	}
	if task.Summary == "" {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.TaskSummaryRequired))
		return
	}

	if task.UserId == 0 {
		userIdRaw, ok := c.Get("userId")
		if !ok || userIdRaw == nil {
			utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserWithoutAccesPermission))
			return
		}

		userId := userIdRaw.(uint32)
		task.UserId = userId
	}

	database.DB.Create(&task)
	utils.SendJSONResponse(c, http.StatusOK, task)
}

func GetTaskById(c *gin.Context) {
	var task models.Task
	id := c.Param("id")
	if id == "" {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserIdRequired))
		return
	}

	database.DB.Preload("User").First(&task, id)
	utils.SendJSONResponse(c, http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")
	if id == "" {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserIdRequired))
		return
	}

	database.DB.First(&task, id)
	if task.ID == 0 {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.TaskNotFound))
		return
	}

	if err := c.ShouldBind(&task); err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v: %v", utils.InvalidJsonProvided, err))
		return
	}

	userIdRaw, ok := c.Get("userId")
	if !ok || userIdRaw == nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserWithoutAccesPermission))
		return
	}

	userId := userIdRaw.(uint32)

	if task.UserId != userId {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserCannotChangeTaskAnotherUser))
		return
	}

	database.DB.Save(&task)
	utils.SendJSONResponse(c, http.StatusOK, "success")
}

func DeleteTask(c *gin.Context) {
	var task models.Task

	isManagerRaw, ok := c.Get("isManager")
	if !ok || isManagerRaw == nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserNotFound))
		return
	}
	isManager := isManagerRaw.(bool)

	if !isManager {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserWithoutAccesPermission))
		return
	}

	id := c.Param("id")
	if id == "" {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserIdRequired))
		return
	}

	database.DB.First(&task, id)
	if task.ID == 0 {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.TaskNotFound))
		return
	}

	if err := c.ShouldBind(&task); err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v: %v", utils.InvalidJsonProvided, err))
		return
	}

	database.DB.Delete(&task)
	utils.SendJSONResponse(c, http.StatusOK, "success")
}

func ExecuteTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")
	if id == "" {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserIdRequired))
		return
	}

	database.DB.First(&task, id)
	if task.ID == 0 {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.TaskNotFound))
		return
	}

	userIdRaw, ok := c.Get("userId")
	if !ok || userIdRaw == nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserWithoutAccesPermission))
		return
	}

	userId := userIdRaw.(uint32)

	if task.UserId != userId {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserCannotChangeTaskAnotherUser))
		return
	}

	task.Done = true
	timeNow := time.Now()
	task.FinishedAt = &timeNow

	database.DB.Save(&task)
	utils.SendJSONResponse(c, http.StatusOK, "success")
}
