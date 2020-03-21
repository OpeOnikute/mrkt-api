package main

import (
	"fmt"
	"net/http"

	"mrkt/db"
	"mrkt/router"
)

// PORT defines the port the application is running on
const PORT = 12345

func main() {
	fmt.Printf("Application listening on port %d\n", PORT)
	db.Connect()
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), router.GetRouter())
}
