package controllers

import (
	"encoding/json"
	"mrkt/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

// CreateUserEndpoint ...
func CreateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	isAdmin := request.URL.Query().Get("isAdmin") == "true"
	result, err := handlers.CreateUser(request.Body, isAdmin)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
	}
	json.NewEncoder(response).Encode(result)

}

// UpdateUserEndpoint ...
func UpdateUserEndpoint(response http.ResponseWriter, request *http.Request) {
	// set response headers
	response.Header().Set("content-type", "application/json")
	// get ID
	params := mux.Vars(request)

	isAdmin := request.URL.Query().Get("isAdmin") == "true"
	user, err := handlers.GetUserByID(params["id"], isAdmin)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	err = json.NewDecoder(request.Body).Decode(&user)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	// update model
	result, err := handlers.UpdateUserByID(params["id"], user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	// send reponse
	json.NewEncoder(response).Encode(result)
}

// DeleteUserEndpoint ...
func DeleteUserEndpoint(response http.ResponseWriter, request *http.Request) {
	// set response headers
	response.Header().Set("content-type", "application/json")
	// get ID
	params := mux.Vars(request)

	isAdmin := request.URL.Query().Get("isAdmin") == "true"
	user, err := handlers.GetUserByID(params["id"], isAdmin)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	// update model
	result, err := handlers.DeleteUserByID(user)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	// send reponse
	json.NewEncoder(response).Encode(result)
}

// GetUsersEndpoint ...
func GetUsersEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")

	isAdmin := request.URL.Query().Get("isAdmin") == "true"
	results, err := handlers.GetAllUsers(isAdmin)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(results)
}

// GetUserEndpoint ...
func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)

	isAdmin := request.URL.Query().Get("isAdmin") == "true"
	user, err := handlers.GetUserByID(params["id"], isAdmin)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(user)
}
