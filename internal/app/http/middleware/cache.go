package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NoCache() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "no-cache, must-revalidate, no-store") // HTTP 1.1.
			w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
			w.Header().Set("Expires", "0")                                         // Proxies.
			next.ServeHTTP(w, r)
		})
	}
}
