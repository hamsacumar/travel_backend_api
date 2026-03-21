package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hamsacumar/travels_backend/adapter/http/router"
)

func main() {

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
