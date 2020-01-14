package types

import (
	"github.com/jinzhu/gorm"
)

type Journal struct {
	gorm.Model
	// Journal has many postings, JournalID is the foreign key
	Postings []Posting

	TransactionID string `gorm:"type:varchar(27);not null;default:''"`

	InitiatedBy uint `gorm:"type:int;not null;default:0"`

	FromID           uint   `gorm:"type:int;not null;default:0"`
	FromEmail        string `gorm:"type:varchar(120);not null;default:''"`
	FromBusinessName string `gorm:"type:varchar(120);not null;default:''"`

	ToID           uint   `gorm:"type:int;not null;default:0"`
	ToEmail        string `gorm:"type:varchar(120);not null;default:''"`
	ToBusinessName string `gorm:"type:varchar(120);not null;default:''"`

	Amount      float64 `gorm:"not null;default:0"`
	Description string  `gorm:"type:varchar(510);not null;default:''"`
	Type        string  `gorm:"type:varchar(31);not null;default:'transfer'"`
	Status      string  `gorm:"type:varchar(31);not null;default:''"`

	CancellationReason string `gorm:"type:varchar(510);not null;default:''"`
}
