package util

import "regexp"

var emailRe = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// IsValidEmail checks if email is valid.
func IsValidEmail(email string) bool {
	if email == "" || !emailRe.MatchString(email) || len(email) > 100 {
		return false
	}
	return true
}
