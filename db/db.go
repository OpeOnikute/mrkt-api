package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database ...
var Database *mongo.Database

// DBCollection ...
type dBCollection struct {
	Entries *mongo.Collection
	Users   *mongo.Collection
}

// Collections ...
var Collections dBCollection

// Connect ...
func Connect() {

	var mongoURL = os.Getenv("MONGO_URL")
	var mongoDB = os.Getenv("MONGO_DATABASE")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	Database = client.Database(mongoDB)

	Collections.Entries = Database.Collection("entries")
	Collections.Users = Database.Collection("users")
}
