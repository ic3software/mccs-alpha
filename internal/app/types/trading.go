package types

import "strings"

type TradingRegisterData struct {
	// Business
	BusinessName       string
	IncType            string
	CompanyNumber      string
	BusinessPhone      string
	Website            string
	Turnover           int
	Description        string
	LocationAddress    string
	LocationCity       string
	LocationRegion     string
	LocationPostalCode string
	LocationCountry    string
	// User
	FirstName string
	LastName  string
	Telephone string
	// Terms
	Authorised string
	// Recaptcha
	RecaptchaSitekey string
}

type TradingUpdateData struct {
	// Business
	BusinessName       string
	IncType            string
	CompanyNumber      string
	BusinessPhone      string
	Website            string
	Turnover           int
	Description        string
	LocationAddress    string
	LocationCity       string
	LocationRegion     string
	LocationPostalCode string
	LocationCountry    string
	// User
	FirstName string
	LastName  string
	Telephone string
	// Recaptcha
	RecaptchaSitekey string
}

func (t *TradingRegisterData) Validate() []string {
	errs := []string{}

	if t.BusinessName == "" {
		errs = append(errs, "Business name is missing.")
	} else if len(t.BusinessName) > 100 {
		errs = append(errs, "Business name cannot exceed 100 characters.")
	}
	if t.IncType == "" {
		errs = append(errs, "Incorporation type is missing.")
	} else if len(t.IncType) > 25 {
		errs = append(errs, "Incorporation type cannot exceed 25 characters.")
	}
	if len(t.CompanyNumber) > 20 {
		errs = append(errs, "Company number cannot exceed 20 characters")
	}
	if t.Description == "" {
		errs = append(errs, "Business description is missing.")
	} else if len(t.Description) > 500 {
		errs = append(errs, "Business description cannot exceed 500 characters.")
	}
	if t.Website != "" && !strings.HasPrefix(t.Website, "http://") && !strings.HasPrefix(t.Website, "https://") {
		errs = append(errs, "Website URL should start with http:// or https://.")
	} else if len(t.Website) > 100 {
		errs = append(errs, "Website URL cannot exceed 100 characters.")
	}
	if len(t.BusinessPhone) > 25 {
		errs = append(errs, "Business phone cannot exceed 25 characters.")
	}
	if t.LocationAddress == "" {
		errs = append(errs, "Business address is missing.")
	} else if len(t.LocationAddress) > 255 {
		errs = append(errs, "Business address cannot exceed 255 characters.")
	}
	if t.LocationCity == "" {
		errs = append(errs, "City/town is missing.")
	} else if len(t.LocationCity) > 50 {
		errs = append(errs, "City/town cannot exceed 50 characters.")
	}
	if t.LocationRegion == "" {
		errs = append(errs, "Region/county is missing.")
	} else if len(t.LocationRegion) > 50 {
		errs = append(errs, "Region/county cannot exceed 50 characters.")
	}
	if t.LocationPostalCode == "" {
		errs = append(errs, "Postcode is missing.")
	} else if len(t.LocationPostalCode) > 10 {
		errs = append(errs, "Postcode cannot exceed 10 characters.")
	}
	if t.LocationCountry == "" {
		errs = append(errs, "Country is missing.")
	} else if len(t.LocationCountry) > 50 {
		errs = append(errs, "Country cannot exceed 50 characters.")
	}

	if t.FirstName == "" {
		errs = append(errs, "First name is missing.")
	} else if len(t.FirstName) > 100 {
		errs = append(errs, "First name cannot exceed 100 characters.")
	}
	if t.LastName == "" {
		errs = append(errs, "Last name is missing.")
	} else if len(t.LastName) > 100 {
		errs = append(errs, "Last name cannot exceed 100 characters.")
	}
	if t.Telephone == "" {
		errs = append(errs, "Telephone is missing.")
	} else if len(t.Telephone) > 25 {
		errs = append(errs, "Telephone cannot exceed 25 characters.")
	}
	if t.Authorised != "on" {
		errs = append(errs, "Please confirm you have read and agree to the Membership Agreement on behalf of your business.")
	}

	return errs
}

func (t *TradingUpdateData) Validate() []string {
	errs := []string{}

	if t.BusinessName == "" {
		errs = append(errs, "Business name is missing.")
	} else if len(t.BusinessName) > 100 {
		errs = append(errs, "Business name cannot exceed 100 characters.")
	}
	if t.IncType == "" {
		errs = append(errs, "Incorporation type is missing.")
	} else if len(t.IncType) > 25 {
		errs = append(errs, "Incorporation type cannot exceed 25 characters.")
	}
	if len(t.CompanyNumber) > 20 {
		errs = append(errs, "Company number cannot exceed 20 characters.")
	}
	if t.Description == "" {
		errs = append(errs, "Business description is missing.")
	} else if len(t.Description) > 500 {
		errs = append(errs, "Business description cannot exceed 500 characters.")
	}
	if t.Website != "" && !strings.HasPrefix(t.Website, "http://") && !strings.HasPrefix(t.Website, "https://") {
		errs = append(errs, "Website URL should start with http:// or https://.")
	} else if len(t.Website) > 100 {
		errs = append(errs, "Website URL cannot exceed 100 characters.")
	}
	if len(t.BusinessPhone) > 25 {
		errs = append(errs, "Business phone cannot exceed 25 characters.")
	}
	if t.LocationAddress == "" {
		errs = append(errs, "Business address is missing.")
	} else if len(t.LocationAddress) > 255 {
		errs = append(errs, "Business address cannot exceed 255 characters.")
	}
	if t.LocationCity == "" {
		errs = append(errs, "City/town is missing.")
	} else if len(t.LocationCity) > 50 {
		errs = append(errs, "City/town cannot exceed 50 characters.")
	}
	if t.LocationRegion == "" {
		errs = append(errs, "Region/county is missing.")
	} else if len(t.LocationRegion) > 50 {
		errs = append(errs, "Region/county cannot exceed 50 characters.")
	}
	if t.LocationPostalCode == "" {
		errs = append(errs, "Postcode is missing.")
	} else if len(t.LocationPostalCode) > 10 {
		errs = append(errs, "Postcode cannot exceed 10 characters.")
	}
	if t.LocationCountry == "" {
		errs = append(errs, "Country is missing.")
	} else if len(t.LocationCountry) > 50 {
		errs = append(errs, "Country cannot exceed 50 characters.")
	}

	if t.FirstName == "" {
		errs = append(errs, "First name is missing.")
	} else if len(t.FirstName) > 100 {
		errs = append(errs, "First name cannot exceed 100 characters.")
	}
	if t.LastName == "" {
		errs = append(errs, "Last name is missing.")
	} else if len(t.LastName) > 100 {
		errs = append(errs, "Last name cannot exceed 100 characters.")
	}
	if t.Telephone == "" {
		errs = append(errs, "Telephone is missing.")
	} else if len(t.Telephone) > 25 {
		errs = append(errs, "Telephone cannot exceed 25 characters.")
	}

	return errs
}
