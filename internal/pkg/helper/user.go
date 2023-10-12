package helper

import (
	"net/http"
	"strconv"

	"github.com/ic3network/mccs-alpha/internal/app/types"
)

func GetRegisterData(r *http.Request) *types.RegisterData {
	return &types.RegisterData{
		Business:        GetBusiness(r),
		User:            GetUser(r),
		ConfirmPassword: r.FormValue("confirm_password"),
		ConfirmEmail:    r.FormValue("confirm_email"),
		Terms:           r.FormValue("terms"),
	}
}

func GetUpdateData(r *http.Request) *types.UpdateAccountData {
	return &types.UpdateAccountData{
		Business:        GetBusiness(r),
		User:            GetUser(r),
		Balance:         &types.BalanceLimit{},
		CurrentPassword: r.FormValue("current_password"),
		ConfirmPassword: r.FormValue("confirm_password"),
	}
}

func GetBusiness(r *http.Request) *types.BusinessData {
	turnover, _ := strconv.Atoi(r.FormValue("turnover"))
	b := &types.BusinessData{
		BusinessName:  r.FormValue("business_name"),  // 100 chars
		IncType:       r.FormValue("inc_type"),       // 25 chars
		CompanyNumber: r.FormValue("company_number"), // 20 chars
		BusinessPhone: r.FormValue("business_phone"), // 25 chars
		Website:       r.FormValue("website"),        // 100 chars
		Turnover:      turnover,                      // 20 chars
		Offers: GetTags(
			r.FormValue("offers"),
		), // 500 chars (max 50 chars per tag)
		Wants: GetTags(
			r.FormValue("wants"),
		), // 500 chars (max 50 chars per tag)
		Description:        r.FormValue("description"),          // 500 chars
		LocationAddress:    r.FormValue("location_address"),     // 255 chars
		LocationCity:       r.FormValue("location_city"),        // 50 chars
		LocationRegion:     r.FormValue("location_region"),      // 50 chars
		LocationPostalCode: r.FormValue("location_postal_code"), // 10 chars
		LocationCountry:    r.FormValue("location_country"),     // 50 chars
		AdminTags:          getAdminTags(r.FormValue("adminTags")),
		Status:             r.FormValue("status"),
	}
	return b
}

func GetUser(r *http.Request) *types.User {
	return &types.User{
		FirstName:         r.FormValue("first_name"),   // 100 chars
		LastName:          r.FormValue("last_name"),    // 100 chars
		Email:             r.FormValue("email"),        // 100 chars
		Telephone:         r.FormValue("telephone"),    // 25 chars
		Password:          r.FormValue("new_password"), // 100 chars
		DailyNotification: r.FormValue("daily_notification") == "true",
	}
}
