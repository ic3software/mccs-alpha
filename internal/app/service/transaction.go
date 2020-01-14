package service

import (
	"fmt"
	"math"
	"time"

	"github.com/ic3network/mccs-alpha/internal/app/repositories/pg"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
)

type transaction struct{}

// Transaction services.
var Transaction = &transaction{}

func (t *transaction) maxBalanceCanBeTransferred(a *types.Account, kind string) (float64, error) {
	if kind == "positive" {
		maxPosBal, err := BalanceLimit.GetMaxPosBalance(a.ID)
		if err != nil {
			return 0, err
		}
		if a.Balance >= 0 {
			return maxPosBal - a.Balance, nil
		}
		return math.Abs(a.Balance) + maxPosBal, nil
	}
	maxNegBal, err := BalanceLimit.GetMaxNegBalance(a.ID)
	if err != nil {
		return 0, err
	}
	if a.Balance >= 0 {
		return a.Balance + maxNegBal, nil
	}
	return maxNegBal - math.Abs(a.Balance), nil
}

func (t *transaction) Propose(
	proposerID,

	fromID,
	fromEmail,
	fromBusinessName,

	toID,
	toEmail,
	toBusinessName string,

	amount float64,
	description string,
) (*types.Transaction, error) {
	// Get the Account IDs using MongoIDs.
	proposer, err := pg.Account.FindByBusinessID(proposerID)
	if err != nil {
		return nil, e.Wrap(err, "service.Transaction.Propose")
	}
	from, err := pg.Account.FindByBusinessID(fromID)
	if err != nil {
		return nil, e.Wrap(err, "service.Transaction.Propose")
	}
	to, err := pg.Account.FindByBusinessID(toID)
	if err != nil {
		return nil, e.Wrap(err, "service.Transaction.Propose")
	}

	// Check the account balance.
	exceed, err := BalanceLimit.IsExceedLimit(from.ID, from.Balance-amount)
	if err != nil {
		return nil, e.Wrap(err, "service.Transaction.Propose")
	}
	if exceed {
		amount, err := t.maxBalanceCanBeTransferred(from, "negative")
		if err != nil {
			return nil, e.Wrap(err, "service.Transaction.Propose")
		}
		return nil, e.CustomMessage("Sender will exceed its credit limit." + " The maximum amount that can be sent is: " + fmt.Sprintf("%.2f", amount))
	}
	exceed, err = BalanceLimit.IsExceedLimit(to.ID, to.Balance+amount)
	if err != nil {
		return nil, e.Wrap(err, "service.Transaction.Propose")
	}
	if exceed {
		amount, err := t.maxBalanceCanBeTransferred(to, "positive")
		if err != nil {
			return nil, e.Wrap(err, "service.Transaction.Propose")
		}
		return nil, e.CustomMessage("Receiver will exceed its maximum balance limit." + " The maximum amount that can be received is: " + fmt.Sprintf("%.2f", amount))
	}

	transaction, err := pg.Transaction.Propose(
		proposer.ID,
		from.ID,
		fromEmail,
		fromBusinessName,
		to.ID,
		toEmail,
		toBusinessName,
		amount,
		description,
	)
	if err != nil {
		return nil, e.Wrap(err, "service.Transaction.Propose")
	}
	return transaction, nil
}

func (t *transaction) Find(transactionID uint) (*types.Transaction, error) {
	transaction, err := pg.Transaction.Find(transactionID)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (t *transaction) FindPendings(accountID uint) ([]*types.Transaction, error) {
	transactions, err := pg.Transaction.FindPendings(accountID)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (t *transaction) Cancel(transactionID uint, reason string) error {
	err := pg.Transaction.Cancel(transactionID, reason)
	if err != nil {
		return err
	}
	return nil
}

func (t *transaction) Accept(
	transactionID uint,
	fromID uint,
	toID uint,
	amount float64,
) error {
	err := pg.Transaction.Accept(
		transactionID,
		fromID,
		toID,
		amount,
	)
	if err != nil {
		return e.Wrap(err, "service.Account.MakeTransfer failed")
	}
	return nil
}

func (t *transaction) FindRecent(accountID uint) ([]*types.Transaction, error) {
	transactions, err := pg.Transaction.FindRecent(accountID)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (t *transaction) FindInRange(accountID uint, dateFrom time.Time, dateTo time.Time, page int) ([]*types.Transaction, int, error) {
	transactions, totalPages, err := pg.Transaction.FindInRange(accountID, dateFrom, dateTo, page)
	if err != nil {
		return nil, 0, err
	}
	return transactions, totalPages, nil
}
