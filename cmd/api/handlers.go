package main

import (
	"net/http"

	"github.com/krutip7/chat-app-server/cmd/api/models"
	"github.com/krutip7/chat-app-server/cmd/api/utils"
)

func (app *Application) HealthCheck(response http.ResponseWriter, request *http.Request) {
	payload := models.HealthStatus{
		Status:  "success",
		Message: "Chat App Backend Server is up and running",
		Version: app.version,
	}

	utils.WriteJSONResponse(response, payload)
}

func (app *Application) Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/health-check", http.StatusPermanentRedirect)
}
