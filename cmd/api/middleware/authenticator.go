package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/krutip7/chat-app-server/cmd/api/auth"
	"github.com/krutip7/chat-app-server/cmd/api/utils"
)

type Authenticator struct {
	Auth *auth.Auth
}

func (authenticator *Authenticator) Authenticate(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		response.Header().Set("Vary", "Authorization")

		tokenString, err := getTokenFromAuthHeader(request.Header["Authorization"])
		if err == nil {
			_, err = authenticator.Auth.VerifyJWT(tokenString)
		}

		if err != nil {
			log.Println("middleware.Authenticator: authentication failed")
			utils.WriteJSONErrorResponse(response, err, http.StatusUnauthorized)
			return
		}

		log.Println("middleware.Authenticator: authentication successful")
		handler.ServeHTTP(response, request)
	})

}

func getTokenFromAuthHeader(authHeader []string) (string, error) {

	if authHeader == nil || len(authHeader) < 1 {
		return "", errors.New("auth header not found")
	}

	headerValues := strings.Split(authHeader[0], " ")
	if len(headerValues) != 2 || headerValues[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	return headerValues[1], nil
}
