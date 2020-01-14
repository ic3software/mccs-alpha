package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Business is the model representation of a business in the data model.
type Business struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DeletedAt time.Time          `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`

	BusinessName       string      `json:"businessName,omitempty" bson:"businessName,omitempty"`
	BusinessPhone      string      `json:"businessPhone,omitempty" bson:"businessPhone,omitempty"`
	IncType            string      `json:"incType,omitempty" bson:"incType,omitempty"`
	CompanyNumber      string      `json:"companyNumber,omitempty" bson:"companyNumber,omitempty"`
	Website            string      `json:"website,omitempty" bson:"website,omitempty"`
	Turnover           int         `json:"turnover,omitempty" bson:"turnover,omitempty"`
	Offers             []*TagField `json:"offers,omitempty" bson:"offers,omitempty"`
	Wants              []*TagField `json:"wants,omitempty" bson:"wants,omitempty"`
	Description        string      `json:"description,omitempty" bson:"description,omitempty"`
	LocationAddress    string      `json:"locationAddress,omitempty" bson:"locationAddress,omitempty"`
	LocationCity       string      `json:"locationCity,omitempty" bson:"locationCity,omitempty"`
	LocationRegion     string      `json:"locationRegion,omitempty" bson:"locationRegion,omitempty"`
	LocationPostalCode string      `json:"locationPostalCode,omitempty" bson:"locationPostalCode,omitempty"`
	LocationCountry    string      `json:"locationCountry,omitempty" bson:"locationCountry,omitempty"`
	Status             string      `json:"status,omitempty" bson:"status,omitempty"`
	AdminTags          []string    `json:"adminTags,omitempty" bson:"adminTags,omitempty"`
	// Timestamp when trading status applied
	MemberStartedAt time.Time `json:"memberStartedAt,omitempty" bson:"memberStartedAt,omitempty"`
}

type TagField struct {
	Name      string    `json:"name,omitempty" bson:"name,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}

// BusinessESRecord is the data that will store into the elastic search.
type BusinessESRecord struct {
	BusinessID      string      `json:"businessID,omitempty"`
	BusinessName    string      `json:"businessName,omitempty"`
	Offers          []*TagField `json:"offers,omitempty"`
	Wants           []*TagField `json:"wants,omitempty"`
	LocationCity    string      `json:"locationCity,omitempty"`
	LocationCountry string      `json:"locationCountry,omitempty"`
	Status          string      `json:"status,omitempty"`
	AdminTags       []string    `json:"adminTags,omitempty"`
}

// Helper types

type BusinessData struct {
	ID                 primitive.ObjectID
	BusinessName       string
	IncType            string
	CompanyNumber      string
	BusinessPhone      string
	Website            string
	Turnover           int
	Offers             []*TagField
	Wants              []*TagField
	OffersAdded        []string
	OffersRemoved      []string
	WantsAdded         []string
	WantsRemoved       []string
	Description        string
	LocationAddress    string
	LocationCity       string
	LocationRegion     string
	LocationPostalCode string
	LocationCountry    string
	Status             string
	AdminTags          []string
}

type SearchCriteria struct {
	TagType          string
	Tags             []*TagField
	CreatedOnOrAfter time.Time

	Statuses              []string // accepted", "pending", rejected", "tradingPending", "tradingAccepted", "tradingRejected"
	BusinessName          string
	LocationCountry       string
	LocationCity          string
	ShowUserFavoritesOnly bool
	FavoriteBusinesses    []primitive.ObjectID
	AdminTag              string
}

type FindBusinessResult struct {
	Businesses      []*Business
	NumberOfResults int
	TotalPages      int
}
