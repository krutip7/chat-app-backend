package middleware

import "net/http"

type MiddlewareHandler func(http.Handler) http.Handler

type Middleware interface {
	Add(MiddlewareHandler) Middleware

	Intercept(http.Handler) http.Handler
}

/*
Usage: Allows to chain middlewares instead of nesting handlers

	authMiddleware := NewMiddleware(baseMiddleware).Add(JWTVerificationMiddleware)
	authRouter = authMiddleware.Intercept(authRouter)
*/
func NewMiddleware(middlewares ...MiddlewareHandler) Middleware {

	return &MiddlewareChain{
		stack: middlewares,
	}
}

type MiddlewareChain struct {
	stack []MiddlewareHandler
}

func (mw *MiddlewareChain) Add(middleware MiddlewareHandler) Middleware {

	mw.stack = append(mw.stack, middleware)

	return mw
}

func (mw *MiddlewareChain) Intercept(handler http.Handler) http.Handler {

	n := len(mw.stack)

	nextHandler := handler
	for i := range n {
		nextHandler = mw.stack[n-i-1](nextHandler)
	}

	return nextHandler
}
