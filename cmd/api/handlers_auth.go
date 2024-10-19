package main

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/krutip7/chat-app-server/cmd/api/models"
	"github.com/krutip7/chat-app-server/cmd/api/utils"
)

func (app *Application) Authenticate(response http.ResponseWriter, request *http.Request) {

	payload := models.LoginRequest{}

	err := utils.ReadJSONRequest(response, request, &payload)
	if err != nil {
		utils.WriteJSONErrorResponse(response, err, http.StatusBadRequest)
		return
	}

	InvalidUserCredentials := errors.New("invalid credentials")

	user, err := app.repo.userRepo.GetUserByEmail(strings.ToLower(payload.Email))
	if err != nil {
		utils.WriteJSONErrorResponse(response, InvalidUserCredentials, http.StatusBadRequest)
	}

	valid, err := user.VerifyPassword(payload.Password)
	if err != nil {
		utils.WriteJSONErrorResponse(response, err, http.StatusInternalServerError)
		return
	} else if !valid {
		utils.WriteJSONErrorResponse(response, InvalidUserCredentials, http.StatusBadRequest)
		return
	}

	tokenPair, err := app.auth.GenerateJWTToken(user)
	if err != nil {
		log.Println(err)
		utils.WriteJSONErrorResponse(response, errors.New("token generation failed"), http.StatusInternalServerError)
		return
	}

	data := models.PostLoginResponse{
		Token: tokenPair.AuthToken,
		User:  *user,
	}

	refreshTokenCookie := app.auth.GetRefreshTokenCookie(tokenPair.RefreshToken)
	http.SetCookie(response, refreshTokenCookie)

	utils.WriteJSONResponse(response, data)
}

func (app *Application) Logout(response http.ResponseWriter, request *http.Request) {
	expiredRefreshCookie := app.auth.ClearRefreshTokenCookie()
	http.SetCookie(response, expiredRefreshCookie)
	utils.WriteJSONResponse(response, nil)
}
