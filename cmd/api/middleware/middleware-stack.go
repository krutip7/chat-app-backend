package middleware

import "net/http"

/*
Usage: Allows to stack middlewares as array instead of nesting handlers

	appMiddleware := MiddlewareStack(LoggingMiddleware, CORSMiddleware)
	router = appMiddleware(router)
*/
func MiddlewareStack(stack ...MiddlewareHandler) MiddlewareHandler {

	return func(handler http.Handler) http.Handler {

		n := len(stack)

		nextHandler := handler
		for i := range n {
			nextHandler = stack[n-i-1](nextHandler)
		}

		return nextHandler
	}
}
