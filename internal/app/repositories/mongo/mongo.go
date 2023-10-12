package mongo

import (
	"context"
	"log"
	"time"

	"github.com/ic3network/mccs-alpha/global"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

func init() {
	global.Init()
	// TODO: set up test docker environment.
	if viper.GetString("env") == "test" {
		return
	}
	db = New()
	registerCollections(db)
}

func registerCollections(db *mongo.Database) {
	Business.Register(db)
	User.Register(db)
	UserAction.Register(db)
	AdminUser.Register(db)
	Tag.Register(db)
	AdminTag.Register(db)
	LostPassword.Register(db)
}

// New returns an initialized JWT instance.
func New() *mongo.Database {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.NewClient(
		options.Client().ApplyURI(viper.GetString("mongo.url")),
	)
	if err != nil {
		log.Fatal(err)
	}

	// connect to mongo
	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	// check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(viper.GetString("mongo.database"))
	return db
}

// For seed/migration/restore data
func DB() *mongo.Database {
	return db
}
