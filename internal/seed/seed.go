package seed

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/ic3network/mccs-alpha/internal/app/repositories/es"
	"github.com/ic3network/mccs-alpha/internal/app/repositories/mongo"
	"github.com/ic3network/mccs-alpha/internal/app/service"
	"github.com/ic3network/mccs-alpha/internal/app/types"
	"github.com/ic3network/mccs-alpha/internal/pkg/bcrypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	businessData  []types.Business
	userData      []types.User
	adminUserData []types.AdminUser
	tagData       []types.Tag
	adminTagData  []types.AdminTag
)

func LoadData() {
	// Load business data.
	data, err := ioutil.ReadFile("internal/seed/data/business.json")
	if err != nil {
		log.Fatal(err)
	}
	businesses := make([]types.Business, 0)
	json.Unmarshal(data, &businesses)
	businessData = businesses

	// Load user data.
	data, err = ioutil.ReadFile("internal/seed/data/user.json")
	if err != nil {
		log.Fatal(err)
	}
	users := make([]types.User, 0)
	json.Unmarshal(data, &users)
	userData = users

	// Load admin user data.
	data, err = ioutil.ReadFile("internal/seed/data/admin_user.json")
	if err != nil {
		log.Fatal(err)
	}
	adminUsers := make([]types.AdminUser, 0)
	json.Unmarshal(data, &adminUsers)
	adminUserData = adminUsers

	// Load user tag data.
	data, err = ioutil.ReadFile("internal/seed/data/tag.json")
	if err != nil {
		log.Fatal(err)
	}
	tags := make([]types.Tag, 0)
	json.Unmarshal(data, &tags)
	tagData = tags

	// Load admin tag data.
	data, err = ioutil.ReadFile("internal/seed/data/admin_tag.json")
	if err != nil {
		log.Fatal(err)
	}
	adminTags := make([]types.AdminTag, 0)
	json.Unmarshal(data, &adminTags)
	adminTagData = adminTags
}

func Run() {
	log.Println("start seeding")
	startTime := time.Now()

	// Generate users and businesses.
	for i, b := range businessData {
		res, err := mongo.DB().Collection("businesses").InsertOne(context.Background(), b)
		if err != nil {
			log.Fatal(err)
		}
		b.ID = res.InsertedID.(primitive.ObjectID)

		bRecord := types.BusinessESRecord{
			BusinessID:      b.ID.Hex(),
			BusinessName:    b.BusinessName,
			Offers:          b.Offers,
			Wants:           b.Wants,
			LocationCity:    b.LocationCity,
			LocationCountry: b.LocationCountry,
			Status:          b.Status,
			AdminTags:       b.AdminTags,
		}
		_, err = es.Client().Index().
			Index("businesses").
			Id(b.ID.Hex()).
			BodyJson(bRecord).
			Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		// PostgresSQL - Create account from business.
		err = service.Account.Create(b.ID.Hex())
		if err != nil {
			log.Fatal(err)
		}

		u := userData[i]
		u.CompanyID = b.ID
		hashedPassword, _ := bcrypt.Hash(u.Password)
		u.Password = hashedPassword
		res, err = mongo.DB().Collection("users").InsertOne(context.Background(), u)
		if err != nil {
			log.Fatal(err)
		}
		u.ID = res.InsertedID.(primitive.ObjectID)

		{
			userID := u.ID.Hex()
			uRecord := types.UserESRecord{
				UserID:    userID,
				FirstName: u.FirstName,
				LastName:  u.LastName,
				Email:     u.Email,
			}
			_, err = es.Client().Index().
				Index("users").
				Id(userID).
				BodyJson(uRecord).
				Do(context.Background())
		}

		if err != nil {
			log.Fatal(err)
		}
	}

	// Generate admin users.
	for _, u := range adminUserData {
		hashedPassword, _ := bcrypt.Hash(u.Password)
		u.Password = hashedPassword
		_, err := mongo.DB().Collection("adminUsers").InsertOne(context.Background(), u)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Generate user tags.
	for _, t := range tagData {
		res, err := mongo.DB().Collection("tags").InsertOne(context.Background(), t)
		if err != nil {
			log.Fatal(err)
		}
		t.ID = res.InsertedID.(primitive.ObjectID)

		// ElasticSearch
		{
			tagID := t.ID.Hex()
			tagRecord := types.TagESRecord{
				TagID:        tagID,
				Name:         t.Name,
				OfferAddedAt: t.OfferAddedAt,
				WantAddedAt:  t.WantAddedAt,
			}
			_, err = es.Client().Index().
				Index("tags").
				Id(tagID).
				BodyJson(tagRecord).
				Do(context.Background())
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Generate admin tags.
	for _, a := range adminTagData {
		_, err := mongo.DB().Collection("adminTags").InsertOne(context.Background(), a)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("took  %v\n", time.Now().Sub(startTime))
}
