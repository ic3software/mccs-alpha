package service

import (
	"github.com/ic3network/mccs-alpha/internal/app/repositories/es"
	"github.com/ic3network/mccs-alpha/internal/app/repositories/mongo"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type trading struct{}

var Trading = &trading{}

func (t *trading) UpdateBusiness(id primitive.ObjectID, data *types.TradingRegisterData) error {
	err := es.Business.UpdateTradingInfo(id, data)
	if err != nil {
		return err
	}
	err = mongo.Business.UpdateTradingInfo(id, data)
	if err != nil {
		return err
	}
	return nil
}

func (t *trading) UpdateUser(id primitive.ObjectID, data *types.TradingRegisterData) error {
	err := es.User.UpdateTradingInfo(id, data)
	if err != nil {
		return err
	}
	err = mongo.User.UpdateTradingInfo(id, data)
	if err != nil {
		return err
	}
	return nil
}
