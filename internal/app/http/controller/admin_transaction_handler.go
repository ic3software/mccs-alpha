package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/global/constant"
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/flash"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/ic3network/mccs-alpha/internal/pkg/log"
	"github.com/ic3network/mccs-alpha/internal/pkg/template"
	"github.com/ic3network/mccs-alpha/internal/pkg/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type adminTransactionHandler struct {
	once *sync.Once
}

var AdminTransactionHandler = newAdminTransactionHandler()

func newAdminTransactionHandler() *adminTransactionHandler {
	return &adminTransactionHandler{
		once: new(sync.Once),
	}
}

func (tr *adminTransactionHandler) RegisterRoutes(
	public *mux.Router,
	private *mux.Router,
	adminPublic *mux.Router,
	adminPrivate *mux.Router,
) {
	tr.once.Do(func() {
		adminPrivate.Path("/transaction").HandlerFunc(tr.transactionPage()).Methods("GET")
		adminPrivate.Path("/transaction").HandlerFunc(tr.transaction()).Methods("POST")

		adminPrivate.Path("/api/pendingTransactions").HandlerFunc(tr.pendingTransactions()).Methods("GET")
		adminPrivate.Path("/api/cancelTransaction").HandlerFunc(tr.cancelTransaction()).Methods("POST")
	})
}

func (tr *adminTransactionHandler) transactionPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/transaction")
	type formData struct {
		FromEmail   string
		ToEmail     string
		Amount      float64
		Description string
	}
	type response struct {
		FormData      formData
		CurBalance    float64
		MaxNegBalance float64
	}
	return func(w http.ResponseWriter, r *http.Request) {
		t.Render(w, r, nil, nil)
	}
}

func (tr *adminTransactionHandler) transaction() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("admin/transaction")
	type formData struct {
		FromEmail   string
		ToEmail     string
		Amount      float64
		Description string
	}
	type response struct {
		FormData      formData
		CurBalance    float64
		MaxNegBalance float64
	}
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		f := formData{
			FromEmail:   r.FormValue("from_email_address"),
			ToEmail:     r.FormValue("to_email_address"),
			Description: r.FormValue("description"),
		}
		res := response{FormData: f}

		// Validate user inputs.
		errorMessages := []string{}
		if !util.IsValidEmail(f.FromEmail) {
			errorMessages = append(errorMessages, "Please enter a valid sender email address.")
		}
		if !util.IsValidEmail(f.ToEmail) {
			errorMessages = append(errorMessages, "Please enter a valid recipient email address.")
		}
		amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
		// Amount should be positive value and with up to two decimal places.
		if err != nil || amount <= 0 || !util.IsDecimalValid(r.FormValue("amount")) {
			errorMessages = append(errorMessages, "Please enter a valid numeric amount to send with up to two decimal places.")
		}
		res.FormData.Amount = amount
		if len(errorMessages) > 0 {
			t.Render(w, r, res, errorMessages)
			return
		}
		f.Amount = amount

		from, err := BusinessHandler.FindByEmail(f.FromEmail)
		if err != nil {
			l.Logger.Info("Transaction failed", zap.Error(err))
			t.Error(w, r, res, err)
			return
		}
		to, err := BusinessHandler.FindByEmail(f.ToEmail)
		if err != nil {
			l.Logger.Info("Transaction failed", zap.Error(err))
			t.Error(w, r, res, err)
			return
		}

		// Only allow transfers with accounts that also have "trading-accepted" status
		if from.Status != constant.Trading.Accepted {
			t.Render(w, r, res, []string{"Sender is not a trading member. You can only make transfers from businesses that have trading member status."})
			return
		} else if to.Status != constant.Trading.Accepted {
			t.Render(w, r, res, []string{"Receiver is not a trading member. You can only make transfers to businesses that have trading member status."})
			return
		}
		if f.FromEmail == f.ToEmail {
			t.Render(w, r, res, []string{"You cannot create a transaction from and to the same account."})
			return
		}

		err = service.AdminTransaction.Create(
			from.ID.Hex(),
			f.FromEmail,
			from.BusinessName,
			to.ID.Hex(),
			f.ToEmail,
			to.BusinessName,
			f.Amount,
			f.Description,
		)
		if err != nil {
			l.Logger.Info("Transaction failed", zap.Error(err))
			t.Error(w, r, res, err)
			return
		}

		go func() {
			objID, _ := primitive.ObjectIDFromHex(r.Header.Get("userID"))
			adminUser, err := service.AdminUser.FindByID(objID)
			if err != nil {
				l.Logger.Error("log.Admin.Transaction failed", zap.Error(err))
				return
			}
			err = service.UserAction.Log(
				log.Admin.Transfer(adminUser, f.FromEmail, f.ToEmail, f.Amount, f.Description),
			)
			if err != nil {
				l.Logger.Error("log.Admin.Transaction failed", zap.Error(err))
			}
		}()

		flash.Success(w, f.FromEmail+" has transferred "+fmt.Sprintf("%.2f", f.Amount)+" Credits to "+f.ToEmail)
		http.Redirect(w, r, "/admin/transaction", http.StatusFound)
	}
}

func (tr *adminTransactionHandler) pendingTransactions() func(http.ResponseWriter, *http.Request) {
	type response struct {
		Transactions []*types.Transaction
	}
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		user, err := UserHandler.FindByBusinessID(q.Get("business_id"))
		if err != nil {
			l.Logger.Error("AdminTransactionHandler.pendingTransactions failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		account, err := AccountHandler.FindByUserID(user.ID.Hex())
		if err != nil {
			l.Logger.Error("AdminTransactionHandler.pendingTransactions failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		transactions, err := service.Transaction.FindPendings(account.ID)
		if err != nil {
			l.Logger.Error("AdminTransactionHandler.pendingTransactions failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res := response{Transactions: transactions}
		js, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (tr *adminTransactionHandler) isInitiatedStatus(w http.ResponseWriter, t *types.Transaction) (bool, error) {
	type response struct {
		Error string `json:"error"`
	}

	if t.Status == constant.Transaction.Completed {
		js, err := json.Marshal(response{Error: "The transaction has already been completed."})
		if err != nil {
			return false, err
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return false, nil
	} else if t.Status == constant.Transaction.Cancelled {
		js, err := json.Marshal(response{Error: "The transaction has already been cancelled."})
		if err != nil {
			return false, err
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return false, nil
	}

	return true, nil
}

func (tr *adminTransactionHandler) cancelTransaction() func(http.ResponseWriter, *http.Request) {
	type request struct {
		TransactionID uint   `json:"id"`
		Reason        string `json:"reason"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			l.Logger.Error("AdminTransactionHandler.cancelTransaction failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		transaction, err := service.Transaction.Find(req.TransactionID)
		if err != nil {
			l.Logger.Error("AdminTransactionHandler.cancelTransaction failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		shouldContinue, err := tr.isInitiatedStatus(w, transaction)
		if err != nil {
			l.Logger.Error("AdminTransactionHandler.cancelTransaction failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !shouldContinue {
			return
		}

		err = service.Transaction.Cancel(req.TransactionID, req.Reason)
		if err != nil {
			l.Logger.Error("AdminTransactionHandler.cancelTransaction failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
