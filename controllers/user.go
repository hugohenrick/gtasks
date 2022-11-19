package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
