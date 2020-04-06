package models

import (
	"time"

	geo "github.com/codingsince1985/geo-golang"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Entry struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title" validate:"required"`
	Description string             `json:"description" bson:"description" validate:"required"`
	UploadedBy  interface{}        `json:"uploadedBy" bson:"uploadedBy"`
	ContentURL  string             `json:"contentURL" bson:"contentURL" validate:"required"`
	ContentType string             `json:"contentType" bson:"contentType" validate:"required"`
	Location    Location           `json:"location" bson:"location" validate:"required"`
	Address     *geo.Address       `json:"address" bson:"address"`
	AlertType   primitive.ObjectID `json:"alertType" bson:"alertType"`
	Status      string             `json:"status" bson:"status"`
	Created     time.Time          `json:"created" bson:"created"`
	Updated     time.Time          `json:"updated" bson:"updated"`
}

type Location struct {
	Type        string     `json:"type" bson:"type"`
	Coordinates [2]float64 `json:"coordinates" bson:"coordinates"`
}

// LocationRanking ...
type LocationRanking struct {
	Average      float64 `json:"average" bson:"average"`
	Text         string  `json:"text" bson:"text"`
	NumIncidents int32   `json:"numIncidents" bson:"numIncidents"`
}

// GetDefaultEntry sets the defaults for entries
func GetDefaultEntry() *Entry {
	defaultLocation := Location{
		Type: "Point",
	}
	return &Entry{
		ID:          primitive.NewObjectID(),
		Location:    defaultLocation,
		ContentType: "image",
		Status:      "enabled",
		Created:     time.Now(),
		Updated:     time.Now(),
	}
}
