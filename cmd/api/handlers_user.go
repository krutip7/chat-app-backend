package main

import (
	"log"
	"net/http"

	"github.com/krutip7/chat-app-server/cmd/api/utils"
)

func (app *Application) GetUsers(response http.ResponseWriter, request *http.Request) {

	users, err := app.repo.userRepo.GetAllUsers()
	if err != nil {
		log.Println(err)
	}

	utils.WriteJSONResponse(response, users)

}
