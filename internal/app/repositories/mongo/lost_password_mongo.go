package mongo

import (
	"context"
	"time"

	"github.com/ic3network/mccs-alpha/internal/pkg/e"

	"github.com/ic3network/mccs-alpha/internal/app/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type lostPassword struct {
	c *mongo.Collection
}

var LostPassword = &lostPassword{}

func (l *lostPassword) Register(db *mongo.Database) {
	l.c = db.Collection("lostPassword")
}

// Create creates a lost password record in the table
func (l *lostPassword) Create(lostPassword *types.LostPassword) error {
	filter := bson.M{"email": lostPassword.Email}
	update := bson.M{"$set": bson.M{
		"email":     lostPassword.Email,
		"token":     lostPassword.Token,
		"tokenUsed": false,
		"createdAt": time.Now(),
	}}
	_, err := l.c.UpdateOne(
		context.Background(),
		filter,
		update,
		options.Update().SetUpsert(true),
	)
	return err
}

func (l *lostPassword) FindByToken(token string) (*types.LostPassword, error) {
	if token == "" {
		return nil, e.New(e.TokenInvalid, "token not found")
	}
	lostPassword := types.LostPassword{}
	err := l.c.FindOne(context.Background(), types.LostPassword{Token: token}).Decode(&lostPassword)
	if err != nil {
		return nil, e.New(e.TokenInvalid, "token not found")
	}
	return &lostPassword, nil
}

func (l *lostPassword) FindByEmail(email string) (*types.LostPassword, error) {
	if email == "" {
		return nil, e.New(e.TokenInvalid, "token not found")
	}
	lostPassword := types.LostPassword{}
	err := l.c.FindOne(context.Background(), types.LostPassword{Email: email}).Decode(&lostPassword)
	if err != nil {
		return nil, e.New(e.TokenInvalid, "token not found")
	}
	return &lostPassword, nil
}

func (l *lostPassword) SetTokenUsed(token string) error {
	filter := bson.M{"token": token}
	update := bson.M{"$set": bson.M{"tokenUsed": true}}
	_, err := l.c.UpdateOne(context.Background(), filter, update)
	return err
}
