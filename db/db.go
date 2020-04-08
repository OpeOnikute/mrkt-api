package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

// Database ...
var Database *mongo.Database

// DBCollection ...
type dBCollection struct {
	Entries    *mongo.Collection
	Users      *mongo.Collection
	AlertTypes *mongo.Collection
}

// Collections ...
var Collections dBCollection

// Connect ...
func Connect() {

	var mongoURL = os.Getenv("MONGO_URL")
	var mongoDB = os.Getenv("MONGO_DATABASE")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(context.TODO())
	if err != nil {
		panic(err)
	}

	Database = client.Database(mongoDB)

	Collections.Entries = Database.Collection("entries")
	Collections.AlertTypes = Database.Collection("alertTypes")
	Collections.Users = Database.Collection("users")

	// Create indexes
	mod := mongo.IndexModel{
		Keys: bson.M{
			"location": "2dsphere", // index in ascending order
		}, Options: nil,
	}
	Collections.Entries.Indexes().CreateOne(ctx, mod)

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// confirm we can connect to the db
	err = client.Ping(ctx, readpref.Primary())

	if err != nil {
		panic(err)
	}
}
