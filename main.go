package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hugohenrick/gtasks/database"
	"github.com/hugohenrick/gtasks/routes"
	"github.com/hugohenrick/gtasks/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func init() {
	utils.LoadEnv()
}

func main() {

	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	httpPort := ":8080"
	if os.Getenv("SERVER_PORT") != "" {
		httpPort = ":" + os.Getenv("SERVER_PORT")
	}

	switch os.Getenv("SERVICE") {
	case "users":
		routes.AddUserRoutes(router)
	case "tasks":
		routes.AddTaskRoutes(router)
	default:
		fmt.Println("SERVICE env var must be one of [ users, tasks ]")
		os.Exit(1)
	}

	server := &http.Server{
		Addr:    httpPort,
		Handler: router,
	}

	database.Conn()
	// start API server
	//go func() {
	fmt.Printf("%s API listening on port %s\n", cases.Title(language.AmericanEnglish).String(os.Getenv("SERVICE")), httpPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		os.Exit(1)
	}
	//}()

}
