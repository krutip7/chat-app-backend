package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	// Init Application Config
	var app = Application{
		domain:  "localhost",
		port:    8080,
		version: "1.0.0",
	}

	// Init Command Flags
	flag.StringVar(&app.dsn, "dsn", os.Getenv("POSTGRES_DSN"), "Postgres DB Connection String")
	flag.Parse()

	// Init Database Connection
	app.InitDBConnection()
	app.InitUserRepository()
	defer app.db.Close()

	// Init Server
	log.Printf("Starting Server on http://%s:%d", app.domain, app.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", app.port), app.routes())
	if err != nil {
		log.Fatal(err)
	}

}
