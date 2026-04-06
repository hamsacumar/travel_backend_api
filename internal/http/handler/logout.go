package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hamsacumar/travel_backend_api/internal/http/response"
)

const logoutLogPrefix = `travels-api.internal.http.handler.logout`

// Logout revokes the current bearer token so it can no longer be used.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	authz := r.Header.Get("Authorization")
	if !strings.HasPrefix(authz, "Bearer ") {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response.ErrorResponse{Message: "missing bearer token"})
		return
	}

	token := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(response.ErrorResponse{Message: "empty bearer token"})
		return
	}

	if h.TokenRepo == nil {
		log.Printf(fmt.Sprintf("[%s] TokenRepo is nil", logoutLogPrefix))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response.ErrorResponse{Message: "server misconfiguration"})
		return
	}

	if err := h.TokenRepo.Revoke(token); err != nil {
		log.Printf(fmt.Sprintf("[%s] revoke error: %v", logoutLogPrefix, err))
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(response.ErrorResponse{Message: "failed to revoke token"})
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response.MessageResponse{Message: "logged out"})
}
