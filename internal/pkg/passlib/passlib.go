package passlib

import (
	"strconv"
	"unicode"
)

var (
	minLen     = 8
	hasLetter  = false
	hasNumber  = false
	hasSpecial = false
)

// Validate validates the given password.
func Validate(password string) []string {
	messages := make([]string, 0, 4)

	for _, ch := range password {
		switch {
		case unicode.IsLetter(ch):
			hasLetter = true
		case unicode.IsNumber(ch):
			hasNumber = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
		}
	}

	if len(password) < minLen {
		messages = append(messages, "Password must be at least "+strconv.Itoa(minLen)+" characters long.")
	} else if len(password) > 100 {
		messages = append(messages, "Password cannot exceed 100 characters.")
	}
	if !hasLetter {
		messages = append(messages, "Password must have at least one letter.")
	}
	if !hasNumber {
		messages = append(messages, "Password must have at least one numeric value.")
	}
	if !hasSpecial {
		messages = append(messages, "Password must have at least one special character.")
	}

	return messages
}
