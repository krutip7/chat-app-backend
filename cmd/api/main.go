package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/krutip7/chat-app-server/cmd/api/middleware"
)

func main() {

	// Init Application Config
	var app = Application{
		domain:  "localhost",
		port:    8080,
		version: "1.0.0",
	}

	app.auth = &Auth{
		Issuer:               app.domain,
		Audience:             app.domain,
		AuthTokenValidity:    15 * time.Minute,
		RefreshTokenValidity: 24 * time.Hour,
		CookieName:           "__Host-refresh_token",
	}

	// Init Command Flags
	flag.IntVar(&app.port, "port", 8080, "Application Server port")
	flag.StringVar(&app.dsn, "dsn", os.Getenv("POSTGRES_DSN"), "Postgres DB Connection String")
	flag.StringVar(&app.jwtSecret, "secret", os.Getenv("JWT_SECRET"), "JWT Signing Secret Key")
	flag.Parse()

	app.auth.SigningKey = []byte(app.jwtSecret)

	// Init Database Connection
	app.InitDBConnection()
	app.InitUserRepository()
	defer app.db.Close()

	// Init Server
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", app.port),
		Handler: middleware.EnableCORS(app.routes()),
	}

	log.Printf("Starting Server on http://%s:%d/", app.domain, app.port)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
