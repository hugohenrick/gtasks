package database

import (
	"fmt"
	"log"
	"os"

	"github.com/hugohenrick/gtasks/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Conn() {
	USER := os.Getenv("DB_USER")
	PASS := os.Getenv("DB_PASSWORD")
	HOST := os.Getenv("DB_HOST")
	DBNAME := os.Getenv("DB_NAME")

	URL := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", USER, PASS, HOST, DBNAME)

	fmt.Println(URL)

	DB, err = gorm.Open(mysql.Open(URL))
	if err != nil {
		log.Panic("Failed to connect to database!")
	}

	fmt.Println("Database connection established")

	DB.AutoMigrate(&models.Task{}, &models.User{})
}
