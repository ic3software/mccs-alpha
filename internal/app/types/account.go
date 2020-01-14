package types

import (
	"github.com/jinzhu/gorm"
)

type Account struct {
	gorm.Model
	// Account has many postings, AccountID is the foreign key
	Postings   []Posting
	BusinessID string  `gorm:"type:varchar(24);not null;unique_index"`
	Balance    float64 `gorm:"not null;default:0"`
}
