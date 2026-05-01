package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hamsacumar/travel_backend_api/internal/http/request"
	"github.com/hamsacumar/travel_backend_api/internal/http/response"
)

const verifyLogPrefix = `travels-api.internal.http.handler.verify`

func DecodeVerifyRequest(_ context.Context, r *http.Request) (request.VerifyInput, error) {
	var req request.VerifyInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] DecodeVerifyRequest error: %v`, verifyLogPrefix, err))
		return req, err
	}
	return req, nil
}

func VerifyEndpoint(h *Handler) func(ctx context.Context, req request.VerifyInput) (res interface{}, err error) {
	return func(ctx context.Context, req request.VerifyInput) (res interface{}, err error) {
		res, err = h.AuthUsecase.Verify(req.Phone, req.Code, req.Role)
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] VerifyEndpoint error: %v`, verifyLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Verify successful for phone: %s`, verifyLogPrefix, req.Phone))
		return res, nil
	}
}

func EncodeVerifyResponse(_ context.Context, w http.ResponseWriter, res interface{}) (response interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] EncodeVerifyResponse error: %v`, verifyLogPrefix, err))
		return nil, err
	}
	return res, nil
}

func (h *Handler) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	req, err := DecodeVerifyRequest(ctx, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response.ErrorResponse{Message: err.Error()})
		return
	}

	res, err := VerifyEndpoint(h)(ctx, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response.ErrorResponse{Message: err.Error()})
		return
	}

	_, err = EncodeVerifyResponse(ctx, w, res)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] VerifyOTP encode error: %v`, verifyLogPrefix, err))
	}
}
