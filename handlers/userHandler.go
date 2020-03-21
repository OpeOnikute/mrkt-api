package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mrkt/db"
	"mrkt/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

// ResError ...
type ResError struct {
	Msg string
}

func (e *ResError) Error() string {
	return fmt.Sprintf("%s", e.Msg)
}

// CreateUser allows you create different types of users by initializing outside the function
func CreateUser(params io.Reader, isAdmin bool) (*mongo.InsertOneResult, error) {
	user := models.GetDefaultUser()
	err := json.NewDecoder(params).Decode(&user)
	if err != nil {
		return nil, err
	}

	// confirm the user doesn't already exist
	if existingUser, _ := GetUserByEmail(user.Email, isAdmin); existingUser.Email == user.Email {
		var newErr ResError
		newErr.Msg = "This user already exists"
		return nil, &newErr
	}

	user.IsAdmin = isAdmin
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.Collections.Users.InsertOne(ctx, user)
}

// GetAllUsers gets all entries
func GetAllUsers(isAdmin bool) ([]models.User, error) {
	var results []models.User
	cursor, err := db.Collections.Users.Find(context.Background(), bson.M{"isAdmin": isAdmin})
	if err != nil {
		return results, err
	}
	err = cursor.All(context.TODO(), &results)
	return results, err
}

// GetUserByID exposes a function to retrieve an user by it's ID
func GetUserByID(requestID string, isAdmin bool) (models.User, error) {
	id, _ := primitive.ObjectIDFromHex(requestID)
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := db.Collections.Users.FindOne(ctx, bson.M{"_id": id, "isAdmin": isAdmin}).Decode(&user)
	return user, err
}

// GetUserByEmail exposes a function to retrieve an user by it's ID
func GetUserByEmail(email string, isAdmin bool) (models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := db.Collections.Users.FindOne(ctx, bson.M{"email": email, "isAdmin": isAdmin}).Decode(&user)
	return user, err
}

// UpdateUserByID exposes a function to update an user by it's ID
func UpdateUserByID(requestID string, user models.User) (*mongo.UpdateResult, error) {

	user.Updated = time.Now()

	update := make(map[string]interface{})
	update["$set"] = user

	id, _ := primitive.ObjectIDFromHex(requestID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := db.Collections.Users.UpdateOne(ctx, bson.M{"_id": id}, update)
	return result, err
}

// DeleteUserByID ...
func DeleteUserByID(user models.User) (*mongo.UpdateResult, error) {

	user.Status = "deleted"
	user.Updated = time.Now()

	update := make(map[string]interface{})
	update["$set"] = user

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := db.Collections.Users.UpdateOne(ctx, bson.M{"_id": user.ID}, update)
	return result, err
}
