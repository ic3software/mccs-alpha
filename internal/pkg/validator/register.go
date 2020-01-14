package validator

import (
	"strings"

	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/util"
)

// ValidateBusiness validates
// BusinessName, Offers and Wants
func ValidateBusiness(b *types.BusinessData) []string {
	errs := []string{}
	if b.BusinessName == "" {
		errs = append(errs, "Business name is missing.")
	} else if len(b.BusinessName) > 100 {
		errs = append(errs, "Business Name cannot exceed 100 characters.")
	}
	if b.Website != "" && !strings.HasPrefix(b.Website, "http://") && !strings.HasPrefix(b.Website, "https://") {
		errs = append(errs, "Website URL should start with http:// or https://.")
	} else if len(b.Website) > 100 {
		errs = append(errs, "Website URL cannot exceed 100 characters.")
	}
	errs = append(errs, validateTagsLimit(b)...)
	return errs
}

// ValidateUser validates
// FirstName, LastName, Email and Email
func ValidateUser(u *types.User) []string {
	errorMessages := []string{}
	u.Email = strings.ToLower(u.Email)
	if u.Email == "" {
		errorMessages = append(errorMessages, "Email is missing.")
	} else if !util.IsValidEmail(u.Email) {
		errorMessages = append(errorMessages, "Email is invalid.")
	} else if len(u.Email) > 100 {
		errorMessages = append(errorMessages, "Email cannot exceed 100 characters.")
	}
	return errorMessages
}
