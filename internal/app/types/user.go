package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// LoginInfo is shared by user and admin user model.
type LoginInfo struct {
	CurrentLoginIP   string    `json:"currentLoginIP,omitempty" bson:"currentLoginIP,omitempty"`
	CurrentLoginDate time.Time `json:"currentLoginDate,omitempty" bson:"currentLoginDate,omitempty"`
	LastLoginIP      string    `json:"lastLoginIP,omitempty" bson:"lastLoginIP,omitempty"`
	LastLoginDate    time.Time `json:"lastLoginDate,omitempty" bson:"lastLoginDate,omitempty"`
}

// User is the model representation of an user in the data model.
type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DeletedAt time.Time          `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`

	FirstName string             `json:"firstName,omitempty" bson:"firstName,omitempty"`
	LastName  string             `json:"lastName,omitempty" bson:"lastName,omitempty"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty"`
	Password  string             `json:"password,omitempty" bson:"password,omitempty"`
	Telephone string             `json:"telephone,omitempty" bson:"telephone,omitempty"`
	CompanyID primitive.ObjectID `json:"companyID,omitempty" bson:"companyID,omitempty"`

	CurrentLoginIP   string    `json:"currentLoginIP,omitempty" bson:"currentLoginIP,omitempty"`
	CurrentLoginDate time.Time `json:"currentLoginDate,omitempty" bson:"currentLoginDate,omitempty"`
	LastLoginIP      string    `json:"lastLoginIP,omitempty" bson:"lastLoginIP,omitempty"`
	LastLoginDate    time.Time `json:"lastLoginDate,omitempty" bson:"lastLoginDate,omitempty"`

	LoginAttempts     int       `json:"loginAttempts,omitempty" bson:"loginAttempts,omitempty"`
	LastLoginFailDate time.Time `json:"lastLoginFailDate,omitempty" bson:"lastLoginFailDate,omitempty"`

	ShowRecentMatchedTags    bool                 `json:"showRecentMatchedTags,omitempty" bson:"showRecentMatchedTags,omitempty"`
	FavoriteBusinesses       []primitive.ObjectID `json:"favoriteBusinesses,omitempty" bson:"favoriteBusinesses,omitempty"`
	DailyNotification        bool                 `json:"dailyNotification,omitempty" bson:"dailyNotification,omitempty"`
	LastNotificationSentDate time.Time            `json:"lastNotificationSentDate,omitempty" bson:"lastNotificationSentDate,omitempty"`
}

// UserESRecord is the data that will store into the elastic search.
type UserESRecord struct {
	UserID    string `json:"userID"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Email     string `json:"email,omitempty"`
}

// Helper types

type FindUserResult struct {
	Users           []*User
	NumberOfResults int
	TotalPages      int
}
