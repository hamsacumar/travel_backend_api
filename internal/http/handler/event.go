package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hamsacumar/travel_backend_api/internal/infrastructure/event"
)

type EventHandler struct {
	EventUC *event.EventUsecase
}

func decodeMeta(_ context.Context, r *http.Request) (map[string]interface{}, error) {
	var m map[string]interface{}
	if r.Body == nil || r.ContentLength == 0 {
		return map[string]interface{}{}, nil
	}
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

func (h *EventHandler) TripCreated(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rideID := vars["ride_id"]
	meta, err := decodeMeta(r.Context(), r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	// Best-effort: any failure will be handled inside usecase by publishing ride_deleted
	_ = h.EventUC.PublishTripCreated(r.Context(), rideID, meta)
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "trip_created requested"})
}

func (h *EventHandler) TripStarted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rideID := vars["ride_id"]
	meta, err := decodeMeta(r.Context(), r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	// Attempt to ensure order, ignore any error (usecase handles failure via ride_deleted)
	_ = h.EventUC.PublishTripCreated(r.Context(), rideID, map[string]interface{}{"note": "ensuring order before start"})
	_ = h.EventUC.PublishTripStarted(r.Context(), rideID, meta)
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "trip_started requested"})
}

func (h *EventHandler) TripCompleted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rideID := vars["ride_id"]
	meta, err := decodeMeta(r.Context(), r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	_ = h.EventUC.PublishTripCompleted(r.Context(), rideID, meta)
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "trip_completed requested"})
}

func (h *EventHandler) RideCancelled(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rideID := vars["ride_id"]
	meta, err := decodeMeta(r.Context(), r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	_ = h.EventUC.PublishRideCancelled(r.Context(), rideID, meta)
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ride_cancelled requested"})
}

// RideCreated allows an internal/admin caller to manually publish ride_created
// Useful for recovery scenarios; will also (re)setup auto trip_created scheduling.
func (h *EventHandler) RideCreated(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rideID := vars["ride_id"]
	meta, err := decodeMeta(r.Context(), r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	_ = h.EventUC.PublishRideCreated(r.Context(), rideID, meta)
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ride_created requested"})
}

func (h *EventHandler) RideDeleted(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rideID := vars["ride_id"]
	meta, err := decodeMeta(r.Context(), r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}
	_ = h.EventUC.PublishRideDeleted(r.Context(), rideID, meta)
	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ride_deleted requested"})
}
