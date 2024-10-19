package models

import "github.com/krutip7/chat-app-server/internals/models"

type PostAuthenticationResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
