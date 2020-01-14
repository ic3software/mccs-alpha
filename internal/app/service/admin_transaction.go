package service

import (
	"github.com/ic3network/mccs-alpha/internal/app/repositories/pg"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
)

type adminTransaction struct{}

var AdminTransaction = &adminTransaction{}

func (a *adminTransaction) Create(
	fromID,
	fromEmail,
	fromBusinessName,

	toID,
	toEmail,
	toBusinessName string,

	amount float64,
	description string,
) error {
	// Get the Account IDs using MongoIDs.
	from, err := pg.Account.FindByBusinessID(fromID)
	if err != nil {
		return e.Wrap(err, "service.Account.MakeTransfer failed")
	}
	to, err := pg.Account.FindByBusinessID(toID)
	if err != nil {
		return e.Wrap(err, "service.Account.MakeTransfer failed")
	}

	// Check the account balance.
	exceed, err := BalanceLimit.IsExceedLimit(from.ID, from.Balance-amount)
	if err != nil {
		return e.Wrap(err, "service.Account.MakeTransfer failed")
	}
	if exceed {
		return e.New(e.ExceedMaxNegBalance, "max negative exceed")
	}
	exceed, err = BalanceLimit.IsExceedLimit(to.ID, to.Balance+amount)
	if err != nil {
		return e.Wrap(err, "service.Account.MakeTransfer failed")
	}
	if exceed {
		return e.New(e.ExceedMaxPosBalance, "max positive exceed")
	}

	err = pg.Transaction.Create(
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
		return e.Wrap(err, "service.Account.MakeTransfer failed")
	}
	return nil
}
