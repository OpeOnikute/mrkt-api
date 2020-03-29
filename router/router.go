package router

import (
	"mrkt/controllers"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var adminController controllers.AdminController
var entriesController controllers.EntriesController
var userController controllers.UsersController

// GetRouter exposes the main router
func GetRouter() http.Handler {
	router := mux.NewRouter()

	entryrouter := router.PathPrefix("/entry").Subrouter()
	entryrouter.HandleFunc("", entriesController.AddEntryEndpoint).Methods("POST")
	entryrouter.HandleFunc("/{id}", entriesController.UpdateEntryEndpoint).Methods("PUT")
	entryrouter.HandleFunc("", entriesController.GetEntriesEndpoint).Methods("GET")
	entryrouter.HandleFunc("/{id}", entriesController.GetEntryEndpoint).Methods("GET")
	entryrouter.HandleFunc("/{id}", entriesController.DeleteEntryEndpoint).Methods("DELETE")

	userrouter := router.PathPrefix("/users").Subrouter()
	userrouter.Use(userController.UserAuthenticationMiddleware)

	userrouter.HandleFunc("/sign-up", userController.SignupEndpoint).Methods("POST")
	userrouter.HandleFunc("/login", userController.LoginEndpoint).Methods("POST")
	userrouter.HandleFunc("/entry", entriesController.AddEntryEndpoint).Methods("POST")

	adminrouter := router.PathPrefix("/admin").Subrouter()
	adminrouter.Use(adminController.AdminAuthenticationMiddleware)

	adminrouter.HandleFunc("", adminController.CreateUserEndpoint).Methods("POST")
	adminrouter.HandleFunc("/login", adminController.AdminLoginEndpoint).Methods("POST")
	adminrouter.HandleFunc("/{id}", adminController.UpdateUserEndpoint).Methods("PUT")
	adminrouter.HandleFunc("", adminController.GetUsersEndpoint).Methods("GET")
	adminrouter.HandleFunc("/{id}", adminController.GetUserEndpoint).Methods("GET")
	adminrouter.HandleFunc("/{id}", adminController.DeleteUserEndpoint).Methods("DELETE")

	return handlers.LoggingHandler(os.Stdout, router)
}
