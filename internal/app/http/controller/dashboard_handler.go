package controller

import (
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/helper"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/ic3network/mccs-alpha/internal/pkg/template"
	"go.uber.org/zap"
)

type dashBoardHandler struct {
	once *sync.Once
}

var DashBoardHandler = newDashBoardHandler()

func newDashBoardHandler() *dashBoardHandler {
	return &dashBoardHandler{
		once: new(sync.Once),
	}
}

func (d *dashBoardHandler) RegisterRoutes(
	public *mux.Router,
	private *mux.Router,
	adminPublic *mux.Router,
	adminPrivate *mux.Router,
) {
	d.once.Do(func() {
		private.Path("/").HandlerFunc(d.dashboardPage()).Methods("GET")
	})
}

func (d *dashBoardHandler) dashboardPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("dashboard")
	type response struct {
		User          *types.User
		Business      *types.Business
		MatchedOffers map[string][]string
		MatchedWants  map[string][]string
		Balance       float64
	}
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := UserHandler.FindByID(r.Header.Get("userID"))

		if err != nil {
			l.Logger.Error("DashboardPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}

		business, err := service.Business.FindByID(user.CompanyID)
		if err != nil {
			l.Logger.Error("DashboardPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}

		lastLoginDate := time.Time{}
		if user.ShowRecentMatchedTags {
			lastLoginDate = user.LastLoginDate
		}

		matchedOffers, err := service.Tag.MatchOffers(
			helper.GetTagNames(business.Offers),
			lastLoginDate,
		)
		if err != nil {
			l.Logger.Error("DashboardPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		matchedWants, err := service.Tag.MatchWants(
			helper.GetTagNames(business.Wants),
			lastLoginDate,
		)
		if err != nil {
			l.Logger.Error("DashboardPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}

		res := response{
			User:          user,
			Business:      business,
			MatchedOffers: matchedOffers,
			MatchedWants:  matchedWants,
		}

		// Get the account balance.
		account, err := service.Account.FindByBusinessID(user.CompanyID.Hex())
		if err != nil {
			l.Logger.Error("DashboardPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		res.Balance = account.Balance

		t.Render(w, r, res, nil)
	}
}
