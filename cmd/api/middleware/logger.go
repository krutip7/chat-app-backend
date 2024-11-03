package middleware

import (
	"bufio"
	"log"
	"net"
	"net/http"
)

type ResponseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (w *ResponseWriterWrapper) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}
func (w *ResponseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}

func LogHTTPExchange(handler http.Handler) http.Handler {

	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {

		wrapper := &ResponseWriterWrapper{
			response, http.StatusOK,
		}

		handler.ServeHTTP(wrapper, request)

		log.Println("middleware.HTTPLogger: ", wrapper.statusCode, request.Method, request.URL)
	})
}
