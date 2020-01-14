package middleware

import (
	"net/http"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"go.uber.org/zap"
)

func Recover() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					buf := make([]byte, 1024)
					runtime.Stack(buf, false)
					l.Logger.Error("recover, error", zap.Any("err", err), zap.ByteString("method", buf))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
