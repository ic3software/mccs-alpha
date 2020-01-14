package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AdminUser is the model representation of an admin user in the data model.
type AdminUser struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DeletedAt time.Time          `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`

	Email    string   `json:"email,omitempty" bson:"email,omitempty"`
	Name     string   `json:"name,omitempty" bson:"name,omitempty"`
	Password string   `json:"password,omitempty" bson:"password,omitempty"`
	Roles    []string `json:"roles,omitempty" bson:"roles,omitempty"`

	CurrentLoginIP   string    `json:"currentLoginIP,omitempty" bson:"currentLoginIP,omitempty"`
	CurrentLoginDate time.Time `json:"currentLoginDate,omitempty" bson:"currentLoginDate,omitempty"`
	LastLoginIP      string    `json:"lastLoginIP,omitempty" bson:"lastLoginIP,omitempty"`
	LastLoginDate    time.Time `json:"lastLoginDate,omitempty" bson:"lastLoginDate,omitempty"`
}
