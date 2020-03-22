package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User ...
type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname" bson:"firstname" validate:"required"`
	Lastname  string             `json:"lastname" bson:"lastname" validate:"required"`
	Email     string             `json:"email" bson:"email" validate:"required,email"`
	Password  string             `json:"password" bson:"password" validate:"required"`
	IsAdmin   bool               `json:"isAdmin" bson:"isAdmin" validate:"required"`
	AdminRole string             `json:"adminRole" bson:"adminRole" validate:"required"` // super, standard
	Status    string             `json:"status" bson:"status"`
	Created   time.Time          `json:"created" bson:"created"`
	Updated   time.Time          `json:"updated" bson:"updated"`
}

// GetDefaultUser sets the defaults for users
func GetDefaultUser() *User {
	return &User{
		Status:  "enabled",
		Created: time.Now(),
		Updated: time.Now(),
	}
}
