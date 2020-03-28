package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User ...
type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"password" bson:"password" validate:"required"`
	IsAdmin   bool               `json:"isAdmin" bson:"isAdmin"`
	AdminRole string             `json:"adminRole" bson:"adminRole"` // super, standard
	Status    string             `json:"status" bson:"status"`
	Created   time.Time          `json:"created" bson:"created"`
	Updated   time.Time          `json:"updated" bson:"updated"`
}

// GetDefaultUser sets the defaults for users
func GetDefaultUser() *User {
	return &User{
		ID:      primitive.NewObjectID(),
		IsAdmin: false,
		Status:  "enabled",
		Created: time.Now(),
		Updated: time.Now(),
	}
}
