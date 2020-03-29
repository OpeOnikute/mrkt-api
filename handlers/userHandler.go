package handlers

import (
	"context"
	"log"
	"mrkt/constants"
	"mrkt/db"
	"mrkt/models"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

type JwtClaim struct {
	UserID   primitive.ObjectID `json:"userID"`
	Username string             `json:"username"`
	IsAdmin  bool               `json:"isAdmin"`
	jwt.StandardClaims
}

var jwtKey = []byte(os.Getenv("JWT_KEY"))

// CreateUser allows you create different types of users by initializing outside the function
func CreateUser(user *models.User) (*mongo.InsertOneResult, error) {
	// confirm the user doesn't already exist
	if existingUser, _ := GetUserByEmail(user.Email, user.IsAdmin); existingUser.Email == user.Email {
		var newErr constants.CustomError
		newErr.Msg = constants.UserExists
		return nil, &newErr
	}

	hash, _ := generatePasswordHash(user.Password)

	user.Password = hash

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

func generatePasswordHash(pwd string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return string(hash), err
	}

	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

// GenerateJWTToken ...
func GenerateJWTToken(user *models.User) (string, error) {

	// Declare the expiration time of the token
	expirationTime := time.Now().Add(24 * time.Hour)
	// Create the JWT claims, which includes the username and expiry time
	claims := &JwtClaim{
		UserID:   user.ID,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	return token.SignedString(jwtKey)
}

// VerifyJWTToken ...
func VerifyJWTToken(tknStr string, isAdmin bool) (bool, *JwtClaim) {

	// remove the bearer part
	tknStr = strings.Replace(tknStr, "Bearer ", "", -1)

	res := true

	claims := &JwtClaim{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		res = false
	}
	if !tkn.Valid {
		res = false
	}
	if claims.IsAdmin != isAdmin {
		res = false
	}
	return res, claims
}
