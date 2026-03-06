package repositories

import (
	"Gin/initializers"
	"Gin/models"
	"context"
	"time"
)

type IUsersRepository interface {
	PostUser(user *models.User) error
	GetUsers() ([]models.User, error)
	GetUser(userID uint) (*models.User, error)
	UpdateUser(user *models.User) error
}

type usersRepository struct{}

func NewUsersRepository() IUsersRepository {
	return &usersRepository{}
}

func (r *usersRepository) PostUser(user *models.User) error {
	result := initializers.DB.Create(user)
	return result.Error
}

func (r *usersRepository) UpdateUser(user *models.User) error {
	result := initializers.DB.Save(user)
	return result.Error
}

func (r *usersRepository) GetUsers() ([]models.User, error) {
	context, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var users []models.User
	result := initializers.DB.WithContext(context).Preload("Company").Preload("Books").Preload("House", nil).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *usersRepository) GetUser(userID uint) (*models.User, error) {
	context, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var user models.User
	result := initializers.DB.WithContext(context).Preload("Company").Preload("Books", nil).Preload("House", nil).Where("id = ?", userID).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
