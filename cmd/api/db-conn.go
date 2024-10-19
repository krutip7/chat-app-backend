package main

import (
	"log"

	pg "github.com/krutip7/chat-app-server/internals/repository/postgres"
)

func (app *Application) InitDBConnection() {
	db, err := pg.InitPostgresDB(app.dsn)

	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
		return
	}

	app.db = db
	log.Printf("Connected to database successfully")
}

func (app *Application) InitUserRepository() {
	userRepo := pg.UserRepository{
		DB: app.db,
	}

	app.repo = Repository{
		userRepo: &userRepo,
	}
}
