package models

import "github.com/krutip7/chat-app-server/internals/models"

type PostLogin struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}
