package main

import "net/http"

func (app *Application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/user", app.GetUser)

	router.HandleFunc("POST /authenticate", app.Authenticate)

	router.HandleFunc("/logout", app.Logout)

	router.HandleFunc("/health-check", app.HealthCheck)

	return router
}
