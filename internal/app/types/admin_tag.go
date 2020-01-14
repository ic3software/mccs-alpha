package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminTag is the model representation of an admin tag in the data model.
type AdminTag struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DeletedAt time.Time          `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`

	Name string `json:"name,omitempty" bson:"name,omitempty"`
}

// Helper types

type FindAdminTagResult struct {
	AdminTags       []*AdminTag
	NumberOfResults int
	TotalPages      int
}
