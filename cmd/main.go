package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"

	"restapi/internal/app/router"
	"restapi/internal/infrastructure/db"
	"restapi/internal/infrastructure/hystrix"

	"github.com/getsentry/sentry-go"
	
)

func init() {
    // Get the current working directory
    cwd, err := os.Getwd()
    if err != nil {
        log.Fatal("Error getting current working directory:", err)
    }

    // Load the .env file from the current working directory
    err = godotenv.Load(filepath.Join(cwd, "../.env"))
    if err != nil {
        log.Fatal("Error loading .env file:", err)
    }

}

func main() {

	err := sentry.Init(sentry.ClientOptions{
        Dsn: "https://dc76d816bcdae49efdb05c0d460486c9@sentry.its.zone/4",
    })
    if err != nil {
        log.Fatalf("sentry.Init: %s", err)
    }

	
    defer sentry.Flush(2 * time.Second)



	dbConn, err := db.NewConnection()
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v", err)
	}
	defer dbConn.Close()

	hystrix.ConfigureCircuitBreaker()
	r := router.NewRouter(dbConn)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

	

	
}





