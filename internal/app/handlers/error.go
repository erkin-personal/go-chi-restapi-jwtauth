package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/getsentry/sentry-go"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
    // Intentionally trigger an error
    err := errors.New("intentional error triggered from /error endpoint")
    sentry.CaptureException(err)

    // Optionally, you can write a response to the client
    http.Error(w, err.Error(), http.StatusInternalServerError)

	log.Printf("INTENTIONAL /error route Captured by SENTRY")
}
