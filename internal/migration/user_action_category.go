package migration

import (
	"context"
	"log"
	"time"

	"github.com/ic3network/mccs-alpha/internal/app/repositories/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// SetUserActionCategory sets all the previous existed user actions' category as "user".
// Since we added a new field to the userAction table (category),
// we need to fill the existing userAction records.
// This migration script starts running at 2019-08-20, we can delete this in the near future.
func SetUserActionCategory() {
	log.Println("start setting user action category")
	startTime := time.Now()
	ctx := context.Background()

	filter := bson.M{
		"category": bson.M{"$exists": false},
	}
	update := bson.M{
		"$set": bson.M{
			"category": "user",
		},
	}

	_, err := mongo.DB().
		Collection("userActions").
		UpdateMany(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("took  %v\n\n", time.Now().Sub(startTime))
}
