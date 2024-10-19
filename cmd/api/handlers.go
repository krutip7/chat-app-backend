package main

import (
	"errors"
	"log"
	"net/http"
	"strings"

	api "github.com/krutip7/chat-app-server/cmd/api/models"
	"github.com/krutip7/chat-app-server/cmd/api/utils"
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

	payload := api.LoginRequest{}

	err := utils.ReadJSONRequest(response, request, payload)
	if err != nil {
		utils.WriteJSONErrorResponse(response, err, http.StatusBadRequest)
		return
	}

	InvalidUserCredentials := errors.New("invalid credentials")

	user, err := app.repo.userRepo.GetUserByEmail(strings.ToLower(payload.Email))
	if err != nil {
		utils.WriteJSONErrorResponse(response, InvalidUserCredentials, http.StatusBadRequest)
	}

	

	tokenPair, err := app.auth.GenerateJWTToken(user)
	if err != nil {
		log.Println(err)
		utils.WriteJSONErrorResponse(response, errors.New("token generation failed"), http.StatusInternalServerError)
		return
	}

	data := api.PostLoginResponse{
		Token: tokenPair.AuthToken,
		User:  *user,
	}

	refreshTokenCookie := app.auth.GetRefreshTokenCookie(tokenPair.RefreshToken)
	http.SetCookie(response, refreshTokenCookie)

	utils.WriteJSONResponse(response, data)
}
