package main

import (
	"errors"
	"log"
	"net/http"

	api "github.com/krutip7/chat-app-server/cmd/api/models"
	"github.com/krutip7/chat-app-server/cmd/api/utils"
	"github.com/krutip7/chat-app-server/internals/models"
)

func (app *Application) GetUser(response http.ResponseWriter, request *http.Request) {

	user, err := app.repo.userRepo.GetUserById(1)
	if err != nil {
		log.Println(err)
	}

	utils.WriteJSONResponse(response, user)

}

func (app *Application) HealthCheck(response http.ResponseWriter, request *http.Request) {
	utils.WriteJSONResponse(response, nil)
}

func (app *Application) Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/health-check", http.StatusPermanentRedirect)
}

func (app *Application) Authenticate(response http.ResponseWriter, request *http.Request) {

	// TODO: Verify user credentials and fetch details from DB
	user := models.User{
		Id:        "1",
		FirstName: "John",
		LastName:  "Doe",
		Username:  "johndoe",
	}

	tokenPair, err := app.auth.GenerateJWTToken(&user)
	if err != nil {
		log.Println(err)
		utils.WriteJSONErrorResponse(response, errors.New("token generation failed"), http.StatusInternalServerError)
		return
	}

	data := api.PostLogin{
		Token: tokenPair.AuthToken,
		User:  user,
	}
	utils.WriteJSONResponse(response, data)
}
