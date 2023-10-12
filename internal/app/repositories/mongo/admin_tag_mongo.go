package mongo

import (
	"context"
	"strings"
	"time"

	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"github.com/ic3network/mccs-alpha/internal/pkg/pagination"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type adminTag struct {
	c *mongo.Collection
}

var AdminTag = &adminTag{}

func (a *adminTag) Register(db *mongo.Database) {
	a.c = db.Collection("adminTags")
}

func (a *adminTag) Create(name string) error {
	if name == "" || len(strings.TrimSpace(name)) == 0 {
		return nil
	}

	filter := bson.M{"name": name}
	update := bson.M{
		"$setOnInsert": bson.M{"name": name, "createdAt": time.Now()},
	}
	_, err := a.c.UpdateOne(
		context.Background(),
		filter,
		update,
		options.Update().SetUpsert(true),
	)
	return err
}

func (a *adminTag) FindByName(name string) (*types.AdminTag, error) {
	adminTag := types.AdminTag{}
	filter := bson.M{
		"name":      name,
		"deletedAt": bson.M{"$exists": false},
	}
	err := a.c.FindOne(context.Background(), filter).Decode(&adminTag)
	if err != nil {
		return nil, e.New(e.BusinessNotFound, "Admin tag not found")
	}
	return &adminTag, nil
}

func (a *adminTag) FindByID(id primitive.ObjectID) (*types.AdminTag, error) {
	adminTag := types.AdminTag{}
	filter := bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}
	err := a.c.FindOne(context.Background(), filter).Decode(&adminTag)
	if err != nil {
		return nil, e.New(e.BusinessNotFound, "Admin tag not found")
	}
	return &adminTag, nil
}

func (a *adminTag) FindTags(
	name string,
	page int64,
) (*types.FindAdminTagResult, error) {
	if page < 0 || page == 0 {
		return nil, e.New(e.InvalidPageNumber, "AdminTagMongo FindTags failed")
	}

	var results []*types.AdminTag

	findOptions := options.Find()
	findOptions.SetSkip(viper.GetInt64("page_size") * (page - 1))
	findOptions.SetLimit(viper.GetInt64("page_size"))

	filter := bson.M{
		"name":      primitive.Regex{Pattern: name, Options: "i"},
		"deletedAt": bson.M{"$exists": false},
	}

	cur, err := a.c.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, e.Wrap(err, "AdminTagMongo FindTags failed")
	}

	for cur.Next(context.TODO()) {
		var elem types.AdminTag
		err := cur.Decode(&elem)
		if err != nil {
			return nil, e.Wrap(err, "AdminTagMongo FindTags failed")
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return nil, e.Wrap(err, "AdminTagMongo FindTags failed")
	}
	cur.Close(context.TODO())

	// Calculate the total page.
	totalCount, err := a.c.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, e.Wrap(err, "AdminTagMongo FindTags failed")
	}
	totalPages := pagination.Pages(totalCount, viper.GetInt64("page_size"))

	return &types.FindAdminTagResult{
		AdminTags:       results,
		NumberOfResults: int(totalCount),
		TotalPages:      totalPages,
	}, nil
}

func (a *adminTag) TagStartWith(prefix string) ([]string, error) {
	var results []string

	filter := bson.M{
		"name":      primitive.Regex{Pattern: "^" + prefix, Options: "i"},
		"deletedAt": bson.M{"$exists": false},
	}

	cur, err := a.c.Find(context.TODO(), filter)
	if err != nil {
		return nil, e.Wrap(err, "mongo.AdminTag.FindTagStartWith failed")
	}

	for cur.Next(context.TODO()) {
		var elem types.AdminTag
		err := cur.Decode(&elem)
		if err != nil {
			return nil, e.Wrap(err, "mongo.AdminTag.FindTagStartWith failed")
		}
		results = append(results, elem.Name)
	}
	if err := cur.Err(); err != nil {
		return nil, e.Wrap(err, "mongo.AdminTag.FindTagStartWith failed")
	}
	cur.Close(context.TODO())

	return results, nil
}

func (a *adminTag) GetAll() ([]*types.AdminTag, error) {
	var results []*types.AdminTag

	filter := bson.M{
		"deletedAt": bson.M{"$exists": false},
	}

	cur, err := a.c.Find(context.TODO(), filter)
	if err != nil {
		return nil, e.Wrap(err, "AdminTagMongo GetAll failed")
	}

	for cur.Next(context.TODO()) {
		var elem types.AdminTag
		err := cur.Decode(&elem)
		if err != nil {
			return nil, e.Wrap(err, "AdminTagMongo GetAll failed")
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return nil, e.Wrap(err, "AdminTagMongo GetAll failed")
	}
	cur.Close(context.TODO())

	return results, nil
}

func (a *adminTag) Update(t *types.AdminTag) error {
	filter := bson.M{"_id": t.ID}
	update := bson.M{"$set": bson.M{
		"name":      t.Name,
		"updatedAt": time.Now(),
	}}
	_, err := a.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "AdminTagMongo Update failed")
	}
	return nil
}

func (a *adminTag) DeleteByID(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"deletedAt": time.Now(),
		"updatedAt": time.Now(),
	}}
	_, err := a.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "AdminTagMongo DeleteByID failed")
	}
	return nil
}
