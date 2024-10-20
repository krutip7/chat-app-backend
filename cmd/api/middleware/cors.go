package middleware

import (
	"net/http"
	"os"
)

func EnableCORS(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		response.Header().Set("Access-Control-Allow-Origin", os.Getenv("ALLOW_ORIGIN"))
		response.Header().Set("Access-Control-Allow-Credentials", "true")

		if request.Method == http.MethodOptions {
			response.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			response.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization, X-CSRF-Token")
			response.WriteHeader(http.StatusAccepted)
			return
		}

		handler.ServeHTTP(response, request)
	})
}
