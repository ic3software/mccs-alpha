package controller

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/ic3network/mccs-alpha/global/constant"
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/email"
	"github.com/ic3network/mccs-alpha/internal/pkg/flash"
	"github.com/ic3network/mccs-alpha/internal/pkg/jsonerror"
	"github.com/ic3network/mccs-alpha/internal/pkg/l"
	"github.com/ic3network/mccs-alpha/internal/pkg/log"
	"github.com/ic3network/mccs-alpha/internal/pkg/template"
	"github.com/unrolled/render"

	"github.com/ic3network/mccs-alpha/internal/pkg/util"
	"go.uber.org/zap"
)

type transactionHandler struct {
	once *sync.Once
}

// TransactionHandler handles transaction related logic.
var TransactionHandler = newTransactionHandler()

func newTransactionHandler() *transactionHandler {
	return &transactionHandler{
		once: new(sync.Once),
	}
}

func (tr *transactionHandler) RegisterRoutes(
	public *mux.Router,
	private *mux.Router,
	adminPublic *mux.Router,
	adminPrivate *mux.Router,
) {
	tr.once.Do(func() {
		private.Path("/transaction").
			HandlerFunc(tr.proposeTransactionPage()).
			Methods("GET")
		private.Path("/transaction").
			HandlerFunc(tr.proposeTransaction()).
			Methods("POST")
		private.Path("/transaction/cancelPropose").
			HandlerFunc(tr.cancelPropose()).
			Methods("GET")

		private.Path("/pending_transactions").
			HandlerFunc(tr.pendingTransactionsPage()).
			Methods("GET")

		private.Path("/api/accountBalance").
			HandlerFunc(tr.getBalance()).
			Methods("GET")
		private.Path("/api/pendingTransactions").
			HandlerFunc(tr.pendingTransactions()).
			Methods("GET")
		private.Path("/api/acceptTransaction").
			HandlerFunc(tr.acceptTransaction()).
			Methods("POST")
		private.Path("/api/cancelTransaction").
			HandlerFunc(tr.cancelTransaction()).
			Methods("POST")
		private.Path("/api/rejectTransaction").
			HandlerFunc(tr.rejectTransaction()).
			Methods("POST")
		private.Path("/api/recentTransactions").
			HandlerFunc(tr.recentTransactions()).
			Methods("GET")
	})
}

func (tr *transactionHandler) getMaxNegBal(
	r *http.Request,
	res *response,
) error {
	// Get the user max negative balance.
	account, err := AccountHandler.FindByUserID(r.Header.Get("userID"))
	if err != nil {
		return err
	}
	maxNegBalance, err := service.BalanceLimit.GetMaxNegBalance(account.ID)
	if err != nil {
		return err
	}
	res.CurBalance = account.Balance
	if res.CurBalance < 0 {
		res.MaxNegBalance = maxNegBalance - math.Abs(res.CurBalance)
	} else {
		res.MaxNegBalance = maxNegBalance
	}
	return nil
}

type formData struct {
	Type        string // "send" or "receive"
	Email       string
	Amount      float64
	Description string
}
type response struct {
	FormData      formData
	CurBalance    float64
	MaxNegBalance float64
}

func (tr *transactionHandler) cancelPropose() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		flash.Info(w, "No transfer has been initiated from your account.")
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (tr *transactionHandler) proposeTransactionPage() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("transaction")
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow access to Transfer screens for users with trading-accepted status
		business, _ := BusinessHandler.FindByUserID(r.Header.Get("userID"))
		if business.Status != constant.Trading.Accepted {
			http.Redirect(w, r, "/", http.StatusFound)
		}

		res := response{}
		err := tr.getMaxNegBal(r, &res)
		if err != nil {
			l.Logger.Info("Transfer failed", zap.Error(err))
			t.Error(w, r, res, err)
			return
		}

		t.Render(w, r, res, nil)
	}
}

type proposeInfo struct {
	FromID,
	FromEmail,
	FromBusinessName,
	FromStatus,
	ToID,
	ToEmail,
	ToBusinessName,
	ToStatus string
}

func (tr *transactionHandler) getProposeInfo(
	kind string,
	initiatorBusiness,
	receiverBusiness *types.Business,
	initiatorEmail,
	receiverEmail string,
) *proposeInfo {
	var fromID, fromEmail, fromBusinessName, fromStatus, toID, toEmail, toBusinessName, toStatus string
	if kind == "send" {
		fromID = initiatorBusiness.ID.Hex()
		fromEmail = initiatorEmail
		fromBusinessName = initiatorBusiness.BusinessName
		fromStatus = initiatorBusiness.Status

		toID = receiverBusiness.ID.Hex()
		toEmail = receiverEmail
		toBusinessName = receiverBusiness.BusinessName
		toStatus = receiverBusiness.Status
	} else if kind == "receive" {
		fromID = receiverBusiness.ID.Hex()
		fromEmail = receiverEmail
		fromBusinessName = receiverBusiness.BusinessName
		fromStatus = receiverBusiness.Status

		toID = initiatorBusiness.ID.Hex()
		toEmail = initiatorEmail
		toBusinessName = initiatorBusiness.BusinessName
		toStatus = initiatorBusiness.Status
	}
	return &proposeInfo{
		fromID,
		fromEmail,
		fromBusinessName,
		fromStatus,
		toID,
		toEmail,
		toBusinessName,
		toStatus,
	}
}

func (tr *transactionHandler) proposeTransaction() func(http.ResponseWriter, *http.Request) {
	t := template.NewView("transaction")
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		f := formData{
			Type:        r.FormValue("transfer_type"),
			Email:       r.FormValue("email_address"),
			Description: r.FormValue("description"),
		}

		res := response{FormData: f}
		err := tr.getMaxNegBal(r, &res)
		if err != nil {
			l.Logger.Info("Transfer failed", zap.Error(err))
			t.Error(w, r, res, err)
			return
		}

		// Validate the user inputs.
		errorMessages := []string{}
		if !util.IsValidEmail(f.Email) {
			errorMessages = append(
				errorMessages,
				"Please enter a valid email address.",
			)
		}
		amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
		// Amount should be positive value and with up to two decimal places.
		if err != nil || amount <= 0 ||
			!util.IsDecimalValid(r.FormValue("amount")) {
			errorMessages = append(
				errorMessages,
				"Please enter a valid numeric amount to send with up to two decimal places.",
			)
		}
		res.FormData.Amount = amount
		if len(errorMessages) > 0 {
			t.Render(w, r, res, errorMessages)
			return
		}
		f.Amount = amount

		initiator, err := UserHandler.FindByID(r.Header.Get("userID"))
		if err != nil {
			l.Logger.Info("Transfer failed", zap.Error(err))
			t.Error(w, r, res, err)
			return
		}
		// Decide the initiator and receiver.
		initiatorBusiness, err := service.Business.FindByID(initiator.CompanyID)
		if err != nil {
			l.Logger.Info("Transfer failed", zap.Error(err))
			t.Error(w, r, res, err)
			return
		}
		receiverBusiness, err := BusinessHandler.FindByEmail(f.Email)
		if err != nil {
			l.Logger.Info("Transfer failed", zap.Error(err))
			t.Error(w, r, res, err)
			return
		}
		proposeInfo := tr.getProposeInfo(
			f.Type,
			initiatorBusiness,
			receiverBusiness,
			initiator.Email,
			f.Email,
		)

		// Only allow transfers with accounts that also have "trading-accepted" status
		if (proposeInfo.FromStatus != constant.Trading.Accepted) ||
			(proposeInfo.ToStatus != constant.Trading.Accepted) {
			t.Render(
				w,
				r,
				res,
				[]string{
					"Recipient is not a trading member. You can only make transfers to other businesses that have trading member status.",
				},
			)
			return
		}

		// Check if the user is doing the transaction to himself.
		if proposeInfo.FromID == proposeInfo.ToID {
			t.Render(
				w,
				r,
				res,
				[]string{"You cannot create a transaction with yourself."},
			)
			return
		}

		transaction, err := service.Transaction.Propose(
			initiator.CompanyID.Hex(),
			proposeInfo.FromID,
			proposeInfo.FromEmail,
			proposeInfo.FromBusinessName,
			proposeInfo.ToID,
			proposeInfo.ToEmail,
			proposeInfo.ToBusinessName,
			f.Amount,
			f.Description,
		)
		if err != nil {
			l.Logger.Info("Proposed failed", zap.Error(err))
			t.Error(w, r, res, err)
			return
		}
		flash.Success(
			w,
			"You have proposed a transfer of "+fmt.Sprintf(
				"%.2f",
				f.Amount,
			)+" Credits with "+f.Email,
		)
		http.Redirect(w, r, "/#transactions", http.StatusFound)

		go func() {
			err := service.UserAction.Log(log.User.ProposeTransfer(
				initiator,
				proposeInfo.FromEmail,
				proposeInfo.ToEmail,
				f.Amount,
				f.Description,
			))
			if err != nil {
				l.Logger.Error("log.User.Transfer failed", zap.Error(err))
			}
		}()
		go func() {
			err := email.Transaction.Initiate(f.Type, transaction)
			if err != nil {
				l.Logger.Error(
					"email.Transaction.Initiate failed",
					zap.Error(err),
				)
			}
		}()
	}
}

func (tr *transactionHandler) getBalance() func(http.ResponseWriter, *http.Request) {
	type response struct {
		Balance float64
	}
	return func(w http.ResponseWriter, r *http.Request) {
		account, err := AccountHandler.FindByUserID(r.Header.Get("userID"))
		if err != nil {
			l.Logger.Error("TransferHandler.getBalance failed", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res := response{Balance: account.Balance}
		js, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (tr *transactionHandler) pendingTransactions() func(http.ResponseWriter, *http.Request) {
	type response struct {
		Transactions []*types.Transaction
	}
	return func(w http.ResponseWriter, r *http.Request) {
		account, err := AccountHandler.FindByUserID(r.Header.Get("userID"))
		if err != nil {
			l.Logger.Error(
				"TransferHandler.pendingTransactions failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		transactions, err := service.Transaction.FindPendings(account.ID)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.pendingTransactions failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res := response{Transactions: transactions}
		js, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (tr *transactionHandler) recentTransactions() func(http.ResponseWriter, *http.Request) {
	type response struct {
		Transactions []*types.Transaction
	}
	return func(w http.ResponseWriter, r *http.Request) {
		account, err := AccountHandler.FindByUserID(r.Header.Get("userID"))
		if err != nil {
			l.Logger.Error(
				"TransferHandler.recentTransactions failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		transactions, err := service.Transaction.FindRecent(account.ID)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.recentTransactions failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res := response{Transactions: transactions}
		js, err := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func (tr *transactionHandler) isInitiatedStatus(
	w http.ResponseWriter,
	t *types.Transaction,
) (bool, error) {
	type response struct {
		Error string `json:"error"`
	}

	if t.Status == constant.Transaction.Completed {
		js, err := json.Marshal(
			response{
				Error: "The transaction has already been completed by the counterparty.",
			},
		)
		if err != nil {
			return false, err
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return false, nil
	} else if t.Status == constant.Transaction.Cancelled {
		js, err := json.Marshal(response{Error: "The transaction has already been cancelled by the counterparty."})
		if err != nil {
			return false, err
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return false, nil
	}

	return true, nil
}

func (tr *transactionHandler) cancelTransaction() func(http.ResponseWriter, *http.Request) {
	type request struct {
		TransactionID uint   `json:"id"`
		Reason        string `json:"reason"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.cancelTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		account, err := AccountHandler.FindByUserID(r.Header.Get("userID"))
		if err != nil {
			l.Logger.Error(
				"TransferHandler.cancelTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		transaction, err := service.Transaction.Find(req.TransactionID)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.cancelTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		shouldContinue, err := tr.isInitiatedStatus(w, transaction)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.cancelTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !shouldContinue {
			return
		}

		if account.ID != transaction.InitiatedBy {
			l.Logger.Error(
				"TransferHandler.cancelTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = service.Transaction.Cancel(req.TransactionID, req.Reason)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.cancelTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		go func() {
			err := email.Transaction.Cancel(transaction, req.Reason)
			if err != nil {
				l.Logger.Error(
					"email.Transaction.Cancel failed",
					zap.Error(err),
				)
			}
		}()
	}
}

func (tr *transactionHandler) rejectTransaction() func(http.ResponseWriter, *http.Request) {
	type request struct {
		TransactionID uint   `json:"id"`
		Reason        string `json:"reason"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.rejectTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		transaction, err := service.Transaction.Find(req.TransactionID)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.cancelTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		shouldContinue, err := tr.isInitiatedStatus(w, transaction)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.cancelTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !shouldContinue {
			return
		}

		err = service.Transaction.Cancel(req.TransactionID, req.Reason)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.rejectTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		go func() {
			err := email.Transaction.Reject(transaction)
			if err != nil {
				l.Logger.Error(
					"email.Transaction.Reject failed",
					zap.Error(err),
				)
			}
		}()
	}
}

func (tr *transactionHandler) acceptTransaction() func(http.ResponseWriter, *http.Request) {
	type request struct {
		TransactionID uint `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.acceptTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		transaction, err := service.Transaction.Find(req.TransactionID)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.acceptTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		shouldContinue, err := tr.isInitiatedStatus(w, transaction)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.cancelTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !shouldContinue {
			return
		}

		from, err := service.Account.FindByID(transaction.FromID)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.acceptTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		to, err := service.Account.FindByID(transaction.ToID)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.acceptTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Check the account balance.
		exceed, err := service.BalanceLimit.IsExceedLimit(
			from.ID,
			from.Balance-transaction.Amount,
		)
		if err != nil {
			l.Logger.Info(
				"TransferHandler.acceptTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if exceed {
			reason := "The sender will exceed its credit limit so this transaction has been cancelled."
			err = service.Transaction.Cancel(req.TransactionID, reason)
			if err != nil {
				l.Logger.Error(
					"TransferHandler.acceptTransaction failed",
					zap.Error(err),
				)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err := jsonerror.New("1", reason)
			render.New().JSON(w, http.StatusInternalServerError, err.Render())
			go func() {
				err := email.Transaction.CancelBySystem(transaction, reason)
				if err != nil {
					l.Logger.Error(
						"email.Transaction.Cancel failed",
						zap.Error(err),
					)
				}
			}()
			return
		}
		exceed, err = service.BalanceLimit.IsExceedLimit(
			to.ID,
			to.Balance+transaction.Amount,
		)
		if err != nil {
			l.Logger.Info(
				"TransferHandler.acceptTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if exceed {
			reason := "The recipient will exceed its maximum positive balance threshold so this transaction has been cancelled."
			err = service.Transaction.Cancel(req.TransactionID, reason)
			if err != nil {
				l.Logger.Error(
					"TransferHandler.acceptTransaction failed",
					zap.Error(err),
				)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			err := jsonerror.New("2", reason)
			render.New().JSON(w, http.StatusInternalServerError, err.Render())
			go func() {
				err := email.Transaction.CancelBySystem(transaction, reason)
				if err != nil {
					l.Logger.Error(
						"email.Transaction.Cancel failed",
						zap.Error(err),
					)
				}
			}()
			return
		}

		err = service.Transaction.Accept(
			transaction.ID,
			from.ID,
			to.ID,
			transaction.Amount,
		)
		if err != nil {
			l.Logger.Error(
				"TransferHandler.acceptTransaction failed",
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

		go func() {
			err := email.Transaction.Accept(transaction)
			if err != nil {
				l.Logger.Error(
					"email.Transaction.Accept failed",
					zap.Error(err),
				)
			}
		}()
	}
}

// pendingTransactionsPage redirects the user to the dashboard (/) page after the user login.
func (tr *transactionHandler) pendingTransactionsPage() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/#transactions", http.StatusFound)
	}
}
