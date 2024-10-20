package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/krutip7/chat-app-server/cmd/api/auth"
	"github.com/krutip7/chat-app-server/internals/repository"

	database "github.com/krutip7/chat-app-server/internals/repository/postgres"
)

type Application struct {
	domain    string
	port      int
	version   string
	dsn       string
	db        *sqlx.DB
	repo      Repository
	auth      *auth.Auth
	jwtSecret string
}

type Repository struct {
	userRepo repository.UserRepository
}

func (app *Application) InitDBConnection() {
	db, err := database.Connect(app.dsn)

	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
		return
	}

	app.db = db
	log.Printf("Connected to database successfully")
}

func (app *Application) InitUserRepository() {
	userRepo := database.UserRepository{
		DB: app.db,
	}

	app.repo = Repository{
		userRepo: &userRepo,
	}
}
