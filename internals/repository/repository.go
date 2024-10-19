package repository

import (
	"github.com/krutip7/chat-app-server/internals/models"
)

type UserRepository interface {
	GetUserById(userId int) (*models.User, error)
	GetUserByEmail(userId string) (*models.User, error)
}
