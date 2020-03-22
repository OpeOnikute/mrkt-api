package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"mrkt/db"
	"mrkt/handlers"
	"mrkt/router"
)

// PORT defines the port the application is running on
const PORT = 12345

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
	db.Connect()
	handlers.InitLogger()
	defer handlers.CloseLoggers()
	fmt.Printf("Application listening on port %d\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), router.GetRouter())
}
