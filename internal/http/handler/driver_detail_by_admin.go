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

type AdminDetailHandler struct {
	DetailUsecase *usecase.DetailUsecase
}

// Decoder
func DecodeGetDriverDetailRequestByAdmin(_ context.Context, r *http.Request) (request.DriverDetailRequest, error) {
	phone := r.URL.Query().Get("phone_number")
	busNumber := r.URL.Query().Get("bus_number")
	if phone == "" && busNumber == "" {
		return request.DriverDetailRequest{}, fmt.Errorf("provide aleast phone_number or bus_number")
	}
	return request.DriverDetailRequest{Phone: phone, BusNumber: busNumber}, nil
}

// Encoder
func EncodeGetDriverDetailResponseByAdmin(_ context.Context, w http.ResponseWriter, res interface{}, status int) (interface{}, error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] EncodeGetDriverDetailResponse Error: %v`, driverDetailLogPrefix, err))
		return nil, err
	}
	return res, nil
}

// Endpoint
func GetDriverDetailEndpointByAdmin(h *AdminDetailHandler) func(ctx context.Context, req request.DriverDetailRequest) (interface{}, error) {
	return func(ctx context.Context, req request.DriverDetailRequest) (interface{}, error) {
		resp, err := h.DetailUsecase.GetDriverDetailByAdmin(req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

// Handler
func (h *AdminDetailHandler) GetDriverDetailByAdmin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := DecodeGetDriverDetailRequestByAdmin(ctx, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	res, err := GetDriverDetailEndpointByAdmin(h)(ctx, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	_, err = EncodeGetDriverDetailResponseByAdmin(ctx, w, res, http.StatusOK)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] GetDriverDetail encode error: %v`, driverDetailLogPrefix, err))
	}
}
