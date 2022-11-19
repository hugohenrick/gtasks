package middlewares

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func formatInterceptorValues(c *gin.Context) string {
	return c.Request.Method + c.FullPath()
}

func validateRequestMethod(interceptorValue string, values map[string]string) bool {
	serviceEnv := os.Getenv("SERVICE")
	for route, service := range values {
		if strings.EqualFold(serviceEnv, service) {
			if strings.EqualFold(route, interceptorValue) {
				return true
			}
		}
	}
	return false
}
