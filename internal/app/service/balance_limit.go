package service

import (
	"math"

	"github.com/ic3network/mccs-alpha/internal/app/repositories/pg"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
)

type balanceLimit struct{}

var BalanceLimit = balanceLimit{}

func (b balanceLimit) FindByAccountID(id uint) (*types.BalanceLimit, error) {
	record, err := pg.BalanceLimit.FindByAccountID(id)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (b balanceLimit) FindByBusinessID(id string) (*types.BalanceLimit, error) {
	account, err := Account.FindByBusinessID(id)
	if err != nil {
		return nil, err
	}
	record, err := pg.BalanceLimit.FindByAccountID(account.ID)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (b balanceLimit) GetMaxPosBalance(id uint) (float64, error) {
	balanceLimitRecord, err := pg.BalanceLimit.FindByAccountID(id)
	if err != nil {
		return 0, e.Wrap(err, "service.BalanceLimit.GetMaxPosBalance failed")
	}
	return balanceLimitRecord.MaxPosBal, nil
}

func (b balanceLimit) GetMaxNegBalance(id uint) (float64, error) {
	balanceLimitRecord, err := pg.BalanceLimit.FindByAccountID(id)
	if err != nil {
		return 0, e.Wrap(err, "service.BalanceLimit.GetMaxNegBalance failed")
	}
	return math.Abs(balanceLimitRecord.MaxNegBal), nil
}

// IsExceedLimit checks whether or not the account exceeds the max positive or max negative limit.
func (b balanceLimit) IsExceedLimit(id uint, balance float64) (bool, error) {
	balanceLimitRecord, err := pg.BalanceLimit.FindByAccountID(id)
	if err != nil {
		return false, e.Wrap(err, "service.BalanceLimit.FindByAccountID failed")
	}
	// MaxNegBal should be positive in the DB.
	if balance < -(math.Abs(balanceLimitRecord.MaxNegBal)) || balance > balanceLimitRecord.MaxPosBal {
		return true, nil
	}
	return false, nil
}

func (b balanceLimit) Update(id uint, maxPosBal float64, maxNegBal float64) error {
	err := pg.BalanceLimit.Update(id, maxPosBal, maxNegBal)
	if err != nil {
		return err
	}
	return nil
}
