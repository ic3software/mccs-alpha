package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LostPassword is the model representation of a lost password in the data model.
type LostPassword struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Token     string             `json:"token,omitempty" bson:"token,omitempty"`
	TokenUsed bool               `json:"tokenUsed,omitempty" bson:"tokenUsed,omitempty"`
}
