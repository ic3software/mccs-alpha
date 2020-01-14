package main

import (
	"github.com/ic3network/mccs-alpha/global"

	"context"
	"log"
	"time"

	"github.com/ic3network/mccs-alpha/internal/app/repositories/mongo"
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	global.Init()
	setUpAccount()
}

// setUpAccount reads the businesses from MongoDB and build up the accounts in PostgreSQL.
func setUpAccount() {
	log.Println("start setting up accounts in PostgreSQL")
	startTime := time.Now()
	ctx := context.Background()

	filter := bson.M{
		"deletedAt": bson.M{"$exists": false},
	}
	cur, err := mongo.DB().Collection("businesses").Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	counter := 0
	for cur.Next(ctx) {
		var b types.Business
		err := cur.Decode(&b)
		if err != nil {
			log.Fatal(err)
		}
		// Create account from business.
		err = service.Account.Create(b.ID.Hex())
		if err != nil {
			log.Fatal(err)
		}
		counter++
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(ctx)

	log.Printf("count %v\n", counter)
	log.Printf("took  %v\n\n", time.Now().Sub(startTime))
}
