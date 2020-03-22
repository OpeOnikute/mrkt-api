package controllers

import (
	"encoding/json"
	"mrkt/constants"
	"mrkt/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

type loginBody struct {
	Email    string
	Password string
}

type loginResponse struct {
	Message string      `json:"message" bson:"message"`
	Data    interface{} `json:"data" bson:"data"`
}

// CreateUserEndpoint ...
func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	isAdmin := request.URL.Query().Get("isAdmin") == "true"
	result, err := handlers.CreateUser(request.Body, isAdmin)

	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccessResponse(response, result)
}

// AdminLoginEndpoint ...
func AdminLoginEndpoint(response http.ResponseWriter, request *http.Request) {
	var body loginBody

	err := json.NewDecoder(request.Body).Decode(&body)

	if err != nil {
		SendErrorResponse(response, http.StatusBadRequest, err.Error())
		return
	}

	user, err := handlers.GetUserByEmail(body.Email, true)

	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	// compare passwords
	if correct := handlers.ComparePasswords(user.Password, []byte(body.Password)); correct != true {
		SendErrorResponse(response, http.StatusUnauthorized, constants.IncorrectCredentials)
		return
	}

	// generate jwt token and send
	token, err := handlers.GenerateJWTToken(user)

	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	data := make(map[string]string)
	data["token"] = token

	res := &loginResponse{
		Message: "Login successful.",
		Data:    data,
	}

	SendSuccessResponse(response, res)
}

// UpdateUserEndpoint ...
func UpdateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	// get ID
	params := mux.Vars(request)

	isAdmin := request.URL.Query().Get("isAdmin") == "true"
	user, err := handlers.GetUserByID(params["id"], isAdmin)

	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.NewDecoder(request.Body).Decode(&user)

	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	// update model
	result, err := handlers.UpdateUserByID(params["id"], user)
	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	SendSuccessResponse(response, result)
}

// DeleteUserEndpoint ...
func DeleteUserEndpoint(response http.ResponseWriter, request *http.Request) {
	// get ID
	params := mux.Vars(request)

	isAdmin := request.URL.Query().Get("isAdmin") == "true"
	user, err := handlers.GetUserByID(params["id"], isAdmin)

	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	// update model
	result, err := handlers.DeleteUserByID(user)
	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	// send reponse
	SendSuccessResponse(response, result)
}

// GetUsersEndpoint ...
func GetUsersEndpoint(response http.ResponseWriter, request *http.Request) {
	isAdmin := request.URL.Query().Get("isAdmin") == "true"
	results, err := handlers.GetAllUsers(isAdmin)

	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}
	SendSuccessResponse(response, results)
}

// GetUserEndpoint ...
func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	isAdmin := request.URL.Query().Get("isAdmin") == "true"
	user, err := handlers.GetUserByID(params["id"], isAdmin)

	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error())
		return
	}
	SendSuccessResponse(response, user)
}

// AdminAuthenticationMiddleware is a Middleware function, which will be called for each request
func AdminAuthenticationMiddleware(next http.Handler) http.Handler {

	unauthenticated := []string{"/admin/login"}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// ensure we are not validating an unauthenticated route
		url := r.URL.String()
		if yes := contains(unauthenticated, url); yes {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")

		if token == "" {
			SendErrorResponse(w, http.StatusForbidden, constants.AccessDenied)
			return
		}

		if valid := handlers.VerifyJWTToken(token); valid {
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Write an error and stop the handler chain
			SendErrorResponse(w, http.StatusForbidden, constants.AccessDenied)
		}
	})
}

// SendErrorResponse ...
func SendErrorResponse(r http.ResponseWriter, status int, message string) {
	r.Header().Set("content-type", "application/json")
	r.WriteHeader(http.StatusInternalServerError)
	r.Write([]byte(`{ "message": "` + message + `" }`))
}

// SendSuccessResponse ...
func SendSuccessResponse(r http.ResponseWriter, result interface{}) {
	r.Header().Set("content-type", "application/json")
	json.NewEncoder(r).Encode(result)
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
