package types

import (
	"github.com/jinzhu/gorm"
)

type BalanceLimit struct {
	gorm.Model
	// `BalanceLimit` belongs to `Account`, `AccountID` is the foreign key
	Account   Account
	AccountID uint    `gorm:"not null;unique_index"`
	MaxNegBal float64 `gorm:"type:int;not null"`
	MaxPosBal float64 `gorm:"type:int;not null"`
}
