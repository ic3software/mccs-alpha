package controller

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/ic3network/mccs-alpha/internal/pkg/template"
	"github.com/ic3network/mccs-alpha/internal/pkg/util"
	"go.uber.org/zap"
)

type adminHistoryHandler struct {
	once *sync.Once
}

var AdminHistoryHandler = newAdminHistoryHandler()

func newAdminHistoryHandler() *adminHistoryHandler {
	return &adminHistoryHandler{
		once: new(sync.Once),
	}
}

func (h *adminHistoryHandler) RegisterRoutes(
	public *mux.Router,
	private *mux.Router,
	adminPublic *mux.Router,
	adminPrivate *mux.Router,
) {
	h.once.Do(func() {
		adminPrivate.Path("/history/{id}").HandlerFunc(h.historyPage()).Methods("GET")
	})
}

func (h *adminHistoryHandler) historyPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/history")
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
		Email        string
		BusinessID   string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bID := vars["id"]
		q := r.URL.Query()

		page, err := strconv.Atoi(q.Get("page"))
		if err != nil {
			l.Logger.Error("controller.AdminHistory.HistoryPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}

		f := formData{
			DateFrom: q.Get("date-from"),
			DateTo:   q.Get("date-to"),
			Page:     page,
		}
		user, err := UserHandler.FindByBusinessID(bID)
		if err != nil {
			l.Logger.Error("controller.AdminHistory.HistoryPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		res := response{FormData: f, BusinessID: bID, Email: user.Email}

		// Get the account balance.
		account, err := service.Account.FindByBusinessID(bID)
		if err != nil {
			l.Logger.Error("controller.AdminHistory.HistoryPage failed", zap.Error(err))
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
			l.Logger.Error("controller.AdminHistory.HistoryPage failed", zap.Error(err))
			t.Error(w, r, nil, err)
			return
		}
		res.Transactions = transactions
		res.TotalPages = totalPages

		t.Render(w, r, res, nil)
	}
}
