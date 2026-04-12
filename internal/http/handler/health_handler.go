package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hamsacumar/travel_backend_api/internal/http/response"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	log.Printf(fmt.Sprintf(`health check working`))
	_ = json.NewEncoder(w).Encode(response.HealthResponse{Status: "ok"})

}
