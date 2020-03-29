package models

import (
	"context"
	"mrkt/constants"
	"mrkt/db"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// AlertType ...
type AlertType struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Name    string             `json:"name" bson:"name" validate:"required"`
	Level   int                `json:"level" bson:"level" validate:"required"`
	Status  string             `json:"status" bson:"status"`
	Created time.Time          `json:"created" bson:"created"`
	Updated time.Time          `json:"updated" bson:"updated"`
}

// AlertModel ...
type AlertModel struct{}

// Create ...
func (model AlertModel) Create(alertType AlertType) (*mongo.InsertOneResult, error) {

	alertType.ID = primitive.NewObjectID()
	alertType.Status = constants.Enabled
	alertType.Created = time.Now()
	alertType.Updated = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.Collections.AlertTypes.InsertOne(ctx, alertType)
}

// FindByID ...
func (model AlertModel) FindByID(id primitive.ObjectID) (AlertType, error) {
	var alertType AlertType
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := db.Collections.AlertTypes.FindOne(ctx, bson.M{"_id": id, "status": constants.Enabled}).Decode(&alertType)
	return alertType, err
}

// FindByName ...
func (model AlertModel) FindByName(name string) (AlertType, error) {
	var alertType AlertType
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := db.Collections.AlertTypes.FindOne(ctx, bson.M{"name": name, "status": constants.Enabled}).Decode(&alertType)
	return alertType, err
}

// FindMany ...
func (model AlertModel) FindMany(q bson.M) ([]AlertType, error) {
	results := []AlertType{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	q["status"] = constants.Enabled

	cursor, err := db.Collections.AlertTypes.Find(ctx, q)
	if err != nil {
		return results, err
	}

	err = cursor.All(context.TODO(), &results)

	return results, err
}

// UpdateByID exposes a function to update an entry by it's ID
func (model AlertModel) UpdateByID(id primitive.ObjectID, alertType AlertType) (*mongo.UpdateResult, error) {

	alertType.Updated = time.Now()

	update := make(map[string]interface{})
	update["$set"] = alertType

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := db.Collections.AlertTypes.UpdateOne(ctx, bson.M{"_id": id}, update)
	return result, err
}

// DeleteByID ...
func (model AlertModel) DeleteByID(id primitive.ObjectID) (*mongo.UpdateResult, error) {

	query := bson.M{"status": "deleted", "updated": time.Now()}

	update := make(map[string]interface{})
	update["$set"] = query

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := db.Collections.AlertTypes.UpdateOne(ctx, bson.M{"_id": id}, update)
	return result, err
}
