package cookie

import "net/http"

// CreateCookie creates the default cookie.
func CreateCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     "mccsToken",
		Value:    value,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	}
}

// ResetCookie resets the default cookie.
func ResetCookie() *http.Cookie {
	return &http.Cookie{
		Name:     "mccsToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
}
