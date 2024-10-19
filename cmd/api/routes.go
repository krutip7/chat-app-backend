package main

import "net/http"

func (app *Application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/", app.Redirect)

	router.HandleFunc("/user", app.GetUser)

	router.HandleFunc("/authenticate", app.Authenticate)

	router.HandleFunc("/health-check", app.HealthCheck)

	return router
}
