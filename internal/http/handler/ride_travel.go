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

const rideLogPrefix = `travels-api.internal.http.handler.ride_driver`

type TravelRideHandler struct {
	RideUsecase *usecase.RideUsecase
}

//func NewRideHandler(rideUsecase *usecase.RideUsecase) *RideHandler {
//	return &RideHandler{RideUsecase: rideUsecase}
//}

// Decoder
func DecodeTravelAddRideRequest(_ context.Context, r *http.Request) (request.TravelRideRequest, error) {
	var req request.TravelRideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf(fmt.Sprintf(`[%s] DecodeAddRideRequest Error: %v`, rideLogPrefix, err))
		return req, err
	}
	return req, nil
}

// Encoder
func EncodeTravelAddRideResponse(_ context.Context, w http.ResponseWriter, res interface{}) (interface{}, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf(fmt.Sprintf(`[%s] EncodeAddRideResponse Error: %v`, rideLogPrefix, err))
		return nil, err
	}
	return res, nil
}

// Endpoint
func AddTravelRideEndpoint(h *TravelRideHandler) func(ctx context.Context, req request.TravelRideRequest) (interface{}, error) {
	return func(ctx context.Context, req request.TravelRideRequest) (interface{}, error) {

		ride, err := h.RideUsecase.TravelAddRide(req, ctx)
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] AddRide usecase error: %v`, rideLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Ride added successfully: %s`, rideLogPrefix, ride.RideID))
		return ride, nil
	}
}

// Handler
func (h *TravelRideHandler) TravelAddRide(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Decode
	req, err := DecodeTravelAddRideRequest(ctx, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	// Endpoint
	res, err := AddTravelRideEndpoint(h)(ctx, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	// Encode
	_, err = EncodeTravelAddRideResponse(ctx, w, res)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] AddRide encode error: %v`, rideLogPrefix, err))
	}
}
