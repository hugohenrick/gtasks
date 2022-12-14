package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hugohenrick/gtasks/database"
	"github.com/hugohenrick/gtasks/middlewares"
	"github.com/hugohenrick/gtasks/rabbitmq"
	"github.com/hugohenrick/gtasks/repository"
	"github.com/hugohenrick/gtasks/routes"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// func init() {
// 	utils.LoadEnv()
// }

func main() {

	fmt.Printf("INICIOOOOO API listening on port")

	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     cors.DefaultConfig().AllowMethods,
		AllowHeaders:     cors.DefaultConfig().AllowHeaders,
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	httpPort := ":8080"
	if os.Getenv("SERVER_PORT") != "" {
		httpPort = ":" + os.Getenv("SERVER_PORT")
	}

	router.Use(middlewares.Authenticate())

	database.Conn()

	//Microservices:
	switch os.Getenv("SERVICE") {
	case "users":
		routes.AddUserRoutes(router)
		repository.UserRepositoryServices = repository.NewUserRepository()
	case "tasks":
		routes.AddTaskRoutes(router)
		repository.TaskRepositoryServices = repository.NewTaskRepository()
	default:
		repository.UserRepositoryServices = repository.NewUserRepository()
		repository.TaskRepositoryServices = repository.NewTaskRepository()
		routes.AddUserRoutes(router)
		routes.AddTaskRoutes(router)
	}

	//Message Broker
	ctx := context.Background()
	rabbitmq.Start(ctx)

	server := &http.Server{
		Addr:    httpPort,
		Handler: router,
	}

	// start API server
	go func() {
		fmt.Printf("%s API listening on port %s\n", cases.Title(language.AmericanEnglish).String(os.Getenv("SERVICE")), httpPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("%s: %s", "server forced to shutdown", err)
		os.Exit(1)
	}
}
