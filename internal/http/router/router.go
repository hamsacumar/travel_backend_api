package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hamsacumar/travel_backend_api/internal/http/handler"
)

func SetupRouter(h *handler.Handler) *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

	//delete revoke token
	//delete expired otp
	//validation //validation for only one number for a role

	// Auth routes
	r.HandleFunc("/register", h.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	r.HandleFunc("/verify", h.VerifyOTP).Methods(http.MethodPost)
	r.HandleFunc("/send-otp", h.SendOTP).Methods(http.MethodPost)
	r.HandleFunc("/logout", h.Logout).Methods(http.MethodPost)

	return r
}

//func AuthMiddleware(secret string, tokenRepo repository.TokenRepository, next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		authz := r.Header.Get("Authorization")
//		if !strings.HasPrefix(authz, "Bearer ") {
//			http.Error(w, "missing bearer token", http.StatusUnauthorized)
//			return
//		}
//		raw := strings.TrimPrefix(authz, "Bearer ")
//
//		// Verify signature and basic claims
//		tok, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
//			if t.Method != jwt.SigningMethodHS256 {
//				return nil, fmt.Errorf("unexpected signing method")
//			}
//			return []byte(secret), nil
//		})
//		if err != nil || !tok.Valid {
//			http.Error(w, "invalid token", http.StatusUnauthorized)
//			return
//		}
//
//		// Check DB revocation
//		dbTok, err := tokenRepo.FindByToken(raw)
//		if err != nil || dbTok == nil || dbTok.Revoked {
//			http.Error(w, "token revoked", http.StatusUnauthorized)
//			return
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}
