package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"mrkt/db"
	"mrkt/models"
	"os"
	"time"

	geo "github.com/codingsince1985/geo-golang"

	"github.com/codingsince1985/geo-golang/google"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// ClearEntries ...
func ClearEntries() (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.Collections.Entries.DeleteMany(ctx, bson.M{})
}

// CreateEntry ...
func CreateEntry(entry *models.Entry) (*mongo.InsertOneResult, error) {

	latFR := entry.Location.Coordinates[0]
	lngFR := entry.Location.Coordinates[1]

	address, err := GetAddressFromCoordinates(latFR, lngFR)

	if err == nil {
		entry.Address = address
	} else {
		ErrorLogger.Error(err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.Collections.Entries.InsertOne(ctx, entry)
}

// GetAddressFromCoordinates ...
func GetAddressFromCoordinates(lat, long float64) (*geo.Address, error) {
	// TODO: Cache results and fetch from cache
	geocoder := google.Geocoder(os.Getenv("GOOGLE_MAPS_KEY"))
	return geocoder.ReverseGeocode(lat, long)
}

// CreateMultipleEntries ...
func CreateMultipleEntries(entries []models.Entry) (*mongo.InsertManyResult, error) {
	// convert to an interface
	data := make([]interface{}, len(entries))
	def := models.GetDefaultEntry()

	for i, entry := range entries {

		latFR := entry.Location.Coordinates[0]
		lngFR := entry.Location.Coordinates[1]

		address, err := GetAddressFromCoordinates(latFR, lngFR)

		if err == nil {
			entry.Address = address
		}

		entry.Status = def.Status
		entry.Created = def.Created
		entry.Updated = def.Updated

		data[i] = entry
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.Collections.Entries.InsertMany(ctx, data)
}

// GetAllEntries gets all entries
func GetAllEntries() ([]models.Entry, error) {
	// init empty array so we don't send null as json response
	results := []models.Entry{}
	cursor, err := db.Collections.Entries.Find(context.Background(), bson.M{})
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
	err := db.Collections.Entries.FindOne(ctx, bson.M{"_id": id}).Decode(&entry)
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

	result, err := db.Collections.Entries.UpdateOne(ctx, bson.M{"_id": id}, update)
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
	result, err := db.Collections.Entries.UpdateOne(ctx, bson.M{"_id": entry.ID}, update)
	return result, err
}

// Seed ...
func Seed() (*mongo.InsertManyResult, error) {
	jsonFile, err := os.Open("./seed_entries.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var entries []models.Entry
	// json.NewDecoder(byteValue).Decode(&entries)

	json.Unmarshal(byteValue, &entries)

	// clear the entries collection
	_, err = ClearEntries()
	if err != nil {
		panic(err)
	}

	return CreateMultipleEntries(entries)
}
