package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/hugohenrick/gtasks/models"
	"github.com/hugohenrick/gtasks/repository"
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

	tasks, err := repository.TaskRepositoryServices.FindTasks(task)
	if err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v: %v", utils.InvalidJsonProvided, err))
		return
	}

	utils.SendJSONResponse(c, http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindWith(&task, binding.JSON); err != nil {
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

	task, err := repository.TaskRepositoryServices.CreateTask(task)
	if err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", err))
		return
	}

	utils.SendJSONResponse(c, http.StatusOK, task)
}

func GetTaskById(c *gin.Context) {
	var task models.Task
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserIdRequired))
		return
	}

	task, err := repository.TaskRepositoryServices.FindTaskById(fmt.Sprint(id))
	if err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", err))
		return
	}

	utils.SendJSONResponse(c, http.StatusOK, task)
}

func UpdateTask(c *gin.Context) {
	var task models.Task
	id := c.Param("id")
	if id == "" {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserIdRequired))
		return
	}

	if err := c.ShouldBindWith(&task, binding.JSON); err != nil {
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

	_, err := repository.TaskRepositoryServices.UpdateTask(fmt.Sprint(id), task)
	if err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", err))
		return
	}

	utils.SendJSONResponse(c, http.StatusOK, "success")
}

func DeleteTask(c *gin.Context) {
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

	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserIdRequired))
		return
	}

	_, err := repository.TaskRepositoryServices.DeleteTask(fmt.Sprint(id))
	if err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", err))
		return
	}

	utils.SendJSONResponse(c, http.StatusOK, "success")
}

func ExecuteTask(c *gin.Context) {
	var task models.Task
	id := strings.TrimSpace(c.Param("id"))
	if id == "" {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserIdRequired))
		return
	}

	userIdRaw, ok := c.Get("userId")
	if !ok || userIdRaw == nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserWithoutAccesPermission))
		return
	}

	userId := userIdRaw.(uint32)

	task, err := repository.TaskRepositoryServices.FindTaskById(fmt.Sprint(id))
	if err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", err))
		return
	}

	if task.UserId != userId {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", utils.UserCannotChangeTaskAnotherUser))
		return
	}

	task.Done = true
	timeNow := time.Now()
	task.FinishedAt = &timeNow

	_, err = repository.TaskRepositoryServices.ExecuteTask(id, task)
	if err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v", err))
		return
	}

	utils.SendJSONResponse(c, http.StatusOK, "success")
}
