package middlewares

import (
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-chi/chi/middleware"
)

func CircuitBreaker(name string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			err := hystrix.Do(name, func() error {
				ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
				next.ServeHTTP(ww, r)
				return nil
			}, nil)

			if err != nil {
				http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			}
		}
		return http.HandlerFunc(fn)
	}
}
