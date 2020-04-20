package controller

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/global/constant"
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/email"
	"github.com/ic3network/mccs-alpha/internal/pkg/helper"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/ic3network/mccs-alpha/internal/pkg/recaptcha"
	"github.com/ic3network/mccs-alpha/internal/pkg/template"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type tradingHandler struct {
	once *sync.Once
}

var TradingHandler = newTradingHandler()

func newTradingHandler() *tradingHandler {
	return &tradingHandler{
		once: new(sync.Once),
	}
}

func (th *tradingHandler) RegisterRoutes(
	public *mux.Router,
	private *mux.Router,
	adminPublic *mux.Router,
	adminPrivate *mux.Router,
) {
	th.once.Do(func() {
		private.Path("/member-signup").HandlerFunc(th.signupPage()).Methods("GET")
		private.Path("/member-signup").HandlerFunc(th.signup()).Methods("POST")

		private.Path("/api/is-trading-member").HandlerFunc(th.isMember()).Methods("GET")
	})
}

func (th *tradingHandler) signupPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("member-signup")
	return func(w http.ResponseWriter, r *http.Request) {
		// Check if the business has the status of "accepted".
		business, err := BusinessHandler.FindByUserID(r.Header.Get("userID"))
		/*if err != nil || business.Status != constant.Business.Accepted {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}*/
		user, err := UserHandler.FindByID(r.Header.Get("userID"))
		if err != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		data := &types.TradingRegisterData{
			BusinessName:       business.BusinessName,
			IncType:            business.IncType,
			CompanyNumber:      business.CompanyNumber,
			BusinessPhone:      business.BusinessPhone,
			Website:            business.Website,
			Turnover:           business.Turnover,
			Description:        business.Description,
			LocationAddress:    business.LocationAddress,
			LocationCity:       business.LocationCity,
			LocationRegion:     business.LocationRegion,
			LocationPostalCode: business.LocationPostalCode,
			LocationCountry:    business.LocationCountry,
			FirstName:          user.FirstName,
			LastName:           user.LastName,
			Telephone:          user.Telephone,
		}
		data.RecaptchaSitekey = viper.GetString("recaptcha.site_key")
		t.Render(w, r, data, nil)
	}
}

func (th *tradingHandler) signup() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("member-signup")
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		// Validate user inputs.
		data := helper.Trading.GetRegisterData(r)
		data.RecaptchaSitekey = viper.GetString("recaptcha.site_key")
		errorMessages := data.Validate()
		if viper.GetString("env") == "production" {
			isValid := recaptcha.Verify(*r)
			if !isValid {
				errorMessages = append(errorMessages, recaptcha.Error()...)
			}
		}
		if len(errorMessages) > 0 {
			l.Logger.Info("TradingHandler.Signup failed", zap.Strings("input invalid", errorMessages))
			t.Render(w, r, data, errorMessages)
			return
		}

		// Update business collection.
		business, err := BusinessHandler.FindByUserID(r.Header.Get("userID"))
		if err != nil {
			l.Logger.Info("TradingHandler.Signup failed", zap.Error(err))
			t.Error(w, r, data, err)
			return
		}
		err = service.Trading.UpdateBusiness(business.ID, data)
		if err != nil {
			l.Logger.Info("TradingHandler.Signup failed", zap.Error(err))
			t.Error(w, r, data, err)
			return
		}

		// Update user collection.
		user, err := UserHandler.FindByID(r.Header.Get("userID"))
		if err != nil {
			l.Logger.Info("TradingHandler.Signup failed", zap.Error(err))
			t.Error(w, r, data, err)
			return
		}
		err = service.Trading.UpdateUser(user.ID, data)
		if err != nil {
			l.Logger.Info("TradingHandler.Signup failed", zap.Error(err))
			t.Error(w, r, data, err)
			return
		}

		// Send thank you email to the User's email address.
		go func() {
			err := email.SendThankYouEmail(data.FirstName, data.LastName, user.Email)
			if err != nil {
				l.Logger.Error("email.SendThankYouEmail failed", zap.Error(err))
			}
		}()
		// Send the to the OCN Admin email address.
		go func() {
			err := email.SendNewMemberSignupEmail(data.BusinessName, user.Email)
			if err != nil {
				l.Logger.Error("email.SendNewMemberSignupEmail failed", zap.Error(err))
			}
		}()

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (th *tradingHandler) isMember() func(http.ResponseWriter, *http.Request) {
	type response struct {
		IsMember bool
	}
	return func(w http.ResponseWriter, r *http.Request) {
		business, err := BusinessHandler.FindByUserID(r.Header.Get("userID"))
		if err != nil {
			l.Logger.Error("TradingHandler.IsMember failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res := response{IsMember: business.Status == constant.Trading.Accepted}
		js, err := json.Marshal(res)
		if err != nil {
			l.Logger.Error("TradingHandler.IsMember failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
