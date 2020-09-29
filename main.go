package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/OpeOnikute/mrkt-api/db"
	"github.com/OpeOnikute/mrkt-api/router"
	"github.com/gorilla/handlers"
)

func main() {
	PORT := os.Getenv("PORT")
	db.Connect()
	fmt.Printf("Application listening on port %s\n", PORT)

	// handle CORS requests
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	r := router.GetRouter()

	http.ListenAndServe(fmt.Sprintf(":%s", PORT), handlers.CORS(originsOk, headersOk, methodsOk)(r))
}
