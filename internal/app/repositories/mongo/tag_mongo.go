package mongo

import (
	"context"
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

type tag struct {
	c *mongo.Collection
}

var Tag = &tag{}

func (t *tag) Register(db *mongo.Database) {
	t.c = db.Collection("tags")
}

// Create creates a tag record in the table
func (t *tag) Create(name string) (primitive.ObjectID, error) {
	filter := bson.M{"name": name}
	update := bson.M{"$setOnInsert": bson.M{
		"name":      name,
		"createdAt": time.Now(),
	}}
	res, err := t.c.UpdateOne(
		context.Background(),
		filter,
		update,
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.UpsertedID.(primitive.ObjectID), nil
}

func (t *tag) UpdateOffer(name string) (primitive.ObjectID, error) {
	filter := bson.M{"name": name}
	update := bson.M{
		"$set": bson.M{
			"offerAddedAt": time.Now(),
			"updatedAt":    time.Now(),
		},
		"$setOnInsert": bson.M{
			"name":      name,
			"createdAt": time.Now(),
		},
	}
	res := t.c.FindOneAndUpdate(
		context.Background(),
		filter,
		update,
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)
	if res.Err() != nil {
		return primitive.ObjectID{}, res.Err()
	}

	tag := types.Tag{}
	err := res.Decode(&tag)
	if err != nil {
		return primitive.ObjectID{}, e.Wrap(err, "TagMongo UpdateOffer failed")
	}
	return tag.ID, nil
}

func (t *tag) UpdateWant(name string) (primitive.ObjectID, error) {
	filter := bson.M{"name": name}
	update := bson.M{
		"$set": bson.M{
			"wantAddedAt": time.Now(),
			"updatedAt":   time.Now(),
		},
		"$setOnInsert": bson.M{
			"name":      name,
			"createdAt": time.Now(),
		},
	}
	res := t.c.FindOneAndUpdate(
		context.Background(),
		filter,
		update,
		options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After),
	)
	if res.Err() != nil {
		return primitive.ObjectID{}, res.Err()
	}

	tag := types.Tag{}
	err := res.Decode(&tag)
	if err != nil {
		return primitive.ObjectID{}, e.Wrap(err, "TagMongo UpdateWant failed")
	}
	return tag.ID, nil
}

func (t *tag) FindByName(name string) (*types.Tag, error) {
	tag := types.Tag{}
	filter := bson.M{
		"name":      name,
		"deletedAt": bson.M{"$exists": false},
	}
	err := t.c.FindOne(context.Background(), filter).Decode(&tag)
	if err != nil {
		return nil, e.New(e.BusinessNotFound, "Tag not found")
	}
	return &tag, nil
}

func (t *tag) FindByID(id primitive.ObjectID) (*types.Tag, error) {
	tag := types.Tag{}
	filter := bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}
	err := t.c.FindOne(context.Background(), filter).Decode(&tag)
	if err != nil {
		return nil, e.New(e.BusinessNotFound, "Tag not found")
	}
	return &tag, nil
}

func (t *tag) FindTags(name string, page int64) (*types.FindTagResult, error) {
	if page < 0 || page == 0 {
		return nil, e.New(e.InvalidPageNumber, "TagMongo FindTags failed")
	}

	var results []*types.Tag

	findOptions := options.Find()
	findOptions.SetSkip(viper.GetInt64("page_size") * (page - 1))
	findOptions.SetLimit(viper.GetInt64("page_size"))

	filter := bson.M{
		"name":      primitive.Regex{Pattern: name, Options: "i"},
		"deletedAt": bson.M{"$exists": false},
	}

	cur, err := t.c.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, e.Wrap(err, "TagMongo FindTags failed")
	}

	for cur.Next(context.TODO()) {
		var elem types.Tag
		err := cur.Decode(&elem)
		if err != nil {
			return nil, e.Wrap(err, "TagMongo FindTags failed")
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return nil, e.Wrap(err, "TagMongo FindTags failed")
	}
	cur.Close(context.TODO())

	totalCount, err := t.c.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, e.Wrap(err, "TagMongo FindTags failed")
	}
	totalPages := pagination.Pages(totalCount, viper.GetInt64("page_size"))

	return &types.FindTagResult{
		Tags:            results,
		NumberOfResults: int(totalCount),
		TotalPages:      totalPages,
	}, nil
}

func (t *tag) Rename(tag *types.Tag) error {
	filter := bson.M{"_id": tag.ID}
	update := bson.M{"$set": bson.M{
		"name":      tag.Name,
		"updatedAt": time.Now(),
	}}
	_, err := t.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "TagMongo Update failed")
	}
	return nil
}

func (t *tag) DeleteByID(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"deletedAt": time.Now(),
		"updatedAt": time.Now(),
	}}
	_, err := t.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "TagMongo DeleteByID failed")
	}
	return nil
}
