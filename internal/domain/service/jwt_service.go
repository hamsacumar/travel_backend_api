package service

type JWTService interface {
	// GenerateToken issues a JWT where `user_id` is the 6-digit user ID, and persists the token keyed by phone.
	// Before saving, any existing token for that phone is deleted (single active token per phone).
	GenerateToken(userID string, phone string, role string) (string, error)
}
