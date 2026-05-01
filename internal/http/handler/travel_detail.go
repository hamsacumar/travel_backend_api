package handler

//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"log"
//	"net/http"
//)
//
//// Decoder
//func DecodeGetTravelDetailRequest(_ context.Context, r *http.Request) (string, error) {
//	phone := r.URL.Query().Get("phone_number")
//	if phone == "" {
//		return "", fmt.Errorf("phone_number is required")
//	}
//	return phone, nil
//}
//
//// Encoder
//func EncodeGetTravelDetailResponse(_ context.Context, w http.ResponseWriter, res interface{}, status int) (interface{}, error) {
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(status)
//	err := json.NewEncoder(w).Encode(res)
//	if err != nil {
//		log.Printf(fmt.Sprintf(`[travels-api.internal.http.handler.travel_detail] EncodeGetTravelDetailResponse Error: %v`, err))
//		return nil, err
//	}
//	return res, nil
//}
//
//// Endpoint
//func GetTravelDetailEndpoint(h *DetailHandler) func(ctx context.Context, phone string) (interface{}, error) {
//	return func(ctx context.Context, phone string) (interface{}, error) {
//		resp, err := h.DetailUsecase.GetTravelDetailByPhone(ctx, phone)
//		if err != nil {
//			return nil, err
//		}
//		return resp, nil
//	}
//}
//
//// Handler
//func (h *DetailHandler) GetTravelDetail(w http.ResponseWriter, r *http.Request) {
//	ctx := r.Context()
//
//	phone, err := DecodeGetTravelDetailRequest(ctx, r)
//	if err != nil {
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusBadRequest)
//		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
//		return
//	}
//
//	res, err := GetTravelDetailEndpoint(h)(ctx, phone)
//	if err != nil {
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusNotFound)
//		_ = json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
//		return
//	}
//
//	_, err = EncodeGetTravelDetailResponse(ctx, w, res, http.StatusOK)
//	if err != nil {
//		log.Printf(fmt.Sprintf(`[travels-api.internal.http.handler.travel_detail] GetTravelDetail encode error: %v`, err))
//	}
//}
