package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	taskMock "github.com/hugohenrick/gtasks/mock"
	"github.com/hugohenrick/gtasks/repository"
	"github.com/hugohenrick/gtasks/routes"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
)

func TestGetTasks(t *testing.T) {
	assert := assert.New(t)
	gin.SetMode(gin.TestMode)

	t.Run("Success: expect correct result", func(t *testing.T) {

		iTaskMock := new(taskMock.ITaskRepository)
		iTaskMock.On("FindTasks", tmock.Anything, tmock.Anything, tmock.Anything).Return(nil, nil)
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
