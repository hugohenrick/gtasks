package repository

import (
	"errors"

	"github.com/hugohenrick/gtasks/database"
	"github.com/hugohenrick/gtasks/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	FindUsers() ([]models.User, error)
	FindUserByEmail(email string) (models.User, error)
	CreateUser(User models.User) (models.User, error)
}

type UserRepository struct {
	Database *gorm.DB
}

var UserRepositoryServices IUserRepository

func NewUserRepository() IUserRepository {
	return &UserRepository{Database: database.DB}
}

func (t *UserRepository) FindUsers() ([]models.User, error) {
	var users []models.User

	result := database.DB.Find(&users)

	if result.RowsAffected == 0 {
		return []models.User{}, errors.New("user data not found")
	}

	return users, nil
}

func (t *UserRepository) FindUserByEmail(email string) (models.User, error) {
	var user models.User

	err := t.Database.Where("email = ?", email).First(&user).Error

	if err != nil {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func (t *UserRepository) CreateUser(user models.User) (models.User, error) {
	result := t.Database.Create(&user)

	if result.RowsAffected == 0 {
		return models.User{}, errors.New("user not created")
	}

	return user, nil
}
