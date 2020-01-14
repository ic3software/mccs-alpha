package pg

import (
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
)

type account struct{}

var Account = &account{}

func (a *account) Create(bID string) error {
	tx := db.Begin()

	account := &types.Account{BusinessID: bID, Balance: 0}
	err := tx.Create(account).Error
	if err != nil {
		tx.Rollback()
		return e.Wrap(err, "pg.Account.Create failed")
	}
	err = BalanceLimit.Create(tx, account.ID)
	if err != nil {
		tx.Rollback()
		return e.Wrap(err, "pg.Account.Create failed")
	}

	return tx.Commit().Error
}

func (a *account) FindByID(accountID uint) (*types.Account, error) {
	var result types.Account
	err := db.Raw(`
	SELECT A.id, A.business_id, A.balance
	FROM accounts AS A
	WHERE A.id = ?
	LIMIT 1
	`, accountID).Scan(&result).Error
	if err != nil {
		return nil, e.Wrap(err, "pg.Account.FindByID")
	}
	return &result, nil
}

func (a *account) FindByBusinessID(businessID string) (*types.Account, error) {
	account := new(types.Account)
	err := db.Where("business_id = ?", businessID).First(account).Error
	if err != nil {
		return nil, e.New(e.UserNotFound, "user not found")
	}
	return account, nil
}
