package services

import (
	"WBA/models"
)

type UserService interface {
	CheckUser(string) (*models.User, error)
}
