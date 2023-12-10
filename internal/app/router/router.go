package router

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"restapi/internal/app/handlers"
	
	"restapi/internal/app/middlewares"
	appMiddleware "restapi/internal/app/middleware"
)

func NewRouter(dbConn *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.SentryMiddleware)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		userHandler := handlers.NewUserHandler(dbConn)

		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.CircuitBreaker("postgres"))
			r.Mount("/users", userHandler.Routes())
		})

		r.Group(func(r chi.Router) {
			r.Use(middlewares.SetMiddlewareAuthentication) 

			// Protected routes
			r.Get("/testGet", handlers.NewCustomHandler().TestGet)
			r.Post("/testPost", handlers.NewCustomHandler().TestPost)
			r.Post("/getCertificate", handlers.NewCustomHandler().GetCertificate)
			r.Post("/checkSignature", handlers.NewCustomHandler().CheckSignature)
			r.Post("/createSignature", handlers.NewCustomHandler().CreateSignature)

			// Add any other protected routes here
		})

		r.Options("/*", handlers.NewCustomHandler().PreFlight) // Middleware for OPTIONS request
        r.Get("/testGet", handlers.NewCustomHandler().TestGet)
        r.Post("/testPost", handlers.NewCustomHandler().TestPost)
        r.Post("/getCertificate", handlers.NewCustomHandler().GetCertificate)
        r.Post("/checkSignature", handlers.NewCustomHandler().CheckSignature)
        r.Post("/createSignature", handlers.NewCustomHandler().CreateSignature)
	})


	
	return r
}




