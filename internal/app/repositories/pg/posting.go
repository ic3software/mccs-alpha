package pg

import (
	"time"

	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
)

type posting struct{}

var Posting = &posting{}

func (t *posting) FindInRange(
	from time.Time,
	to time.Time,
) ([]*types.Posting, error) {
	var result []*types.Posting
	err := db.Raw(`
	SELECT P.amount, P.created_at
	FROM postings AS P
	WHERE P.created_at BETWEEN ? AND ?
	ORDER BY P.created_at DESC
	`, from, to).Scan(&result).Error
	if err != nil {
		return nil, e.Wrap(err, "pg.Posting.FindInRange failed")
	}
	return result, nil
}
