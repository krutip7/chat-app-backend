package main

import (
	"net/http"

	"github.com/krutip7/chat-app-server/cmd/api/middleware"
	"github.com/krutip7/chat-app-server/cmd/api/sockethub"
)

func (app *Application) routes() http.Handler {

	baseRouter := configureBaseRouter(app)

	baseMW := middleware.MiddlewareStack(
		middleware.LogHTTPExchange,
		middleware.EnableCORS,
	)

	authenticator := middleware.Authenticator{Auth: app.auth}
	authMW := middleware.MiddlewareStack(authenticator.Authenticate)

	internalRouter := configureInternalRouter(app)
	baseRouter.Handle("/", authMW(internalRouter))

	return baseMW(baseRouter)
}

func configureBaseRouter(app *Application) *http.ServeMux {
	router := http.NewServeMux()

	socketHandler := sockethub.NewWebSocketHandler(app.auth)

	router.HandleFunc("/ws", socketHandler.HandleConnection)

	router.HandleFunc("/user", app.GetUser)

	router.HandleFunc("POST /authenticate", app.Authenticate)

	router.HandleFunc("/refresh", app.RefreshCookie)

	router.HandleFunc("/logout", app.Logout)

	router.HandleFunc("/health-check", app.HealthCheck)

	return router
}

func configureInternalRouter(app *Application) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /users", app.GetUsers)

	return router
}
