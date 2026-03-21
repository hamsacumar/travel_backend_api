package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hamsacumar/travel_backend_api/adapter/http/router"
	"github.com/hamsacumar/travel_backend_api/adapter/repository"
)

func main() {

	// Initialize global DB
	repository.Connect()

	// Setup router (no need to pass db, use database.DB anywhere)

	r := router.SetupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
