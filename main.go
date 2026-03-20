package main

import (
	"log"
	"net/http"

	"github.com/hamsacumar/travels_backend/adapter/http/router"
)

func main() {

	r := router.SetupRouter()

	log.Println("Server running on port 8080")

	http.ListenAndServe(":8080", r)

}
