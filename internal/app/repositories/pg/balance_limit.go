package pg

import (
	"math"

	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

type balanceLimit struct{}

var BalanceLimit = balanceLimit{}

func (b balanceLimit) Create(tx *gorm.DB, accountID uint) error {
	balance := &types.BalanceLimit{
		AccountID: accountID,
		MaxNegBal: viper.GetFloat64("transaction.maxNegBal"),
		MaxPosBal: viper.GetFloat64("transaction.maxPosBal"),
	}
	err := tx.Create(balance).Error
	if err != nil {
		return e.Wrap(err, "pg.BalanceLimit.Create failed")
	}
	return nil
}

func (b balanceLimit) FindByAccountID(
	accountID uint,
) (*types.BalanceLimit, error) {
	balance := new(types.BalanceLimit)
	err := db.Where("account_id = ?", accountID).First(balance).Error
	if err != nil {
		return nil, e.Wrap(err, "pg.BalanceLimit.FindByAccountID failed")
	}
	return balance, nil
}

func (b balanceLimit) Update(
	id uint,
	maxPosBal float64,
	maxNegBal float64,
) error {
	if math.Abs(maxNegBal) == 0 {
		maxNegBal = 0
	} else {
		maxNegBal = math.Abs(maxNegBal)
	}

	err := db.
		Model(&types.BalanceLimit{}).
		Where("account_id = ?", id).
		Updates(map[string]interface{}{
			"max_pos_bal": math.Abs(maxPosBal),
			"max_neg_bal": maxNegBal,
		}).Error
	if err != nil {
		return e.Wrap(err, "pg.BalanceLimit.Update failed")
	}
	return nil
}
