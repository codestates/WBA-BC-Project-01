package services

import (
	"WBA/models"
)

type UserService interface {
	CreateUser(*models.User) error
	GetUser(string) (*models.User, error)
}
