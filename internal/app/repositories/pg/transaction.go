package pg

import (
	"time"

	"github.com/ic3network/mccs-alpha/global/constant"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"github.com/ic3network/mccs-alpha/internal/pkg/pagination"
	"github.com/jinzhu/gorm"
	"github.com/segmentio/ksuid"
	"github.com/spf13/viper"
)

type transaction struct{}

var Transaction = &transaction{}

// Create makes a transaction directly.
func (t *transaction) Create(
	fromID uint,
	fromEmail string,
	fromBusinessName string,

	toID uint,
	toEmail string,
	toBusinessName string,

	amount float64,
	desc string,
) error {
	tx := db.Begin()

	journalRecord := &types.Journal{
		TransactionID:    ksuid.New().String(),
		FromID:           fromID,
		FromEmail:        fromEmail,
		FromBusinessName: fromBusinessName,
		ToID:             toID,
		ToEmail:          toEmail,
		ToBusinessName:   toBusinessName,
		Amount:           amount,
		Description:      desc,
		Type:             constant.Journal.Transfer,
		Status:           constant.Transaction.Completed,
	}
	err := tx.Create(journalRecord).Error
	if err != nil {
		tx.Rollback()
		return e.Wrap(err, "pg.Transaction.Create")
	}

	journalID := journalRecord.ID

	// Create postings.
	err = tx.Create(
		&types.Posting{
			AccountID: fromID,
			JournalID: journalID,
			Amount:    -amount,
		},
	).Error
	if err != nil {
		tx.Rollback()
		return e.Wrap(err, "pg.Transaction.Create")
	}
	err = tx.Create(
		&types.Posting{AccountID: toID, JournalID: journalID, Amount: amount},
	).Error
	if err != nil {
		tx.Rollback()
		return e.Wrap(err, "pg.Transaction.Create")
	}

	// Update accounts' balance.
	err = tx.Model(&types.Account{}).
		Where("id = ?", fromID).
		Update("balance", gorm.Expr("balance - ?", amount)).
		Error
	if err != nil {
		tx.Rollback()
		return e.Wrap(err, "pg.Transaction.Create")
	}
	err = tx.Model(&types.Account{}).
		Where("id = ?", toID).
		Update("balance", gorm.Expr("balance + ?", amount)).
		Error
	if err != nil {
		tx.Rollback()
		return e.Wrap(err, "pg.Transaction.Create")
	}

	return tx.Commit().Error
}

// Propose proposes a transaction.
func (t *transaction) Propose(
	initiatedBy uint,

	fromID uint,
	fromEmail string,
	fromBusinessName string,

	toID uint,
	toEmail string,
	toBusinessName string,

	amount float64,
	desc string,
) (*types.Transaction, error) {
	journalRecord := &types.Journal{
		TransactionID:    ksuid.New().String(),
		InitiatedBy:      initiatedBy,
		FromID:           fromID,
		FromEmail:        fromEmail,
		FromBusinessName: fromBusinessName,
		ToID:             toID,
		ToEmail:          toEmail,
		ToBusinessName:   toBusinessName,
		Amount:           amount,
		Description:      desc,
		Type:             constant.Journal.Transfer,
		Status:           constant.Transaction.Initiated,
	}
	err := db.Create(journalRecord).Error
	if err != nil {
		return nil, e.Wrap(err, "pg.Transaction.Create failed")
	}
	return &types.Transaction{
		TransactionID:    journalRecord.TransactionID,
		InitiatedBy:      initiatedBy,
		FromID:           fromID,
		FromEmail:        fromEmail,
		FromBusinessName: fromBusinessName,
		ToID:             toID,
		ToEmail:          toEmail,
		ToBusinessName:   toBusinessName,
		Amount:           amount,
		Description:      desc,
		Status:           journalRecord.Status,
	}, nil
}

// Find finds a transaction.
func (t *transaction) Find(transactionID uint) (*types.Transaction, error) {
	var result types.Transaction
	err := db.Raw(`
	SELECT
		J.id, J.transaction_id, J.initiated_by, J.from_id, J.from_email, J.from_business_name,
		J.to_id, J.to_email, J.to_business_name, J.amount, J.status
	FROM journals AS J
	WHERE J.id = ?
	LIMIT 1
	`, transactionID).Scan(&result).Error

	if err != nil {
		return nil, e.Wrap(err, "pg.Transaction.Find failed")
	}
	return &result, nil
}

// Cancel cancels a transaction.
func (t *transaction) Cancel(transactionID uint, reason string) error {
	err := db.Exec(`
	UPDATE journals
	SET status=?, cancellation_reason = ?, updated_at=?
	WHERE id=?
	`, constant.Transaction.Cancelled, reason, time.Now(), transactionID).Error

	if err != nil {
		return e.Wrap(err, "pg.Transaction.Cancel failed")
	}
	return nil
}

// Accept accepts a transaction.
func (t *transaction) Accept(
	transactionID uint,
	fromID uint,
	toID uint,
	amount float64,
) error {
	tx := db.Begin()

	// Create postings.
	err := tx.Create(&types.Posting{
		AccountID: fromID,
		JournalID: transactionID,
		Amount:    -amount,
	}).Error
	if err != nil {
		tx.Rollback()
		return e.Wrap(err, "pg.Transaction.Accept")
	}
	err = tx.Create(&types.Posting{
		AccountID: toID,
		JournalID: transactionID,
		Amount:    amount,
	}).Error
	if err != nil {
		tx.Rollback()
		return e.Wrap(err, "pg.Transaction.Accept")
	}

	// Update accounts' balance.
	err = tx.Model(&types.Account{}).
		Where("id = ?", fromID).
		Update("balance", gorm.Expr("balance - ?", amount)).
		Error
	if err != nil {
		tx.Rollback()
		return e.Wrap(err, "pg.Transaction.Accept")
	}
	err = tx.Model(&types.Account{}).
		Where("id = ?", toID).
		Update("balance", gorm.Expr("balance + ?", amount)).
		Error
	if err != nil {
		tx.Rollback()
		return e.Wrap(err, "pg.Transaction.Accept")
	}

	// Update the transaction status.
	err = tx.Exec(`
	UPDATE journals
	SET status=?, updated_at=?
	WHERE id=?
	`, constant.Transaction.Completed, time.Now(), transactionID).Error
	if err != nil {
		tx.Rollback()
		return e.Wrap(err, "pg.Transaction.Accept")
	}

	return tx.Commit().Error
}

// FindPendings finds the pending transactions.
func (t *transaction) FindPendings(id uint) ([]*types.Transaction, error) {
	var result []*types.Transaction
	err := db.Raw(`
	SELECT
		J.id, J.transaction_id, CAST((CASE WHEN J.initiated_by = ? THEN 1 ELSE 0 END) AS BIT) AS "is_initiator",
		J.id, J.initiated_by, J.from_id, J.from_email, J.to_id, J.from_business_name, J.to_business_name,
		J.to_email, J.amount, J.description, J.created_at
	FROM journals AS J
	WHERE (J.from_id = ? OR J.to_id = ?) AND J.status = ?
	ORDER BY J.created_at DESC
	`, id, id, id, constant.Transaction.Initiated).Scan(&result).Error

	if err != nil {
		return nil, e.Wrap(err, "pg.Transaction.FindPendingTransactions failed")
	}
	return result, nil
}

// FindRecent finds the recent 3 completed transactions.
func (t *transaction) FindRecent(id uint) ([]*types.Transaction, error) {
	var result []*types.Transaction
	err := db.Raw(`
	SELECT J.transaction_id, J.from_email, J.to_email, J.from_business_name, J.to_business_name, J.description, P.amount, P.created_at
	FROM postings AS P
	INNER JOIN journals AS J ON J."id" = P."journal_id"
	WHERE P.account_id = ?
	ORDER BY P.created_at DESC
	LIMIT ?
	`, id, 3).Scan(&result).Error

	if err != nil {
		return nil, e.Wrap(err, "pg.Transaction.FindRecent failed")
	}
	return result, nil
}

// FindInRange finds the completed transactions in specific time range.
func (t *transaction) FindInRange(
	id uint,
	dateFrom time.Time,
	dateTo time.Time,
	page int,
) ([]*types.Transaction, int, error) {
	limit := viper.GetInt("page_size")
	offset := viper.GetInt("page_size") * (page - 1)

	if dateFrom.IsZero() {
		dateFrom = constant.Date.DefaultFrom
	}
	if dateTo.IsZero() {
		dateTo = constant.Date.DefaultTo
	}

	// Add 24 hours to include the end date.
	dateTo = dateTo.Add(24 * time.Hour)

	var result []*types.Transaction
	err := db.Raw(`
	SELECT J.transaction_id, J.from_email, J.to_email, J.from_business_name, J.to_business_name, J.description, P.amount, P.created_at
	FROM postings AS P
	INNER JOIN journals AS J ON J."id" = P."journal_id"
	WHERE P.account_id = ? AND (P.created_at BETWEEN ? AND ?)
	ORDER BY P.created_at DESC
	LIMIT ? OFFSET ?
	`, id, dateFrom, dateTo, limit, offset).Scan(&result).Error

	var numberOfResults int64
	db.Model(&types.Posting{}).
		Where("account_id = ? AND (created_at BETWEEN ? AND ?)", id, dateFrom, dateTo).
		Count(&numberOfResults)
	totalPages := pagination.Pages(numberOfResults, viper.GetInt64("page_size"))

	if err != nil {
		return nil, 0, e.Wrap(err, "pg.Transaction.Find failed")
	}
	return result, totalPages, nil
}
