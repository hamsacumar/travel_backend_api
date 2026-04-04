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

const sendOTPLogPrefix = `travels-api.internal.http.handler.send_otp`

func DecodeSendOTPRequest(_ context.Context, r *http.Request) (request.SendOTPInput, error) {
	var req request.SendOTPInput
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] DecodeSendOTPRequest error: %v`, sendOTPLogPrefix, err))
		return req, err
	}
	return req, nil
}

func SendOTPEndpoint(h *Handler) func(ctx context.Context, req request.SendOTPInput) (res interface{}, err error) {
	return func(ctx context.Context, req request.SendOTPInput) (res interface{}, err error) {
		// Reuse Login usecase to generate and send OTP for a phone number
		res, err = h.AuthUsecase.SendOTP(req.Phone)
		if err != nil {
			log.Printf(fmt.Sprintf(`[%s] SendOTPEndpoint error: %v`, sendOTPLogPrefix, err))
			return nil, err
		}
		log.Printf(fmt.Sprintf(`[%s] OTP sent for phone: %s`, sendOTPLogPrefix, req.Phone))
		return res, nil
	}
}

func EncodeSendOTPResponse(_ context.Context, w http.ResponseWriter, res interface{}) (response interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] EncodeSendOTPResponse error: %v`, sendOTPLogPrefix, err))
		return nil, err
	}
	return res, nil
}

// SendOTP handles POST /send-otp
func (h *Handler) SendOTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	req, err := DecodeSendOTPRequest(ctx, r)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response.ErrorResponse{Message: err.Error()})
		return
	}

	res, err := SendOTPEndpoint(h)(ctx, req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response.ErrorResponse{Message: err.Error()})
		return
	}

	_, err = EncodeSendOTPResponse(ctx, w, res)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] SendOTP encode error: %v`, sendOTPLogPrefix, err))
	}
}
