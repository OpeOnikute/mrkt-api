package controllers

import (
	"encoding/json"
	"fmt"
	"mrkt/constants"
	"mrkt/handlers"
	"mrkt/models"
	"net/http"

	"github.com/gorilla/mux"
)

// AddEntryEndpoint ...
func AddEntryEndpoint(response http.ResponseWriter, request *http.Request) {

	entry := models.GetDefaultEntry()

	err := json.NewDecoder(request.Body).Decode(&entry)
	if err != nil {
		SendErrorResponse(response, http.StatusBadRequest, err.Error(), defaultRes)
		return
	}

	if ok, errors := validateRequest(entry); !ok {
		SendErrorResponse(response, http.StatusBadRequest, constants.InvalidParams, errors)
		return
	}

	fmt.Println(entry)

	result, err := handlers.CreateEntry(entry)
	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error(), defaultRes)
	}
	SendSuccessResponse(response, result)
}

// UpdateEntryEndpoint ...
func UpdateEntryEndpoint(response http.ResponseWriter, request *http.Request) {
	// set response headers
	response.Header().Set("content-type", "application/json")
	// get ID
	params := mux.Vars(request)
	// get entry
	entry, err := handlers.GetEntryByID(params["id"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	err = json.NewDecoder(request.Body).Decode(&entry)

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	// update model
	result, err := handlers.UpdateEntryByID(params["id"], entry)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	// send reponse
	json.NewEncoder(response).Encode(result)
}

// DeleteEntryEndpoint ...
func DeleteEntryEndpoint(response http.ResponseWriter, request *http.Request) {
	// set response headers
	response.Header().Set("content-type", "application/json")
	// get ID
	params := mux.Vars(request)
	// get entry
	entry, err := handlers.GetEntryByID(params["id"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	// update model
	result, err := handlers.DeleteEntryByID(entry)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	// send reponse
	json.NewEncoder(response).Encode(result)
}

// GetEntriesEndpoint ...
func GetEntriesEndpoint(response http.ResponseWriter, request *http.Request) {
	results, err := handlers.GetAllEntries()
	if err != nil {
		SendErrorResponse(response, http.StatusInternalServerError, err.Error(), defaultRes)
		return
	}
	SendSuccessResponse(response, results)
}

// GetEntryEndpoint ...
func GetEntryEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	entry, err := handlers.GetEntryByID(params["id"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}
	json.NewEncoder(response).Encode(entry)
}
