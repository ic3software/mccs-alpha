package service

import (
	"time"

	"github.com/ic3network/mccs-alpha/internal/app/repositories/es"
	"github.com/ic3network/mccs-alpha/internal/app/repositories/mongo"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type business struct{}

var Business = &business{}

func (b *business) FindByID(id primitive.ObjectID) (*types.Business, error) {
	bs, err := mongo.Business.FindByID(id)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (b *business) Create(business *types.BusinessData) (primitive.ObjectID, error) {
	id, err := mongo.Business.Create(business)
	if err != nil {
		return primitive.ObjectID{}, e.Wrap(err, "create business failed")
	}
	err = es.Business.Create(id, business)
	if err != nil {
		return primitive.ObjectID{}, e.Wrap(err, "create business failed")
	}
	return id, nil
}

func (b *business) UpdateBusiness(
	id primitive.ObjectID,
	business *types.BusinessData,
	isAdmin bool,
) error {
	err := es.Business.UpdateBusiness(id, business)
	if err != nil {
		return e.Wrap(err, "update business failed")
	}
	err = mongo.Business.UpdateBusiness(id, business, isAdmin)
	if err != nil {
		return e.Wrap(err, "update business failed")
	}
	return nil
}

func (b *business) SetMemberStartedAt(id primitive.ObjectID) error {
	err := mongo.Business.SetMemberStartedAt(id)
	if err != nil {
		return err
	}
	return nil
}

func (b *business) UpdateAllTagsCreatedAt(id primitive.ObjectID, t time.Time) error {
	err := es.Business.UpdateAllTagsCreatedAt(id, t)
	if err != nil {
		return e.Wrap(err, "BusinessService UpdateAllTagsCreatedAt failed")
	}
	err = mongo.Business.UpdateAllTagsCreatedAt(id, t)
	if err != nil {
		return e.Wrap(err, "BusinessService UpdateAllTagsCreatedAt failed")
	}
	return nil
}

func (b *business) FindBusiness(c *types.SearchCriteria, page int64) (*types.FindBusinessResult, error) {
	ids, numberOfResults, totalPages, err := es.Business.Find(c, page)
	if err != nil {
		return nil, e.Wrap(err, "BusinessService FindBusiness failed")
	}
	businesses, err := mongo.Business.FindByIDs(ids)
	if err != nil {
		return nil, e.Wrap(err, "BusinessService FindBusiness failed")
	}
	return &types.FindBusinessResult{
		Businesses:      businesses,
		NumberOfResults: numberOfResults,
		TotalPages:      totalPages,
	}, nil
}

func (b *business) DeleteByID(id primitive.ObjectID) error {
	err := es.Business.Delete(id.Hex())
	if err != nil {
		return e.Wrap(err, "delete business by id failed")
	}
	err = mongo.Business.DeleteByID(id)
	if err != nil {
		return e.Wrap(err, "delete business by id failed")
	}
	return nil
}

func (b *business) RenameTag(old string, new string) error {
	err := es.Business.RenameTag(old, new)
	if err != nil {
		return e.Wrap(err, "BusinessMongo RenameTag failed")
	}
	err = mongo.Business.RenameTag(old, new)
	if err != nil {
		return e.Wrap(err, "BusinessMongo RenameTag failed")
	}
	return nil
}

func (b *business) RenameAdminTag(old string, new string) error {
	err := es.Business.RenameAdminTag(old, new)
	if err != nil {
		return e.Wrap(err, "BusinessMongo RenameAdminTag failed")
	}
	err = mongo.Business.RenameAdminTag(old, new)
	if err != nil {
		return e.Wrap(err, "BusinessMongo RenameAdminTag failed")
	}
	return nil
}

func (b *business) DeleteTag(name string) error {
	err := es.Business.DeleteTag(name)
	if err != nil {
		return e.Wrap(err, "BusinessMongo DeleteTag failed")
	}
	err = mongo.Business.DeleteTag(name)
	if err != nil {
		return e.Wrap(err, "BusinessMongo DeleteTag failed")
	}
	return nil
}

func (b *business) DeleteAdminTags(name string) error {
	err := es.Business.DeleteAdminTags(name)
	if err != nil {
		return e.Wrap(err, "BusinessMongo DeleteAdminTags failed")
	}
	err = mongo.Business.DeleteAdminTags(name)
	if err != nil {
		return e.Wrap(err, "BusinessMongo DeleteAdminTags failed")
	}
	return nil
}
