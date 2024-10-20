package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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
	if err != nil || user == nil {
		utils.WriteJSONErrorResponse(response, InvalidUserCredentials, http.StatusBadRequest)
		return
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

	data := models.PostAuthenticationResponse{
		Token: tokenPair.AuthToken,
		User:  *user,
	}

	refreshTokenCookie := app.auth.GetRefreshTokenCookie(tokenPair.RefreshToken)
	http.SetCookie(response, refreshTokenCookie)

	utils.WriteJSONResponse(response, data)
}

func (app *Application) RefreshCookie(response http.ResponseWriter, request *http.Request) {

	var claims jwt.Claims
	var refreshToken string
	for _, cookie := range request.Cookies() {
		if cookie.Name == app.auth.CookieName {
			refreshToken = cookie.Value
		}
	}

	UnauthorizedUser := errors.New("unauthorized user")

	oldToken, err := jwt.ParseWithClaims(refreshToken, claims, func(_ *jwt.Token) (any, error) { return app.jwtSecret, nil })
	if err != nil {
		utils.WriteJSONErrorResponse(response, UnauthorizedUser, http.StatusUnauthorized)
		return
	}

	sub, err := oldToken.Claims.GetSubject()
	if err != nil {
		utils.WriteJSONErrorResponse(response, UnauthorizedUser, http.StatusUnauthorized)
		return
	}

	userId, err := strconv.Atoi(sub)
	if err != nil {
		utils.WriteJSONErrorResponse(response, UnauthorizedUser, http.StatusUnauthorized)
		return
	}

	user, err := app.repo.userRepo.GetUserById(userId)
	if err != nil {
		utils.WriteJSONErrorResponse(response, UnauthorizedUser, http.StatusUnauthorized)
		return
	}

	tokenPair, err := app.auth.GenerateJWTToken(user)
	if err != nil {
		utils.WriteJSONErrorResponse(response, err, http.StatusInternalServerError)
		return
	}

	data := models.PostAuthenticationResponse{
		Token: tokenPair.AuthToken,
		User:  *user,
	}

	refreshCookie := app.auth.GetRefreshTokenCookie(tokenPair.RefreshToken)
	http.SetCookie(response, refreshCookie)

	utils.WriteJSONResponse(response, data)
}

func (app *Application) Logout(response http.ResponseWriter, request *http.Request) {
	expiredRefreshCookie := app.auth.ClearRefreshTokenCookie()
	http.SetCookie(response, expiredRefreshCookie)
	utils.WriteJSONResponse(response, nil)
}
