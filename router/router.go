package router

import (
	"mrkt/controllers"

	"github.com/gorilla/mux"
)

// GetRouter exposes the main router
func GetRouter() *mux.Router {
	router := mux.NewRouter()

	entryrouter := router.PathPrefix("/entry").Subrouter()
	entryrouter.HandleFunc("", controllers.AddEntryEndpoint).Methods("POST")
	entryrouter.HandleFunc("/{id}", controllers.UpdateEntryEndpoint).Methods("PUT")
	entryrouter.HandleFunc("", controllers.GetEntriesEndpoint).Methods("GET")
	entryrouter.HandleFunc("/{id}", controllers.GetEntryEndpoint).Methods("GET")
	entryrouter.HandleFunc("/{id}", controllers.DeleteEntryEndpoint).Methods("DELETE")

	return router
}
