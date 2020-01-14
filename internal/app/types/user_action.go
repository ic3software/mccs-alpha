package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserAction is the model representation of an user action in the data model.
type UserAction struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DeletedAt time.Time          `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`

	UserID        primitive.ObjectID `json:"userID,omitempty" bson:"userID,omitempty"`
	Email         string             `json:"email,omitempty" bson:"email,omitempty"`
	Action        string             `json:"action,omitempty" bson:"action,omitempty"`
	ActionDetails string             `json:"actionDetails,omitempty" bson:"actionDetails,omitempty"`
	Category      string             `json:"category,omitempty" bson:"category,omitempty"`
}

type UserActionSearchCriteria struct {
	Email    string
	Category string
	DateFrom time.Time
	DateTo   time.Time
}
