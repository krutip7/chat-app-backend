package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	var app = Application{
		domain:  "localhost",
		port:    8080,
		version: "1.0.0",
	}

	log.Printf("Starting Server on http://%s:%d", app.domain, app.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", app.port), app.routes())
	if err != nil {
		log.Fatal(err)
	}

}
