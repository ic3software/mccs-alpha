package service

import (
	"github.com/ic3network/mccs-alpha/internal/app/repositories/mongo"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type adminTag struct{}

var AdminTag = &adminTag{}

func (a *adminTag) Create(name string) error {
	err := mongo.AdminTag.Create(name)
	if err != nil {
		return e.Wrap(err, "create admin tag failed")
	}
	return nil
}

func (a *adminTag) FindByName(name string) (*types.AdminTag, error) {
	adminTag, err := mongo.AdminTag.FindByName(name)
	if err != nil {
		return nil, e.Wrap(err, "AdminTagService FindByName failed")
	}
	return adminTag, nil
}

func (a *adminTag) FindByID(id primitive.ObjectID) (*types.AdminTag, error) {
	adminTag, err := mongo.AdminTag.FindByID(id)
	if err != nil {
		return nil, e.Wrap(err, "AdminTagService FindByID failed")
	}
	return adminTag, nil
}

func (a *adminTag) FindTags(
	name string,
	page int64,
) (*types.FindAdminTagResult, error) {
	result, err := mongo.AdminTag.FindTags(name, page)
	if err != nil {
		return nil, e.Wrap(err, "AdminTagService FindTags failed")
	}
	return result, nil
}

func (a *adminTag) TagStartWith(prefix string) ([]string, error) {
	tags, err := mongo.AdminTag.TagStartWith(prefix)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (a *adminTag) GetAll() ([]*types.AdminTag, error) {
	adminTags, err := mongo.AdminTag.GetAll()
	if err != nil {
		return nil, e.Wrap(err, "AdminTagService GetAll failed")
	}
	return adminTags, nil
}

func (a *adminTag) Update(tag *types.AdminTag) error {
	err := mongo.AdminTag.Update(tag)
	if err != nil {
		return e.Wrap(err, "AdminTagService Update failed")
	}
	return nil
}

func (a *adminTag) DeleteByID(id primitive.ObjectID) error {
	err := mongo.AdminTag.DeleteByID(id)
	if err != nil {
		return e.Wrap(err, "AdminTagService DeleteByID failed")
	}
	return nil
}
