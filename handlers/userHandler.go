package handlers

import (
	"context"
	"log"
	"github.com/OpeOnikute/mrkt-api/constants"
	"github.com/OpeOnikute/mrkt-api/db"
	"github.com/OpeOnikute/mrkt-api/models"
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

// GetUser exposes a function to retrieve an user using any query
func GetUser(q bson.M) (models.User, error) {
	q["status"] = "enabled"
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := db.Collections.Users.FindOne(ctx, q).Decode(&user)
	return user, err
}

// GetUserByID exposes a function to retrieve an user by it's ID
func GetUserByID(requestID string, isAdmin bool) (models.User, error) {
	id, _ := primitive.ObjectIDFromHex(requestID)
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err := db.Collections.Users.FindOne(ctx, bson.M{"_id": id, "isAdmin": isAdmin}).Decode(&user)

	ranking, err := getUserRanking(user)
	if err != nil {
		return user, err
	}

	user.Ranking.LastUpdated = ranking.LastUpdated
	user.Ranking.Rank = ranking.Rank

	// Hide password
	user.Password = ""

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

// ComputeHighestAlphaRanking calculates the top alpha using the
// number of incidents reported. This highest number of incidents
// reported is then used in the percentile distribution to get
// the ranks of other users when we compute them.
// If two users tie, the user is selected at random as we limit
// the result to one user.
func ComputeHighestAlphaRanking() {

	results := []bson.M{}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	matchStage := bson.M{"$match": bson.M{"status": "enabled"}}
	lookupStage := bson.M{
		"$lookup": bson.M{
			"from": "entries",
			"let":  bson.M{"user_id": "$_id"},
			"pipeline": []bson.M{
				bson.M{
					"$match": bson.M{
						"$expr": bson.M{
							"$eq": []string{"$uploadedBy", "$$user_id"},
						},
					},
				},
				bson.M{
					"$match": bson.M{
						"$expr": bson.M{
							"$eq": []string{"$status", "enabled"},
						},
					},
				},
			},
			"as": "entries",
		},
	}
	unwindStage := bson.M{
		"$unwind": bson.M{
			"path":                       "$entries",
			"preserveNullAndEmptyArrays": false,
		},
	}
	groupStage := bson.M{
		"$group": bson.M{
			"_id": "$email",
			"count": bson.M{
				"$sum": 1,
			},
		},
	}

	sortStage := bson.M{
		"$sort": bson.M{
			"count": -1,
		},
	}

	limitStage := bson.M{
		"$limit": 1,
	}

	pipeline := []bson.M{matchStage, lookupStage, unwindStage, groupStage, sortStage, limitStage}

	cursor, _ := db.Collections.Users.Aggregate(ctx, pipeline)
	_ = cursor.All(context.TODO(), &results)

	if len(results) > 0 {
		topAlpha := results[0]
		// make type assertion to convert interface to string and int
		email := topAlpha["_id"].(string)
		count := topAlpha["count"].(int32)
		changeTopAlpha(email, count)
	}
}

func changeTopAlpha(email string, numIncidents int32) {
	// find and update the current top alpha to false
	currentAlpha, err := GetUser(bson.M{"ranking.isTopAlpha": true})

	if (err != nil) && (err != mongo.ErrNoDocuments) {
		panic(err)
	}

	// if user the same user isn't the new alpha, commot am
	if currentAlpha.Email != "" && currentAlpha.Email != email {
		_, err := removeTopAlpha(currentAlpha)
		if err != nil {
			panic(err)
		}
		log.Printf("Removed old alpha: %s", currentAlpha.Email)
	}

	newAlpha, err := GetUser(bson.M{"email": email})

	// update this top alpha to true and add number of incidents
	_, err = addNewAlpha(newAlpha, numIncidents)
	if err != nil {
		panic(err)
	}
	log.Printf("New alpha added: %s", email)
}

func addNewAlpha(user models.User, numIncidents int32) (*mongo.UpdateResult, error) {
	user.Ranking.IsTopAlpha = true
	user.Ranking.NumIncidents = numIncidents
	user.Ranking.Rank = constants.ALPHA_RANK
	user.Ranking.LastUpdated = time.Now()
	stringID := user.ID.Hex()
	return UpdateUserByID(stringID, user)
}

func removeTopAlpha(user models.User) (*mongo.UpdateResult, error) {
	user.Ranking.IsTopAlpha = false
	stringID := user.ID.Hex()
	return UpdateUserByID(stringID, user)
}

func computeUserRanking(user models.User) (*models.Ranking, error) {

	var rank int

	currentAlpha, err := GetUser(bson.M{"ranking.isTopAlpha": true})

	if err != nil {
		return &user.Ranking, err
	}

	// calculate where the user lies in the spectrum
	entries, err := GetAllEntries(bson.M{"uploadedBy": user.ID, "status": "enabled"})

	if err != nil {
		return &user.Ranking, err
	}

	alphaEntries := currentAlpha.Ranking.NumIncidents
	userEntries := len(entries)

	percentile := (int32(userEntries) / alphaEntries) * 100

	if percentile >= 80 {
		rank = constants.ALPHA_RANK
	} else if 40 <= percentile && percentile <= 79 {
		rank = constants.BETA_RANK
	} else {
		rank = constants.PUP_RANK
	}

	// persist new rank to DB
	user.Ranking.LastUpdated = time.Now()
	user.Ranking.Rank = rank
	stringID := user.ID.Hex()

	if _, err = UpdateUserByID(stringID, user); err != nil {
		return &user.Ranking, err
	}

	return &user.Ranking, nil
}

// getUserRanking ...
func getUserRanking(user models.User) (*models.Ranking, error) {

	ranking := user.Ranking

	// if the rank has already been calculated today, return it
	lastUpdated := ranking.LastUpdated

	// get midnight timestamp
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	if lastUpdated.Before(midnight) {
		// compute user ranking based on alpha
		newRanking, err := computeUserRanking(user)
		if err != nil {
			return &ranking, err
		}
		return newRanking, err
	}

	return &ranking, nil
}

// GetAlphaValue fetches the computed value for the highest ranking alpha
// for the day.
func GetAlphaValue() {
	// todo: cache this
	// get top alpha
	// return number of incidents
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
