package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/krutip7/chat-app-server/internals/repository"
)

type Application struct {
	domain  string
	port    int
	version string
	dsn     string
	db      *sqlx.DB
	repo    Repository
}

type Repository struct {
	userRepo repository.UserRepository
}
