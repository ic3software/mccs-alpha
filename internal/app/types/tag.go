package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tag is the model representation of a tag in the data model.
type Tag struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DeletedAt time.Time          `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`

	Name         string    `json:"name,omitempty" bson:"name,omitempty"`
	OfferAddedAt time.Time `json:"offerAddedAt,omitempty" bson:"offerAddedAt,omitempty"`
	WantAddedAt  time.Time `json:"wantAddedAt,omitempty" bson:"wantAddedAt,omitempty"`
}

type TagESRecord struct {
	TagID        string    `json:"tagID,omitempty"`
	Name         string    `json:"name,omitempty"`
	OfferAddedAt time.Time `json:"offerAddedAt,omitempty"`
	WantAddedAt  time.Time `json:"wantAddedAt,omitempty"`
}

// Helper types

type FindTagResult struct {
	Tags            []*Tag
	NumberOfResults int
	TotalPages      int
}

type MatchedTags struct {
	MatchedOffers map[string][]string
	MatchedWants  map[string][]string
}
