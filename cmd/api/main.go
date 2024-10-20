package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/krutip7/chat-app-server/cmd/api/auth"
)

func main() {

	// Init Application Config
	var app = Application{
		version: "1.0.0",
	}

	app.auth = &auth.Auth{
		Issuer:               app.domain,
		Audience:             app.domain,
		AuthTokenValidity:    15 * time.Minute,
		RefreshTokenValidity: 24 * time.Hour,
		CookieName:           "refresh_token",
	}

	// Init Command Flags
	flag.IntVar(&app.port, "port", 8080, "Application Server port")
	flag.StringVar(&app.domain, "domain", os.Getenv("APPLICATION_DOMAIN"), "Application Server Domain")
	flag.StringVar(&app.dsn, "dsn", os.Getenv("POSTGRES_DSN"), "Postgres DB Connection String")
	flag.StringVar(&app.jwtSecret, "secret", os.Getenv("JWT_SECRET"), "JWT Signing Secret Key")
	flag.Parse()

	app.auth.SigningKey = []byte(app.jwtSecret)

	// Init Database Connection
	app.InitDBConnection()
	app.InitUserRepository()
	defer app.db.Close()

	// Init Server
	log.Printf("Starting Server on http://%s:%d/", app.domain, app.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", app.port), app.routes())
	if err != nil {
		log.Fatal(err)
	}

}
