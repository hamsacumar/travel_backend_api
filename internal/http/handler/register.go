package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hamsacumar/travel_backend_api/internal/domain/repository"
	"github.com/hamsacumar/travel_backend_api/internal/http/request"
	"github.com/hamsacumar/travel_backend_api/internal/http/response"
	"github.com/hamsacumar/travel_backend_api/internal/usecase"
)

const (
	registerLogPrefix = `travels-api.internal.http.handler.register`
)

type Handler struct {
	AuthUsecase *usecase.AuthUsecase
	TokenRepo   repository.TokenRepository
}

func DecodeSignUpRequest(_ context.Context, r *http.Request) (request.SignUpInput, error) {
	var req request.SignUpInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] DecodeSignUpRequest Error: %v`, registerLogPrefix, err))
		return req, err
	}
	return req, nil
}

func EncodeRegisterResponse(_ context.Context, w http.ResponseWriter, res interface{}) (response interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] EncodeRegisterResponse Error: %v`, registerLogPrefix, err))
		return nil, err
	}
	return res, nil
}

func RegisterEndpoint(h *Handler) func(ctx context.Context, req request.SignUpInput) (res interface{}, err error) {
	return func(ctx context.Context, req request.SignUpInput) (res interface{}, err error) {
		res, err = h.AuthUsecase.Register(req)
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] Register usecase error: %v`, registerLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] Register successful for client: %s`, registerLogPrefix, res))
		return res, nil
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	req, err := DecodeSignUpRequest(ctx, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response.ErrorResponse{Message: err.Error()})
		return
	}

	res, err := RegisterEndpoint(h)(ctx, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response.ErrorResponse{Message: err.Error()})
		return
	}

	_, err = EncodeRegisterResponse(ctx, w, res)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] Register encode error: %v`, registerLogPrefix, err))
	}
}
