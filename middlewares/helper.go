package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func formatInterceptorValues(c *gin.Context) string {
	return c.Request.Method + c.FullPath()
}

func validateRequestMethod(interceptorValue string, values []string) bool {
	for _, route := range values {
		if strings.EqualFold(route, interceptorValue) {
			return true
		}

	}
	return false
}
