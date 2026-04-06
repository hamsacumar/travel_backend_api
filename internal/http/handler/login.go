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

const loginLogPrefix = `travels-api.internal.http.handler.login`

func DecodeLoginRequest(_ context.Context, r *http.Request) (request.LoginInput, error) {
	var req request.LoginInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] DecodeLoginRequest error: %v`, loginLogPrefix, err))
		return req, err
	}
	return req, nil
}

func LoginEndpoint(h *Handler) func(ctx context.Context, req request.LoginInput) (res interface{}, err error) {
	return func(ctx context.Context, req request.LoginInput) (res interface{}, err error) {
		res, err = h.AuthUsecase.Login(req.Phone)
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] LoginEndpoint error: %v`, loginLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Login successful, OTP sent for phone: %s`, loginLogPrefix, req.Phone))
		return res, nil
	}
}

func EncodeLoginResponse(_ context.Context, w http.ResponseWriter, res interface{}) (response interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] EncodeLoginResponse error: %v`, loginLogPrefix, err))
		return nil, err
	}
	return res, nil
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	req, err := DecodeLoginRequest(ctx, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response.ErrorResponse{Message: err.Error()})
		return
	}

	res, err := LoginEndpoint(h)(ctx, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(response.ErrorResponse{Message: err.Error()})
		return
	}

	_, err = EncodeLoginResponse(ctx, w, res)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] Login encode error: %v`, loginLogPrefix, err))
	}
}
