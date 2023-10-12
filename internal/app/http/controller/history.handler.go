package controller

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/global/constant"
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/ic3network/mccs-alpha/internal/pkg/template"
	"github.com/ic3network/mccs-alpha/internal/pkg/util"
	"go.uber.org/zap"
)

type historyHandler struct {
	once *sync.Once
}

var HistoryHandler = newHistoryHandler()

func newHistoryHandler() *historyHandler {
	return &historyHandler{
		once: new(sync.Once),
	}
}

func (h *historyHandler) RegisterRoutes(
	public *mux.Router,
	private *mux.Router,
	adminPublic *mux.Router,
	adminPrivate *mux.Router,
) {
	h.once.Do(func() {
		private.Path("/history").HandlerFunc(h.historyPage()).Methods("GET")
		private.Path("/history/search").
			HandlerFunc(h.searchHistory()).
			Methods("GET")
	})
}

func (h *historyHandler) historyPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("history")
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow access to History screens for users with trading-accepted status
		business, _ := BusinessHandler.FindByUserID(r.Header.Get("userID"))
		if business.Status != constant.Trading.Accepted {
			http.Redirect(w, r, "/", http.StatusFound)
		}
		t.Render(w, r, nil, nil)
	}
}

func (h *historyHandler) searchHistory() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("history")
	type formData struct {
		DateFrom string
		DateTo   string
		Page     int
	}
	type response struct {
		FormData     formData
		TotalPages   int
		Balance      float64
		Transactions []*types.Transaction
	}
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		page, err := strconv.Atoi(q.Get("page"))
		if err != nil {
			l.Logger.Error(
				"controller.History.HistoryPage failed",
				zap.Error(err),
			)
			t.Error(w, r, nil, err)
			return
		}

		f := formData{
			DateFrom: q.Get("date-from"),
			DateTo:   q.Get("date-to"),
			Page:     page,
		}
		res := response{FormData: f}

		user, err := UserHandler.FindByID(r.Header.Get("userID"))
		if err != nil {
			l.Logger.Error(
				"controller.History.HistoryPage failed",
				zap.Error(err),
			)
			t.Error(w, r, nil, err)
			return
		}

		// Get the account balance.
		account, err := service.Account.FindByBusinessID(user.CompanyID.Hex())
		if err != nil {
			l.Logger.Error(
				"controller.History.HistoryPage failed",
				zap.Error(err),
			)
			t.Error(w, r, nil, err)
			return
		}
		res.Balance = account.Balance

		// Get the recent transactions.
		transactions, totalPages, err := service.Transaction.FindInRange(
			account.ID,
			util.ParseTime(f.DateFrom),
			util.ParseTime(f.DateTo),
			page,
		)
		if err != nil {
			l.Logger.Error(
				"controller.History.HistoryPage failed",
				zap.Error(err),
			)
			t.Error(w, r, nil, err)
			return
		}
		res.Transactions = transactions
		res.TotalPages = totalPages

		t.Render(w, r, res, nil)
	}
}
