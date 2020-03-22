package router

import (
	"mrkt/controllers"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// GetRouter exposes the main router
func GetRouter() http.Handler {
	router := mux.NewRouter()

	entryrouter := router.PathPrefix("/entry").Subrouter()
	entryrouter.HandleFunc("", controllers.AddEntryEndpoint).Methods("POST")
	entryrouter.HandleFunc("/{id}", controllers.UpdateEntryEndpoint).Methods("PUT")
	entryrouter.HandleFunc("", controllers.GetEntriesEndpoint).Methods("GET")
	entryrouter.HandleFunc("/{id}", controllers.GetEntryEndpoint).Methods("GET")
	entryrouter.HandleFunc("/{id}", controllers.DeleteEntryEndpoint).Methods("DELETE")

	adminrouter := router.PathPrefix("/admin").Subrouter()
	adminrouter.Use(controllers.AdminAuthenticationMiddleware)

	adminrouter.HandleFunc("", controllers.CreateUserEndpoint).Methods("POST")
	adminrouter.HandleFunc("/login", controllers.AdminLoginEndpoint).Methods("POST")
	adminrouter.HandleFunc("/{id}", controllers.UpdateUserEndpoint).Methods("PUT")
	adminrouter.HandleFunc("", controllers.GetUsersEndpoint).Methods("GET")
	adminrouter.HandleFunc("/{id}", controllers.GetUserEndpoint).Methods("GET")
	adminrouter.HandleFunc("/{id}", controllers.DeleteUserEndpoint).Methods("DELETE")

	return handlers.LoggingHandler(os.Stdout, router)
}
