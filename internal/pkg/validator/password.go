package validator

import (
	"github.com/ic3network/mccs-alpha/internal/pkg/passlib"
)

func ValidatePassword(password string, confirmPassword string) []string {
	errorMessages := []string{}

	if password == "" {
		errorMessages = append(errorMessages, "Please enter a password.")
	} else if password != confirmPassword {
		errorMessages = append(errorMessages, "Password and confirmation password do not match.")
	} else {
		errorMessages = append(errorMessages, passlib.Validate(password)...)
	}

	return errorMessages
}

func validateUpdatePassword(currentPass string, newPass string, confirmPass string) []string {
	errorMessages := []string{}

	if currentPass == "" && newPass == "" && confirmPass == "" {
		return errorMessages
	}

	if currentPass == "" {
		errorMessages = append(errorMessages, "Please enter your current password.")
	} else if newPass != confirmPass {
		errorMessages = append(errorMessages, "New password and confirmation password do not match.")
	} else {
		errorMessages = append(errorMessages, passlib.Validate(newPass)...)
	}

	return errorMessages
}
