package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/OpeOnikute/mrkt-api/db"
	"github.com/OpeOnikute/mrkt-api/router"
)

func main() {
	PORT := os.Getenv("PORT")
	db.Connect()
	fmt.Printf("Application listening on port %s\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%s", PORT), router.GetRouter())
}
