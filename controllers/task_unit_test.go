package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	taskMock "github.com/hugohenrick/gtasks/mock"
	"github.com/hugohenrick/gtasks/models"
	"github.com/hugohenrick/gtasks/repository"
	"github.com/hugohenrick/gtasks/routes"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
)

func TestGetTasks(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	t.Run("Failed: user not found", func(t *testing.T) {
		expectMsgError := `{"error":"user not found"}`

		iTaskMock := new(taskMock.ITaskRepository)
		iTaskMock.On("FindTasks", tmock.Anything).Return(nil, nil)
		repository.TaskRepositoryServices = iTaskMock

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		router.Use(func(c *gin.Context) {})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodGet, "/task", nil)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(expectMsgError, w.Body.String())
	})

	t.Run("Success: expect correct result", func(t *testing.T) {

		iTaskMock := new(taskMock.ITaskRepository)
		iTaskMock.On("FindTasks", tmock.Anything).Return(nil, nil)
		repository.TaskRepositoryServices = iTaskMock

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {
			c.Set("isManager", false)
			c.Set("userId", uint32(1))
		})

		router.Use(func(c *gin.Context) {})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodGet, "/task", nil)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusOK, w.Code)
	})
}

func TestGetTaskById(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	taskModel := models.Task{
		ID:      1,
		Title:   "Test Title",
		Summary: "Test Summary",
		UserId:  1,
	}

	t.Run("Failed: user id is required", func(t *testing.T) {
		expectMsgError := `{"error":"user id is required"}`

		iTaskMock := new(taskMock.ITaskRepository)
		iTaskMock.On("FindTaskById", tmock.Anything).Return(taskModel, nil)
		repository.TaskRepositoryServices = iTaskMock

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)

		router.Use(func(c *gin.Context) {})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodGet, "/task/ ", nil)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(expectMsgError, w.Body.String())
	})

	t.Run("Success: expect correct result", func(t *testing.T) {
		iTaskMock := new(taskMock.ITaskRepository)
		iTaskMock.On("FindTaskById", tmock.Anything).Return(taskModel, nil)
		repository.TaskRepositoryServices = iTaskMock

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {
			c.Set("isManager", false)
			c.Set("userId", uint32(1))
		})

		router.Use(func(c *gin.Context) {})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodGet, "/task/1", nil)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusOK, w.Code)
	})
}

func TestCreateTask(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	t.Run("Failed: task title is required", func(t *testing.T) {
		// expect error msg
		expectMsgError := `{"error":"task title is required"}`

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {})
		routes.AddTaskRoutes(router)

		taskModel := models.Task{}
		data, _ := json.Marshal(taskModel)
		reader := bytes.NewReader(data)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodPost, "/task", reader)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(expectMsgError, w.Body.String())
	})

	t.Run("Failed: task summary is required", func(t *testing.T) {
		// expect error msg
		expectMsgError := `{"error":"task summary is required"}`

		taskModel := models.Task{
			ID:    1,
			Title: "Test Title",
		}
		data, _ := json.Marshal(taskModel)
		body := bytes.NewBuffer(data)

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodPost, "/task", body)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(expectMsgError, w.Body.String())
	})

	t.Run("Failed: user without access permission", func(t *testing.T) {
		// expect error msg
		expectMsgError := `{"error":"user without access permission"}`

		taskModel := models.Task{
			ID:      1,
			Title:   "Test Title",
			Summary: "Test Summary",
		}
		data, _ := json.Marshal(taskModel)
		body := bytes.NewBuffer(data)

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodPost, "/task", body)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(expectMsgError, w.Body.String())
	})

	t.Run("Succes: create new task", func(t *testing.T) {
		taskModel := models.Task{
			ID:      1,
			Title:   "Test Title",
			Summary: "Test Summary",
		}

		iTaskMock := new(taskMock.ITaskRepository)
		iTaskMock.On("CreateTask", tmock.Anything).Return(taskModel, nil)
		repository.TaskRepositoryServices = iTaskMock

		data, _ := json.Marshal(taskModel)
		body := bytes.NewBuffer(data)

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {
			c.Set("userId", uint32(1))
		})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodPost, "/task", body)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusOK, w.Code)
	})
}

func TestUpdateTask(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	t.Run("Failed: task title is required", func(t *testing.T) {
		// expect error msg
		expectMsgError := `{"error":"task title is required"}`

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {})
		routes.AddTaskRoutes(router)

		taskModel := models.Task{}
		data, _ := json.Marshal(taskModel)
		reader := bytes.NewReader(data)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodPatch, "/task/1", reader)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(expectMsgError, w.Body.String())
	})

	t.Run("Failed: task summary is required", func(t *testing.T) {
		// expect error msg
		expectMsgError := `{"error":"task summary is required"}`

		taskModel := models.Task{
			ID:    1,
			Title: "Test Title",
		}
		data, _ := json.Marshal(taskModel)
		body := bytes.NewBuffer(data)

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodPatch, "/task/1", body)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(expectMsgError, w.Body.String())
	})

	t.Run("Failed: user without access permission", func(t *testing.T) {
		// expect error msg
		expectMsgError := `{"error":"user without access permission"}`

		taskModel := models.Task{
			ID:      1,
			Title:   "Test Title",
			Summary: "Test Summary",
		}
		data, _ := json.Marshal(taskModel)
		body := bytes.NewBuffer(data)

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodPatch, "/task/1", body)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(expectMsgError, w.Body.String())
	})

	t.Run("Failed: user cannot change a task of another user", func(t *testing.T) {
		// expect error msg
		expectMsgError := `{"error":"user cannot change a task of another user"}`

		taskModel := models.Task{
			ID:      1,
			Title:   "Test Title",
			Summary: "Test Summary",
		}
		data, _ := json.Marshal(taskModel)
		body := bytes.NewBuffer(data)

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {
			c.Set("userId", uint32(1))
		})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodPatch, "/task/1", body)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(expectMsgError, w.Body.String())
	})

	t.Run("Succes: update task", func(t *testing.T) {
		taskModel := models.Task{
			ID:      1,
			Title:   "Test Title",
			Summary: "Test Summary",
			UserId:  1,
		}

		iTaskMock := new(taskMock.ITaskRepository)
		iTaskMock.On("UpdateTask", tmock.Anything, tmock.Anything).Return(taskModel, nil)
		repository.TaskRepositoryServices = iTaskMock

		data, _ := json.Marshal(taskModel)
		body := bytes.NewBuffer(data)

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {
			c.Set("userId", uint32(1))
		})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodPatch, "/task/1", body)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusOK, w.Code)
	})
}

func TestDeleteTask(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)
	var numRecord int64 = 1

	t.Run("Failed: user not found", func(t *testing.T) {
		expectMsgError := `{"error":"user not found"}`

		iTaskMock := new(taskMock.ITaskRepository)
		var numRecord int64 = 1
		iTaskMock.On("DeleteTask", tmock.Anything).Return(numRecord, nil)
		repository.TaskRepositoryServices = iTaskMock

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodDelete, "/task/1", nil)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(expectMsgError, w.Body.String())
	})

	t.Run("Failed: user without access permission", func(t *testing.T) {
		expectMsgError := `{"error":"user without access permission"}`

		iTaskMock := new(taskMock.ITaskRepository)
		iTaskMock.On("DeleteTask", tmock.Anything).Return(numRecord, nil)
		repository.TaskRepositoryServices = iTaskMock

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {
			c.Set("isManager", false)
			c.Set("userId", uint32(1))
		})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodDelete, "/task/1", nil)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(expectMsgError, w.Body.String())
	})

	t.Run("Failed: user id is required", func(t *testing.T) {
		expectMsgError := `{"error":"user id is required"}`

		iTaskMock := new(taskMock.ITaskRepository)
		iTaskMock.On("DeleteTask", tmock.Anything).Return(numRecord, nil)
		repository.TaskRepositoryServices = iTaskMock

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {
			c.Set("isManager", true)
			c.Set("userId", uint32(1))
		})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodDelete, "/task/ ", nil)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(expectMsgError, w.Body.String())
	})

	t.Run("Success: expect correct result", func(t *testing.T) {
		iTaskMock := new(taskMock.ITaskRepository)
		iTaskMock.On("DeleteTask", tmock.Anything).Return(numRecord, nil)
		repository.TaskRepositoryServices = iTaskMock

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {
			c.Set("isManager", true)
			c.Set("userId", uint32(1))
		})

		router.Use(func(c *gin.Context) {})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodDelete, "/task/1", nil)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusOK, w.Code)
	})
}

func TestExecuteTask(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	t.Run("Failed: user without access permission", func(t *testing.T) {
		// expect error msg
		expectMsgError := `{"error":"user without access permission"}`

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {})
		routes.AddTaskRoutes(router)

		taskModel := models.Task{}
		data, _ := json.Marshal(taskModel)
		reader := bytes.NewReader(data)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodPatch, "/task/execute/1", reader)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusBadRequest, w.Code)
		assert.Equal(expectMsgError, w.Body.String())
	})

	t.Run("Succes: execute task", func(t *testing.T) {
		taskModel := models.Task{
			ID:      1,
			Title:   "Test Title",
			Summary: "Test Summary",
			UserId:  1,
		}

		iTaskMock := new(taskMock.ITaskRepository)
		iTaskMock.On("FindTaskById", tmock.Anything).Return(taskModel, nil)
		iTaskMock.On("ExecuteTask", tmock.Anything, tmock.Anything).Return(taskModel, nil)
		repository.TaskRepositoryServices = iTaskMock

		data, _ := json.Marshal(taskModel)
		body := bytes.NewBuffer(data)

		w := httptest.NewRecorder()
		c, router := gin.CreateTestContext(w)
		router.Use(func(c *gin.Context) {
			c.Set("userId", uint32(1))
		})

		routes.AddTaskRoutes(router)

		// creating a request to send on endpoint call
		c.Request, _ = http.NewRequest(http.MethodPatch, "/task/execute/1", body)

		// endpoint call
		router.ServeHTTP(w, c.Request)

		// asserts
		assert.Equal(http.StatusOK, w.Code)
	})
}
