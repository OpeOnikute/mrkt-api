package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/OpeOnikute/mrkt-api/constants"
	"github.com/OpeOnikute/mrkt-api/db"
	"github.com/OpeOnikute/mrkt-api/models"

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
		lg := new(Logger)
		lg.Error(err.Error())
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
func GetAllEntries(query bson.M) ([]models.Entry, error) {
	// init empty array so we don't send null as json response
	results := []models.Entry{}
	cursor, err := db.Collections.Entries.Find(context.Background(), query)
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

// GetLocationRanking houses the core logic to classify how safe a location is.
func GetLocationRanking(lat, long float64) (models.LocationRanking, error) {

	var result float64
	var text string
	var ranking models.LocationRanking
	var newErr constants.CustomError

	// 5km
	maxDistance := 5000

	dayAvg := 5

	// get x days ago as date. we only care about incidents reported
	// up to x days ago.
	xDaysAgo := time.Now().AddDate(0, 0, -dayAvg)

	// toDate := time.Date(2014, time.November, 5, 0, 0, 0, 0, time.UTC)

	results := []bson.M{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	geoNearStage := bson.M{
		"$geoNear": bson.M{
			"near": bson.M{
				"type":        "Point",
				"coordinates": []float64{lat, long},
			},
			"distanceField": "dist.calculated",
			"maxDistance":   maxDistance,
			"query": bson.M{
				"status":  constants.Enabled,
				"created": bson.M{"$gte": xDaysAgo},
			},
			"includeLocs": "dist.location",
			"spherical":   false,
		},
	}

	lookupStage := bson.M{
		"$lookup": bson.M{
			"from":         "alertTypes",
			"localField":   "alertType",
			"foreignField": "_id",
			"as":           "alertType",
		},
	}
	unwindStage := bson.M{
		"$unwind": bson.M{
			"path":                       "$alertType",
			"preserveNullAndEmptyArrays": false,
		},
	}
	matchStage := bson.M{
		"$match": bson.M{
			"alertType.level": bson.M{
				"$gte": 3,
			},
		},
	}
	groupStage := bson.M{
		"$group": bson.M{
			"_id": "",
			"count": bson.M{
				"$sum": 1,
			},
		},
	}

	projectStage := bson.M{
		"$project": bson.M{
			"_id": 0,
			"avg": bson.M{
				"$divide": []interface{}{"$count", dayAvg},
			},
			"count": 1,
		},
	}

	pipeline := []bson.M{geoNearStage, lookupStage, unwindStage, matchStage, groupStage, projectStage}

	cursor, _ := db.Collections.Entries.Aggregate(ctx, pipeline)
	_ = cursor.All(context.TODO(), &results)

	// If there are no results, no incident was found. Location is safe.
	if len(results) <= 0 {
		ranking.Average = 0
		ranking.Text = constants.LOCATION_SAFE
		ranking.NumIncidents = 0
		newErr.Msg = "Aggregation failed."
		return ranking, nil
	}

	result, ok := results[0]["avg"].(float64)

	if !ok {
		newErr.Msg = "Failed to convert aggregation result to float."
		return ranking, &newErr
	}

	if 0 <= result && result <= 0.4 {
		text = constants.LOCATION_SAFE
	} else if 0.5 <= result && result <= 0.9 {
		text = constants.LOCATION_WARNING
	} else if result >= 1 {
		text = constants.LOCATION_UNSAFE
	} else {
		// handle negative cases
		text = constants.LOCATION_UNKNOWN
	}

	ranking.Text = text
	ranking.Average = result
	ranking.NumIncidents = results[0]["count"].(int32)

	return ranking, nil
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
