package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hamsacumar/travel_backend_api/internal/http/request"
	"github.com/hamsacumar/travel_backend_api/internal/usecase"
)

const ridedriverLogPrefix = `travels-api.internal.http.handler.ride_driver`

type RideHandler struct {
	RideUsecase *usecase.RideUsecase
}

//func NewRideHandler(rideUsecase *usecase.RideUsecase) *RideHandler {
//	return &RideHandler{RideUsecase: rideUsecase}
//}

// Decoder
func DecodeAddRideRequest(_ context.Context, r *http.Request) (request.AddRideRequest, error) {
	var req request.AddRideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf(fmt.Sprintf(`[%s] DecodeAddRideRequest Error: %v`, ridedriverLogPrefix, err))
		return req, err
	}
	return req, nil
}

// Encoder
func EncodeAddRideResponse(_ context.Context, w http.ResponseWriter, res interface{}) (interface{}, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf(fmt.Sprintf(`[%s] EncodeAddRideResponse Error: %v`, ridedriverLogPrefix, err))
		return nil, err
	}
	return res, nil
}

// Endpoint
func AddRideEndpoint(h *RideHandler) func(ctx context.Context, req request.AddRideRequest) (interface{}, error) {
	return func(ctx context.Context, req request.AddRideRequest) (interface{}, error) {

		ride, err := h.RideUsecase.AddRide(req, ctx)
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] AddRide usecase error: %v`, ridedriverLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Ride added successfully: %s`, ridedriverLogPrefix, ride.RideID))
		return ride, nil
	}
}

// Handler
func (h *RideHandler) AddRide(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode
	req, err := DecodeAddRideRequest(ctx, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	// Endpoint
	res, err := AddRideEndpoint(h)(ctx, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	// Encode
	_, err = EncodeAddRideResponse(ctx, w, res)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] AddRide encode error: %v`, rideLogPrefix, err))
	}
}
