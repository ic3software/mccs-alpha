package validator

import (
	"github.com/ic3network/mccs-alpha/internal/app/types"
)

func Register(d *types.RegisterData) []string {
	errorMessages := []string{}
	errorMessages = append(errorMessages, ValidateBusiness(d.Business)...)
	errorMessages = append(errorMessages, ValidateUser(d.User)...)
	errorMessages = append(
		errorMessages,
		ValidatePassword(d.User.Password, d.ConfirmPassword)...)

	if d.User.Email != d.ConfirmEmail {
		errorMessages = append(
			errorMessages,
			"The email addresses you entered do not match.",
		)
	}
	if d.Terms != "on" {
		errorMessages = append(
			errorMessages,
			"Please confirm you accept to have your business listed in OCN's directory.",
		)
	}
	return errorMessages
}

func Account(d *types.UpdateAccountData) []string {
	errorMessages := []string{}
	errorMessages = append(errorMessages, ValidateBusiness(d.Business)...)
	errorMessages = append(errorMessages, ValidateUser(d.User)...)
	errorMessages = append(
		errorMessages,
		validateUpdatePassword(
			d.CurrentPassword,
			d.User.Password,
			d.ConfirmPassword,
		)...)
	return errorMessages
}

func UpdateBusiness(b *types.BusinessData) []string {
	errorMessages := []string{}
	errorMessages = append(errorMessages, ValidateBusiness(b)...)
	return errorMessages
}
