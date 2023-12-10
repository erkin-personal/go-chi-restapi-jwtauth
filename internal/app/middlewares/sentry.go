package middlewares

import (
    "net/http"
    "github.com/getsentry/sentry-go"
)

func SentryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        hub := sentry.GetHubFromContext(r.Context())
        if hub == nil {
            hub = sentry.CurrentHub().Clone()
            r = r.WithContext(sentry.SetHubOnContext(r.Context(), hub))
        }
        hub.Scope().SetRequest(r)

        defer func() {
            if err := recover(); err != nil {
                hub.Recover(err)
                // Optionally, write a generic error response to the client
            }
        }()

        next.ServeHTTP(w, r)
    })
}
