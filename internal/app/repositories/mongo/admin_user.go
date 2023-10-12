package mongo

import (
	"context"
	"strings"
	"time"

	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/e"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type adminUser struct {
	c *mongo.Collection
}

var AdminUser = &adminUser{}

func (a *adminUser) Register(db *mongo.Database) {
	a.c = db.Collection("adminUsers")
}

func (a *adminUser) FindByEmail(email string) (*types.AdminUser, error) {
	email = strings.ToLower(email)

	if email == "" {
		return &types.AdminUser{}, e.New(e.UserNotFound, "admin user not found")
	}
	user := types.AdminUser{}
	filter := bson.M{
		"email":     email,
		"deletedAt": bson.M{"$exists": false},
	}
	err := a.c.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, e.New(e.UserNotFound, "admin user not found")
	}
	return &user, nil
}

func (a *adminUser) FindByID(id primitive.ObjectID) (*types.AdminUser, error) {
	adminUser := types.AdminUser{}
	filter := bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}
	err := a.c.FindOne(context.Background(), filter).Decode(&adminUser)
	if err != nil {
		return nil, e.New(e.UserNotFound, "admin user not found")
	}
	return &adminUser, nil
}

func (a *adminUser) GetLoginInfo(
	id primitive.ObjectID,
) (*types.LoginInfo, error) {
	loginInfo := &types.LoginInfo{}
	filter := bson.M{"_id": id}
	projection := bson.M{
		"currentLoginIP":   1,
		"currentLoginDate": 1,
		"lastLoginIP":      1,
		"lastLoginDate":    1,
	}
	findOneOptions := options.FindOne()
	findOneOptions.SetProjection(projection)
	err := a.c.FindOne(context.Background(), filter, findOneOptions).
		Decode(&loginInfo)
	if err != nil {
		return nil, e.Wrap(err, "AdminUserMongo GetLoginInfo failed")
	}
	return loginInfo, nil
}

func (a *adminUser) UpdateLoginInfo(
	id primitive.ObjectID,
	i *types.LoginInfo,
) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"currentLoginIP":   i.CurrentLoginIP,
		"currentLoginDate": time.Now(),
		"lastLoginIP":      i.LastLoginIP,
		"lastLoginDate":    i.LastLoginDate,
		"updatedAt":        time.Now(),
	}}
	_, err := a.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "AdminUserMongo UpdateLoginInfo failed")
	}
	return nil
}
