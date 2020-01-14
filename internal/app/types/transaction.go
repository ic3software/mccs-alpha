package types

import (
	"time"
)

type Transaction struct {
	ID               uint // Journal ID
	TransactionID    string
	IsInitiator      bool
	InitiatedBy      uint
	FromID           uint
	FromEmail        string
	FromBusinessName string
	ToID             uint
	ToEmail          string
	ToBusinessName   string
	Amount           float64
	Description      string
	Status           string
	CreatedAt        time.Time
}
