package middleware

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/internal/pkg/jwt"
)

func GetLoggedInUser() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("mccsToken")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			mccsToken := cookie.Value
			claims, err := jwt.NewJWTManager().ValidateToken(mccsToken)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			r.Header.Set("userID", claims.UserID)
			r.Header.Set("admin", strconv.FormatBool(claims.Admin))
			next.ServeHTTP(w, r)
		})
	}
}

func RequireUser() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := r.Header.Get("userID")
			// When user is not logged in.
			// 	1. If it's on the root page, redirect to find businesses page.
			// 	2. If it's on the other page, redirect to the targeting page after logging in.
			if userID == "" {
				if url.QueryEscape(r.URL.String()) == url.QueryEscape("/") {
					http.Redirect(w, r, "/businesses/search?page=1", http.StatusFound)
				} else {
					http.Redirect(w, r, "/login?redirect_login="+url.QueryEscape(r.URL.String()), http.StatusFound)
				}
				return
			}
			admin, _ := strconv.ParseBool(r.Header.Get("admin"))
			if admin == true {
				// Redirect to admin page if user is an admin.
				http.Redirect(w, r, "/admin", http.StatusFound)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAdmin() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			admin, err := strconv.ParseBool(r.Header.Get("admin"))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if admin != true {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
