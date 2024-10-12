package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	api "chat-app-server/cmd/api/models"
)

func (app *Application) HealthCheck(response http.ResponseWriter, request *http.Request) {
	responsePayload := api.HealthCheckResponse{
		Status:  "success",
		Message: "Chat App Backend Server is up and running",
		Version: app.version,
	}

	jsonResponse, err := json.Marshal(responsePayload)
	if err != nil {
		fmt.Println(err)
	}

	response.Header().Add("Content-Type", "application/json")

	response.Write(jsonResponse)
}

func (app *Application) Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/health-check", http.StatusPermanentRedirect)
}
