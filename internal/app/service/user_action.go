package service

import (
	"github.com/ic3network/mccs-alpha/internal/app/repositories/mongo"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
)

type userAction struct{}

var UserAction = &userAction{}

func (u *userAction) Log(log *types.UserAction) error {
	if log == nil {
		return nil
	}
	err := mongo.UserAction.Log(log)
	if err != nil {
		return e.Wrap(err, "UserActionService Log failed")
	}
	return nil
}

func (u *userAction) Find(c *types.UserActionSearchCriteria, page int64) ([]*types.UserAction, int, error) {
	userActions, totalPages, err := mongo.UserAction.Find(c, page)
	if err != nil {
		return nil, 0, e.Wrap(err, "UserActionService Find failed")
	}
	return userActions, totalPages, nil
}
