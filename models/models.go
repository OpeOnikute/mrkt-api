package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

type Entry struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	UploadedBy  string             `json:"uploadedBy,omitempty" bson:"uploadedBy,omitempty"`
	Content     string             `json:"content,omitempty" bson:"content,omitempty"`
	ContentType string             `json:"contentType,omitempty" bson:"contentType,omitempty"`
	Location    Location           `json:"location,omitempty" bson:"location,omitempty"`
	Status      string             `json:"status,omitempty" bson:"status,omitempty"`
	Created     time.Time          `json:"created,omitempty" bson:"created,omitempty"`
	Updated     time.Time          `json:"updated,omitempty" bson:"updated,omitempty"`
}

type Location struct {
	Type        string `json:"type,omitempty" bson:"type,omitempty"`
	Coordinates [2]int `json:"coordinates,omitempty" bson:"coordinates,omitempty"`
}

// GetDefaultEntry sets the defaults for entries
func GetDefaultEntry() *Entry {
	defaultLocation := Location{
		Type: "point",
	}
	return &Entry{
		Location:    defaultLocation,
		ContentType: "image",
		Status:      "enabled",
		Created:     time.Now(),
		Updated:     time.Now(),
	}
}
