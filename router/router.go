package router

import (
	"net/http"
	"os"

	"github.com/OpeOnikute/mrkt-api/constants"
	"github.com/OpeOnikute/mrkt-api/controllers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var adminController controllers.AdminController
var entriesController controllers.EntriesController
var userController controllers.UsersController

// GetRouter exposes the main router
func GetRouter() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		data := make(map[string]interface{})
		data["message"] = constants.WELCOME_MESSAGE
		controllers.SendSuccessResponse(response, data)
	}).Methods("GET")

	entryrouter := router.PathPrefix("/entry").Subrouter()
	entryrouter.HandleFunc("", entriesController.AddEntryEndpoint).Methods("POST")
	entryrouter.HandleFunc("", entriesController.GetEntriesEndpoint).Methods("GET")
	entryrouter.HandleFunc("/{id}", entriesController.GetEntryEndpoint).Methods("GET")

	locationrouter := router.PathPrefix("/location").Subrouter()
	locationrouter.HandleFunc("/safety", entriesController.GetLocationRanking).Methods("GET")

	userrouter := router.PathPrefix("/users").Subrouter()
	userrouter.Use(userController.UserAuthenticationMiddleware)

	userrouter.HandleFunc("/sign-up", userController.SignupEndpoint).Methods("POST")
	userrouter.HandleFunc("/login", userController.LoginEndpoint).Methods("POST")
	userrouter.HandleFunc("/dashboard", userController.DashboardEndpoint).Methods("GET")
	userrouter.HandleFunc("/entry", entriesController.AddEntryEndpoint).Methods("POST")
	userrouter.HandleFunc("/entry/{id}", entriesController.UpdateEntryEndpoint).Methods("PUT")
	userrouter.HandleFunc("/entry", entriesController.GetEntriesEndpoint).Methods("GET")
	userrouter.HandleFunc("/entry/{id}", entriesController.GetEntryEndpoint).Methods("GET")
	userrouter.HandleFunc("/entry/{id}", entriesController.DeleteEntryEndpoint).Methods("DELETE")

	adminrouter := router.PathPrefix("/admin").Subrouter()
	adminrouter.Use(adminController.AdminAuthenticationMiddleware)

	adminrouter.HandleFunc("", adminController.CreateUserEndpoint).Methods("POST")
	adminrouter.HandleFunc("/login", adminController.AdminLoginEndpoint).Methods("POST")
	adminrouter.HandleFunc("/users/{id}", adminController.UpdateUserEndpoint).Methods("PUT")
	adminrouter.HandleFunc("/users", adminController.GetUsersEndpoint).Methods("GET")
	adminrouter.HandleFunc("/users/{id}", adminController.GetUserEndpoint).Methods("GET")
	adminrouter.HandleFunc("/users/{id}", adminController.DeleteUserEndpoint).Methods("DELETE")
	adminrouter.HandleFunc("/alert-type", adminController.CreateAlertTypeEndpoint).Methods("POST")
	adminrouter.HandleFunc("/alert-type/{id}", adminController.UpdateAlertTypeEndpoint).Methods("PUT")
	adminrouter.HandleFunc("/alert-type", adminController.GetAlertTypesEndpoint).Methods("GET")
	adminrouter.HandleFunc("/alert-type/{id}", adminController.GetAlertTypeEndpoint).Methods("GET")
	adminrouter.HandleFunc("/alert-type/{id}", adminController.DeleteAlertTypeEndpoint).Methods("DELETE")

	return handlers.LoggingHandler(os.Stdout, router)
}
