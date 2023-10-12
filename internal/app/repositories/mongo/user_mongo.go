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

type user struct {
	c *mongo.Collection
}

var User = &user{}

func (u *user) Register(db *mongo.Database) {
	u.c = db.Collection("users")
}

func (u *user) FindByID(id primitive.ObjectID) (*types.User, error) {
	ctx := context.Background()
	user := types.User{}
	filter := bson.M{
		"_id":       id,
		"deletedAt": bson.M{"$exists": false},
	}
	err := u.c.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, e.New(e.UserNotFound, "user not found")
	}
	return &user, nil
}

func (u *user) FindByEmail(email string) (*types.User, error) {
	email = strings.ToLower(email)
	if email == "" {
		return &types.User{}, e.New(e.UserNotFound, "user not found")
	}
	user := types.User{}

	filter := bson.M{
		"email":     email,
		"deletedAt": bson.M{"$exists": false},
	}
	err := u.c.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, e.New(e.UserNotFound, "user not found")
	}
	return &user, nil
}

func (u *user) FindByBusinessID(id primitive.ObjectID) (*types.User, error) {
	user := types.User{}
	filter := bson.M{
		"companyID": id,
		"deletedAt": bson.M{"$exists": false},
	}
	err := u.c.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		return nil, e.New(e.UserNotFound, "user not found")
	}
	return &user, nil
}

func (u *user) UpdateTradingInfo(
	id primitive.ObjectID,
	data *types.TradingRegisterData,
) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"firstName": data.FirstName,
		"lastName":  data.LastName,
		"telephone": data.Telephone,
		"updatedAt": time.Now(),
	}}
	_, err := u.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "UserMongo UpdateTradingInfo failed")
	}
	return nil
}

// Create creates a user record in the table
func (u *user) Create(user *types.User) error {
	user.Email = strings.ToLower(user.Email)

	doc := bson.M{
		"firstName": user.FirstName,
		"lastName":  user.LastName,
		"email":     user.Email,
		"password":  user.Password,
		"telephone": user.Telephone,
		"companyID": user.CompanyID,
		"createdAt": time.Now(),
	}
	res, err := u.c.InsertOne(context.Background(), doc)
	if err != nil {
		return err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

func (u *user) FindByDailyNotification() ([]*types.User, error) {
	filter := bson.M{
		"dailyNotification": true,
		"deletedAt":         bson.M{"$exists": false},
	}
	projection := bson.M{
		"_id":                      1,
		"email":                    1,
		"companyID":                1,
		"lastNotificationSentDate": 1,
		"dailyNotification":        1,
	}
	findOptions := options.Find()
	findOptions.SetProjection(projection)
	cur, err := u.c.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, e.Wrap(err, "UserMongo FindByDailyNotification failed")
	}

	var users []*types.User
	for cur.Next(context.TODO()) {
		var elem types.User
		err := cur.Decode(&elem)
		if err != nil {
			return nil, e.Wrap(err, "UserMongo FindByDailyNotification failed")
		}
		users = append(users, &elem)
	}
	if err := cur.Err(); err != nil {
		return nil, e.Wrap(err, "UserMongo FindByDailyNotification failed")
	}
	cur.Close(context.TODO())

	return users, nil
}

func (u *user) FindByIDs(ids []string) ([]*types.User, error) {
	var results []*types.User

	objectIDs, err := toObjectIDs(ids)
	if err != nil {
		return nil, e.Wrap(err, "find user failed")
	}

	pipeline := newFindByIDsPipeline(objectIDs)
	cur, err := u.c.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, e.Wrap(err, "find user failed")
	}

	for cur.Next(context.TODO()) {
		var elem types.User
		err := cur.Decode(&elem)
		if err != nil {
			return nil, e.Wrap(err, "find user failed")
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		return nil, e.Wrap(err, "find user failed")
	}
	cur.Close(context.TODO())

	return results, nil
}

func (u *user) UpdatePassword(user *types.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{"password": user.Password, "updatedAt": time.Now()},
	}
	_, err := u.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "UpdatePassword failed")
	}
	return nil
}

func (u *user) UpdateUserInfo(user *types.User) error {
	user.Email = strings.ToLower(user.Email)

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"email":             user.Email,
		"firstName":         user.FirstName,
		"lastName":          user.LastName,
		"telephone":         user.Telephone,
		"dailyNotification": user.DailyNotification,
		"updatedAt":         time.Now(),
	}}
	_, err := u.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "UserMongo UpdateUserInfo failed")
	}
	return nil
}

func (u *user) AdminUpdateUser(user *types.User) error {
	user.Email = strings.ToLower(user.Email)

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"email":             user.Email,
		"firstName":         user.FirstName,
		"lastName":          user.LastName,
		"telephone":         user.Telephone,
		"dailyNotification": user.DailyNotification,
		"updatedAt":         time.Now(),
	}}
	_, err := u.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "UserMongo AdminUpdateUser failed")
	}
	return nil
}

func (u *user) UpdateLoginAttempts(
	email string,
	attempts int,
	lockUser bool,
) error {
	filter := bson.M{"email": email}
	set := bson.M{
		"loginAttempts": attempts,
		"updatedAt":     time.Now(),
	}
	if lockUser {
		set["lastLoginFailDate"] = time.Now()
	}
	update := bson.M{"$set": set}
	_, err := u.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "UserMongo UpdateLoginAttempts failed")
	}
	return nil
}

func (u *user) GetLoginInfo(id primitive.ObjectID) (*types.LoginInfo, error) {
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
	err := u.c.FindOne(context.Background(), filter, findOneOptions).
		Decode(&loginInfo)
	if err != nil {
		return nil, e.Wrap(err, "UserMongo GetLoginInfo failed")
	}
	return loginInfo, nil
}

func (u *user) UpdateLoginInfo(
	id primitive.ObjectID,
	i *types.LoginInfo,
) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"currentLoginIP":           i.CurrentLoginIP,
		"currentLoginDate":         time.Now(),
		"lastNotificationSentDate": time.Now(),
		"lastLoginIP":              i.LastLoginIP,
		"lastLoginDate":            i.LastLoginDate,
		"updatedAt":                time.Now(),
	}}
	_, err := u.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "update user login info failed")
	}
	return nil
}

func (u *user) UpdateLastNotificationSentDate(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"lastNotificationSentDate": time.Now(),
		"updatedAt":                time.Now(),
	}}
	_, err := u.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "UserMongo UpdateLastNotificationSentDate failed")
	}
	return nil
}

func (u *user) DeleteByID(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"deletedAt": time.Now(),
		"updatedAt": time.Now(),
	}}
	_, err := u.c.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		return e.Wrap(err, "delete user failed")
	}
	return nil
}

// APIs

func (u *user) ToggleShowRecentMatchedTags(id primitive.ObjectID) error {
	var res struct {
		ShowRecentMatchedTags bool `bson:"showRecentMatchedTags,omitempty"`
	}

	filter := bson.M{"_id": id}
	projection := bson.M{"showRecentMatchedTags": 1}
	findOneOptions := options.FindOne()
	findOneOptions.SetProjection(projection)
	err := u.c.FindOne(context.Background(), filter, findOneOptions).
		Decode(&res)
	if err != nil {
		return e.Wrap(err, "UserMongo ToggleShowRecentMatchedTags failed")
	}

	filter = bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{"showRecentMatchedTags": !res.ShowRecentMatchedTags},
	}
	_, err = u.c.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return e.Wrap(err, "UserMongo ToggleShowRecentMatchedTags failed")
	}
	return nil
}

func (u *user) AddToFavoriteBusinesses(uID, bID primitive.ObjectID) error {
	filter := bson.M{"_id": uID}
	update := bson.M{"$addToSet": bson.M{"favoriteBusinesses": bID}}
	_, err := u.c.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return e.Wrap(err, "UserMongo AddToFavoriteBusinesses failed")
	}
	return nil
}

func (u *user) RemoveFromFavoriteBusinesses(uID, bID primitive.ObjectID) error {
	filter := bson.M{"_id": uID}
	update := bson.M{"$pull": bson.M{"favoriteBusinesses": bID}}
	_, err := u.c.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return e.Wrap(err, "UserMongo RemoveFromFavoriteBusinesses failed")
	}
	return nil
}
