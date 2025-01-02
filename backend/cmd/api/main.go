package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const port = 8080

type application struct {
	DSN          string
	Domain       string
	DB           repository.DatabaseRepo
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	APIKey       string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found. Using system environment variables.")
	}

	app := application{
		DSN:          mustGetEnv("DSN"),
		JWTSecret:    mustGetEnv("JWT_SECRET"),
		JWTIssuer:    mustGetEnv("JWT_ISSUER"),
		JWTAudience:  mustGetEnv("JWT_AUDIENCE"),
		CookieDomain: mustGetEnv("COOKIE_DOMAIN"),
		Domain:       mustGetEnv("DOMAIN"),
		APIKey:       mustGetEnv("API_KEY"),
	}

	// connect to the database
	conn, err := app.connectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	app.DB = &dbrepo.PostgresDBRepo{DB: conn}
	defer app.DB.Connection().Close()

	// configure authentication
	app.auth = Auth{
		Issuer:       app.JWTIssuer,
		Audience:     app.JWTAudience,
		Secret:       app.JWTSecret,
		TokenExpiry:  time.Minute * 15,
		RefreshExpiry: time.Hour * 24,
		CookiePath:   "/",
		CookieName:   "app_refresh_token",
		CookieDomain: app.CookieDomain,
	}

	log.Println("Starting application on port", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Error: Required environment variable %s is not set", key)
	}
	return value
}
