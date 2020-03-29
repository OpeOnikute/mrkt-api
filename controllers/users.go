package controllers

import (
	"context"
	"encoding/json"
	"mrkt/constants"
	"mrkt/handlers"
	"mrkt/models"
	"net/http"
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
		SendErrorResponse(response, http.StatusInternalServerError, err.Error(), defaultRes)
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
