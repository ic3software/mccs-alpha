package helper

import (
	"net/http"
	"strconv"

	"github.com/ic3network/mccs-alpha/internal/app/types"
)

type trading struct{}

var Trading = &trading{}

func (t *trading) GetRegisterData(r *http.Request) *types.TradingRegisterData {
	turnover, _ := strconv.Atoi(r.FormValue("turnover"))
	return &types.TradingRegisterData{
		BusinessName:       r.FormValue("business_name"),        // 100 chars
		IncType:            r.FormValue("inc_type"),             // 25 chars
		CompanyNumber:      r.FormValue("company_number"),       // 20 chars
		BusinessPhone:      r.FormValue("business_phone"),       // 25 chars
		Website:            r.FormValue("website"),              // 100 chars
		Turnover:           turnover,                            // 20 chars
		Description:        r.FormValue("description"),          // 500 chars
		LocationAddress:    r.FormValue("location_address"),     // 255 chars
		LocationCity:       r.FormValue("location_city"),        // 50 chars
		LocationRegion:     r.FormValue("location_region"),      // 50 chars
		LocationPostalCode: r.FormValue("location_postal_code"), // 10 chars
		LocationCountry:    r.FormValue("location_country"),     // 50 chars
		FirstName:          r.FormValue("first_name"),           // 100 chars
		LastName:           r.FormValue("last_name"),            // 100 chars
		Telephone:          r.FormValue("telephone"),            // 25 chars
		Authorised:         r.FormValue("authorised"),
	}
}

func (t *trading) GetUpdateData(r *http.Request) *types.TradingUpdateData {
	turnover, _ := strconv.Atoi(r.FormValue("turnover"))
	return &types.TradingUpdateData{
		BusinessName:       r.FormValue("business_name"),        // 100 chars
		IncType:            r.FormValue("inc_type"),             // 25 chars
		CompanyNumber:      r.FormValue("company_number"),       // 20 chars
		BusinessPhone:      r.FormValue("business_phone"),       // 25 chars
		Website:            r.FormValue("website"),              // 100 chars
		Turnover:           turnover,                            // 20 chars
		Description:        r.FormValue("description"),          // 500 chars
		LocationAddress:    r.FormValue("location_address"),     // 255 chars
		LocationCity:       r.FormValue("location_city"),        // 50 chars
		LocationRegion:     r.FormValue("location_region"),      // 50 chars
		LocationPostalCode: r.FormValue("location_postal_code"), // 10 chars
		LocationCountry:    r.FormValue("location_country"),     // 50 chars
		FirstName:          r.FormValue("first_name"),           // 100 chars
		LastName:           r.FormValue("last_name"),            // 100 chars
		Telephone:          r.FormValue("telephone"),            // 25 chars
	}
}
