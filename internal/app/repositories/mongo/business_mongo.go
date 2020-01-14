package mongo

import (
	"context"
	"time"

	"github.com/ic3network/mccs-alpha/global/constant"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/helper"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type business struct {
	c *mongo.Collection
}

var Business = &business{}

func (b *business) Register(db *mongo.Database) {
	b.c = db.Collection("businesses")
}

func (b *business) FindByID(id primitive.ObjectID) (*types.Business, error) {
	ctx := context.Background()
	business := types.Business{}
	filter := bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}
	err := b.c.FindOne(ctx, filter).Decode(&business)
	if err != nil {
		return nil, e.New(e.BusinessNotFound, "business not found")
	}
	return &business, nil
}

func (b *business) UpdateTradingInfo(id primitive.ObjectID, data *types.TradingRegisterData) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"businessName":       data.BusinessName,
		"incType":            data.IncType,
		"companyNumber":      data.CompanyNumber,
		"businessPhone":      data.BusinessPhone,
		"website":            data.Website,
		"turnover":           data.Turnover,
		"description":        data.Description,
		"locationAddress":    data.LocationAddress,
		"locationCity":       data.LocationCity,
		"locationRegion":     data.LocationRegion,
		"locationPostalCode": data.LocationPostalCode,
		"locationCountry":    data.LocationCountry,
		"status":             constant.Trading.Pending,
		"updatedAt":          time.Now(),
	}}
	_, err := b.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "BusinessMongo UpdateTradingInfo failed")
	}
	return nil
}

func (b *business) UpdateBusiness(
	id primitive.ObjectID,
	data *types.BusinessData,
	isAdmin bool,
) error {
	updates := []bson.M{}

	u := bson.M{
		"businessName":       data.BusinessName,
		"businessPhone":      data.BusinessPhone,
		"incType":            data.IncType,
		"companyNumber":      data.CompanyNumber,
		"website":            data.Website,
		"turnover":           data.Turnover,
		"description":        data.Description,
		"locationAddress":    data.LocationAddress,
		"locationCity":       data.LocationCity,
		"locationRegion":     data.LocationRegion,
		"locationPostalCode": data.LocationPostalCode,
		"locationCountry":    data.LocationCountry,
		"updatedAt":          time.Now(),
	}
	if data.Status != "" {
		u["status"] = data.Status
	}
	if isAdmin {
		u["adminTags"] = data.AdminTags
	}
	updates = append(updates, bson.M{"$set": u})

	push := bson.M{}
	if len(data.OffersAdded) != 0 {
		push["offers"] = bson.M{"$each": helper.ToTagFields(data.OffersAdded)}
	}
	if len(data.WantsAdded) != 0 {
		push["wants"] = bson.M{"$each": helper.ToTagFields(data.WantsAdded)}
	}
	if len(push) != 0 {
		updates = append(updates, bson.M{"$push": push})
	}

	pull := bson.M{}
	if len(data.OffersRemoved) != 0 {
		pull["offers"] = bson.M{"name": bson.M{"$in": data.OffersRemoved}}
	}
	if len(data.WantsRemoved) != 0 {
		pull["wants"] = bson.M{"name": bson.M{"$in": data.WantsRemoved}}
	}
	if len(pull) != 0 {
		updates = append(updates, bson.M{"$pull": pull})
	}

	var writes []mongo.WriteModel
	for _, upd := range updates {
		model := mongo.NewUpdateOneModel().SetFilter(bson.M{"_id": id}).SetUpdate(upd)
		writes = append(writes, model)
	}

	_, err := b.c.BulkWrite(context.Background(), writes)
	if err != nil {
		return e.Wrap(err, "businessMongo updateBusiness failed")
	}
	return nil
}

func (b *business) SetMemberStartedAt(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"memberStartedAt": time.Now(),
	}}
	_, err := b.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "BusinessMongo SetMemberStartedAt failed")
	}
	return nil
}

// Create creates a business record in the table.
func (b *business) Create(data *types.BusinessData) (primitive.ObjectID, error) {
	doc := bson.M{
		"businessName":       data.BusinessName,
		"businessPhone":      data.BusinessPhone,
		"incType":            data.IncType,
		"companyNumber":      data.CompanyNumber,
		"website":            data.Website,
		"turnover":           data.Turnover,
		"offers":             data.Offers,
		"wants":              data.Wants,
		"description":        data.Description,
		"locationAddress":    data.LocationAddress,
		"locationCity":       data.LocationCity,
		"locationRegion":     data.LocationRegion,
		"locationPostalCode": data.LocationPostalCode,
		"locationCountry":    data.LocationCountry,
		"status":             constant.Business.Pending,
		"createdAt":          time.Now(),
	}
	res, err := b.c.InsertOne(context.Background(), doc)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (b *business) UpdateAllTagsCreatedAt(id primitive.ObjectID, t time.Time) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"offers.$[].createdAt": t,
		"wants.$[].createdAt":  t,
	}}
	_, err := b.c.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return e.Wrap(err, "businessMongo UpdateAllTagsCreatedAt failed")
	}
	return nil
}

func (b *business) FindByIDs(ids []string) ([]*types.Business, error) {
	var results []*types.Business

	objectIDs, err := toObjectIDs(ids)
	if err != nil {
		return nil, e.Wrap(err, "find business failed")
	}

	pipeline := newFindByIDsPipeline(objectIDs)
	cur, err := b.c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, e.Wrap(err, "find business failed")
	}

	for cur.Next(context.TODO()) {
		var elem types.Business
		err := cur.Decode(&elem)
		if err != nil {
			return nil, e.Wrap(err, "find business failed")
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return nil, e.Wrap(err, "find business failed")
	}
	cur.Close(context.TODO())

	return results, nil
}

func (b *business) DeleteByID(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"deletedAt": time.Now()}}
	_, err := b.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "delete business failed")
	}
	return nil
}

func (b *business) RenameTag(old string, new string) error {
	err := b.updateOffers(old, new)
	if err != nil {
		return err
	}
	err = b.updateWants(old, new)
	if err != nil {
		return err
	}
	return nil
}

func (b *business) updateOffers(old string, new string) error {
	filter := bson.M{"offers.name": old}
	update := bson.M{
		"$set": bson.M{
			"offers.$.name": new,
			"updatedAt":     time.Now(),
		},
	}
	_, err := b.c.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return e.Wrap(err, "updateOffers failed")
	}
	return nil
}

func (b *business) updateWants(old string, new string) error {
	filter := bson.M{"wants.name": old}
	update := bson.M{
		"$set": bson.M{
			"wants.$.name": new,
			"updatedAt":    time.Now(),
		},
	}
	_, err := b.c.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return e.Wrap(err, "updateWants failed")
	}
	return nil
}

func (b *business) RenameAdminTag(old string, new string) error {
	// Push the new tag tag name.
	filter := bson.M{"adminTags": old}
	update := bson.M{
		"$push": bson.M{"adminTags": new},
		"$set":  bson.M{"updatedAt": time.Now()},
	}
	_, err := b.c.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return e.Wrap(err, "RenameAdminTag failed")
	}
	// Delete the old tag name.
	filter = bson.M{"adminTags": old}
	update = bson.M{
		"$pull": bson.M{"adminTags": old},
		"$set":  bson.M{"updatedAt": time.Now()},
	}
	_, err = b.c.UpdateMany(context.Background(), filter, update)
	if err != nil {
		return e.Wrap(err, "RenameAdminTag failed")
	}
	return nil
}

func (b *business) DeleteTag(name string) error {
	filter := bson.M{
		"$or": []interface{}{
			bson.M{"offers.name": name},
			bson.M{"wants.name": name},
		},
	}
	update := bson.M{
		"$pull": bson.M{
			"offers": bson.M{"name": name},
			"wants":  bson.M{"name": name},
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}
	_, err := b.c.UpdateMany(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "DeleteTag failed")
	}
	return nil
}

func (b *business) DeleteAdminTags(name string) error {
	filter := bson.M{
		"$or": []interface{}{
			bson.M{"adminTags": name},
		},
	}
	update := bson.M{
		"$pull": bson.M{
			"adminTags": name,
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}
	_, err := b.c.UpdateMany(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "DeleteAdminTags failed")
	}
	return nil
}
