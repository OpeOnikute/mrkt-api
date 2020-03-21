package handlers

import (
	"context"
	"encoding/json"
	"io"
	"mrkt/db"
	"mrkt/models"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// CreateEntry ...
func CreateEntry(params io.Reader) (*mongo.InsertOneResult, error) {
	entry := models.GetDefaultEntry()
	err := json.NewDecoder(params).Decode(&entry)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.Collection.InsertOne(ctx, entry)
}

// GetAllEntries gets all entries
func GetAllEntries() ([]models.Entry, error) {
	var results []models.Entry
	cursor, err := db.Collection.Find(context.Background(), bson.M{})
	if err != nil {
		return results, err
	}
	err = cursor.All(context.TODO(), &results)
	return results, err
}

// GetEntryByID exposes a function to retrieve an entry by it's ID
func GetEntryByID(requestID string) (models.Entry, error) {
	id, _ := primitive.ObjectIDFromHex(requestID)
	var entry models.Entry
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := db.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&entry)
	return entry, err
}

// UpdateEntryByID exposes a function to update an entry by it's ID
func UpdateEntryByID(requestID string, entry models.Entry) (*mongo.UpdateResult, error) {

	entry.Updated = time.Now()

	update := make(map[string]interface{})
	update["$set"] = entry

	id, _ := primitive.ObjectIDFromHex(requestID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := db.Collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	return result, err
}

// DeleteEntryByID ...
func DeleteEntryByID(entry models.Entry) (*mongo.UpdateResult, error) {

	entry.Status = "deleted"
	entry.Updated = time.Now()

	update := make(map[string]interface{})
	update["$set"] = entry

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := db.Collection.UpdateOne(ctx, bson.M{"_id": entry.ID}, update)
	return result, err
}
