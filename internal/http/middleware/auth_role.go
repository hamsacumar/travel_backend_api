package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

const AuthMiddlewareLogPrefix = `travels-api.internal.http.middleware.auth_role`

// RoleAuthMiddleware returns a middleware that checks for one or more allowed roles in the JWT token.
func RoleAuthMiddleware(allowedRoles []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			// Log missing secret
			log.Printf("%s: JWT_SECRET not set in environment", AuthMiddlewareLogPrefix)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		// Extract Authorization Header and Checks for the presence of a Bearer
		authz := r.Header.Get("Authorization")
		if !strings.HasPrefix(authz, "Bearer ") {
			log.Printf("%s: missing bearer token", AuthMiddlewareLogPrefix)
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		// Removes the "Bearer " and Ensures the signing method and Uses the secret to validate
		raw := strings.TrimPrefix(authz, "Bearer ")
		tok, err := jwt.Parse(raw, func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				log.Printf("%s: unexpected signing method: %v", AuthMiddlewareLogPrefix, t.Method)
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(secret), nil
		})
		if err != nil || !tok.Valid {
			log.Printf("%s: invalid token: %v", AuthMiddlewareLogPrefix, err)
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Extracts the claims (payload) from the token
		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("%s: invalid claims structure", AuthMiddlewareLogPrefix)
			http.Error(w, "invalid claims", http.StatusUnauthorized)
			return
		}

		// Checks the role and user_id from the claims
		role, _ := claims["role"].(string)
		userID, _ := claims["user_id"].(string)
		allowed := false
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				allowed = true
				break
			}
		}
		if !allowed {
			log.Printf("%s: forbidden role: got '%s', want one of %v (user_id=%s)", AuthMiddlewareLogPrefix, role, allowedRoles, userID)
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		if userID == " " {
			log.Printf("%s:user id=%s is missing)", AuthMiddlewareLogPrefix, userID)
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		log.Printf("%s: authorized user_id=%s with role=%s", AuthMiddlewareLogPrefix, userID, role)

		// Sets the context with the role and user_id
		ctx := context.WithValue(r.Context(), "role", role)
		ctx = context.WithValue(ctx, "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
