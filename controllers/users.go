package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/OpeOnikute/mrkt-api/constants"
	"github.com/OpeOnikute/mrkt-api/handlers"
	"github.com/OpeOnikute/mrkt-api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

// UsersController ...
type UsersController struct{}

// SignupEndpoint ...
func (c UsersController) SignupEndpoint(response http.ResponseWriter, request *http.Request) {

	user := models.GetDefaultUser()

	if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error(), defaultRes)
		return
	}

	user.IsAdmin = false

	if ok, errors := validateRequest(user); !ok {
		SendErrorResponse(response, http.StatusBadRequest, constants.InvalidParams, errors)
		return
	}

	if _, err := handlers.CreateUser(user); err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error(), defaultRes)
		return
	}

	// generate jwt token and send
	token, err := handlers.GenerateJWTToken(user)

	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error(), defaultRes)
		return
	}

	data := map[string]string{"token": token}

	SendSuccessResponse(response, data)
}

// LoginEndpoint ...
func (c UsersController) LoginEndpoint(response http.ResponseWriter, request *http.Request) {
	var body loginBody

	err := json.NewDecoder(request.Body).Decode(&body)

	if err != nil {
		SendErrorResponse(response, http.StatusBadRequest, err.Error(), defaultRes)
		return
	}

	user, err := handlers.GetUserByEmail(body.Email, false)
	if err != nil {
		SendQueryErrorResponse(response, err, "user")
		return
	}

	// compare passwords
	if correct := handlers.ComparePasswords(user.Password, []byte(body.Password)); correct != true {
		SendErrorResponse(response, http.StatusUnauthorized, constants.IncorrectCredentials, defaultRes)
		return
	}

	// generate jwt token and send
	token, err := handlers.GenerateJWTToken(&user)

	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error(), defaultRes)
		return
	}

	data := map[string]string{"token": token}

	// data := make(map[string]string)
	// data["token"] = token

	SendSuccessResponse(response, data)
}

// DashboardEndpoint ...
func (c UsersController) DashboardEndpoint(response http.ResponseWriter, request *http.Request) {

	userID := request.Context().Value("UserID")
	if userID == nil {
		SendErrorResponse(response, http.StatusInternalServerError, "Something went wrong whilst fetching your data. Please try again.", defaultRes)
		return
	}

	// convert interface to object id
	id, ok := userID.(primitive.ObjectID)

	if !ok {
		SendErrorResponse(response, http.StatusInternalServerError, "Something went wrong whilst fetching your data. Please try again.", defaultRes)
		return
	}

	user, err := handlers.GetUserByID(id.Hex(), false)

	if err != nil {
		SendQueryErrorResponse(response, err, "user")
		return
	}

	entries, err := handlers.GetAllEntries(bson.M{"uploadedBy": id})

	if err != nil {
		SendQueryErrorResponse(response, err, "entry")
		return
	}

	data := make(map[string]interface{})
	data["user"] = user
	data["entries"] = entries

	SendSuccessResponse(response, data)
}

// UserAuthenticationMiddleware is a Middleware function, which will be called for each request
func (c UsersController) UserAuthenticationMiddleware(next http.Handler) http.Handler {

	unauthenticated := []string{"/users/login", "/users/sign-up"}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ensure we are not validating an unauthenticated route
		url := r.URL.String()
		if yes := contains(unauthenticated, url); yes {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")

		if token == "" {
			SendErrorResponse(w, http.StatusForbidden, constants.AccessDenied, defaultRes)
			return
		}

		if valid, claim := handlers.VerifyJWTToken(token, false); valid {
			// Pass down the request to the next middleware (or final handler)
			ctx := context.WithValue(r.Context(), "UserID", claim.UserID) // nolint
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			// Write an error and stop the handler chain
			SendErrorResponse(w, http.StatusForbidden, constants.AccessDenied, defaultRes)
		}
	})
}
