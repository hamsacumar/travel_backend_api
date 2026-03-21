package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hamsacumar/travel_backend_api/adapter/http/handler"
)

func SetupRouter() *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

	return r
}
