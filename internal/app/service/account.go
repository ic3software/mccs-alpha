package service

import (
	"github.com/ic3network/mccs-alpha/internal/app/repositories/pg"
	"github.com/ic3network/mccs-alpha/internal/app/types"
)

type account struct{}

var Account = &account{}

func (a *account) Create(bID string) error {
	err := pg.Account.Create(bID)
	if err != nil {
		return err
	}
	return nil
}

func (a *account) FindByID(accountID uint) (*types.Account, error) {
	account, err := pg.Account.FindByID(accountID)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a *account) FindByBusinessID(businessID string) (*types.Account, error) {
	account, err := pg.Account.FindByBusinessID(businessID)
	if err != nil {
		return nil, err
	}
	return account, nil
}
