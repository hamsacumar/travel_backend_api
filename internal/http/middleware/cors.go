package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORSMiddleware adds CORS headers and handles preflight OPTIONS requests.
// Allowed origins:
// - http://localhost:4200
// - https://transitx-portal.web.app
// - value from env ALLOWED_CORS_ORIGINS (comma-separated list), if provided
func CORSMiddleware() func(http.Handler) http.Handler {
	// Build allowed origins set
	allowed := map[string]struct{}{
		"http://localhost:4200":           {},
		"https://transitx-portal.web.app": {},
	}
	if extra := strings.TrimSpace(os.Getenv("ALLOWED_CORS_ORIGINS")); extra != "" {
		parts := strings.Split(extra, ",")
		for _, p := range parts {
			o := strings.TrimSpace(p)
			if o != "" {
				allowed[o] = struct{}{}
			}
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" {
				if _, ok := allowed[origin]; ok {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					// Set credentials true to allow cookies/authorization headers when needed
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					w.Header().Set("Vary", "Origin")
				}
				// Always advertise allowed headers/methods for simplicity (browsers only use if origin matched)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, X-Correlation-Id")
			}

			// Short-circuit preflight
			if r.Method == http.MethodOptions {
				// No body for preflight
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
