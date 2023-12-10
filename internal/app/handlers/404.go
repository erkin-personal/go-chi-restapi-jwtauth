package handlers

import (
    "net/http"
    "github.com/getsentry/sentry-go"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    // Capture the 404 error with Sentry
    sentry.CaptureMessage("404 - Page Not Found: " + r.URL.Path)

    // Respond with a 404 error
    http.NotFound(w, r)
}
