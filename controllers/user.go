package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/hugohenrick/gtasks/database"
	"github.com/hugohenrick/gtasks/models"
	"github.com/hugohenrick/gtasks/utils"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBind(&user); err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v: %v", utils.InvalidJsonProvided, err))
		return
	}

	hashPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashPassword

	database.DB.Create(&user)
	utils.SendJSONResponse(c, http.StatusOK, user)
}

func GetUsers(c *gin.Context) {
	var users []models.User

	database.DB.Find(&users)
	utils.SendJSONResponse(c, http.StatusOK, users)
}

//LoginUser : Generates JWT Token for validated user
func LoginUser(c *gin.Context) {
	var user models.UserLogin
	var hmacSampleSecret []byte

	if err := c.ShouldBindJSON(&user); err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v: %v", utils.InvalidJsonProvided, err))
		return
	}

	var dbUser models.User
	email := user.Email
	password := user.Password

	err := database.DB.Where("email = ?", email).First(&dbUser).Error
	if err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v: %v", utils.UserNotFound, err))
		return
	}

	hashErr := utils.CheckPasswordHash(password, dbUser.Password)
	if hashErr != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v: %v", utils.UserInvalidCredentials, err))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": dbUser,
		"exp":  time.Now().Add(time.Minute * 30).Unix(),
	})

	tokenString, err := token.SignedString(hmacSampleSecret)
	if err != nil {
		utils.SendJSONError(c, http.StatusBadRequest, fmt.Errorf("%v: %v", utils.UserFailedGetToken, err))
		return
	}

	utils.SendJSONResponse(c, http.StatusCreated, gin.H{
		"message": "Token generated sucessfully",
		"token":   tokenString,
	})
}
