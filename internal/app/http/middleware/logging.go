package middleware

import (
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/internal/pkg/ip"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
)

// Logging middleware logs messages.
func Logging() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			defer func() {
				uri := r.RequestURI
				// Skip for the health check and static requests.
				if uri == "/health" || uri == "/ram" || uri == "/cpu" ||
					uri == "/disk" ||
					strings.HasPrefix(uri, "/static") {
					return
				}
				elapse := time.Now().Sub(startTime)
				l.Logger.Info("request",
					zap.String("ip", ip.FromRequest(r)),
					zap.String("method", r.Method),
					zap.String("uri", uri),
					zap.String("userAgent", r.UserAgent()),
					zap.Duration("responseTime", elapse))
			}()

			next.ServeHTTP(w, r)
		})
	}
}
