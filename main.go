package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/OpeOnikute/mrkt-api/db"
	"github.com/OpeOnikute/mrkt-api/router"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("No .env file found")
	}
	PORT := os.Getenv("PORT")
	db.Connect()
	fmt.Printf("Application listening on port %s\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%s", PORT), router.GetRouter())
}
