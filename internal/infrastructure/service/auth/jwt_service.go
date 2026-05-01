package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hamsacumar/travel_backend_api/internal/domain/entity"
	"github.com/hamsacumar/travel_backend_api/internal/domain/repository"
)

const jwtLogPrefix = `travels-api.internal.infrastructure.service.auth`

type JWTService struct {
	Secret    string
	tokenRepo repository.TokenRepository
}

func NewJWTService(tokenRepo repository.TokenRepository) *JWTService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Printf(fmt.Sprintf(`[%s] JWT_SECRET not set, using default secret`, jwtLogPrefix))
		secret = "default_secret"
	}
	log.Printf(fmt.Sprintf(`[%s] JWTService initialized`, jwtLogPrefix))
	return &JWTService{Secret: secret, tokenRepo: tokenRepo}
}

// GenerateToken issues a JWT where user_id is the 6-digit user ID, and persists the token keyed by phone.
// Before saving, any existing token for that phone is deleted (single active token per phone).
func (j *JWTService) GenerateToken(userID string, phone string, role string) (token string, err error) {
	log.Printf(fmt.Sprintf(`[%s] GenerateToken started for userID: %s phone: %s role: %s`, jwtLogPrefix, userID, phone, role))

	if userID == "" {
		log.Printf(fmt.Sprintf(`[%s] GenerateToken error: userID is empty`, jwtLogPrefix))
		return "", fmt.Errorf("userID is required")
	}
	if len(userID) != 6 {
		log.Printf(fmt.Sprintf(`[%s] GenerateToken error: userID must be 6 digits`, jwtLogPrefix))
		return "", fmt.Errorf("userID must be 6 digits")
	}
	if phone == "" {
		log.Printf(fmt.Sprintf(`[%s] GenerateToken error: phone is empty`, jwtLogPrefix))
		return "", fmt.Errorf("phone is required")
	}
	if role == "" {
		log.Printf(fmt.Sprintf(`[%s] GenerateToken error: role is empty`, jwtLogPrefix))
		return "", fmt.Errorf("role is required")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"phone":   phone,
		"role":    role,
		"exp":     time.Now().Add(24 * 180 * time.Hour).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = t.SignedString([]byte(j.Secret))
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] GenerateToken sign error: %v`, jwtLogPrefix, err))
		return "", err
	}

	tok := entity.Token{
		Phone:     phone,
		Role:      role,
		Token:     token,
		ExpiresAt: nil,
		Revoked:   false,
		CreatedAt: time.Now(),
	}
	if err := j.tokenRepo.Save(tok); err != nil {
		log.Printf(fmt.Sprintf(`[%s] GenerateToken save error: %v`, jwtLogPrefix, err))
		return "", err
	}
	log.Printf(fmt.Sprintf(`[%s] GenerateToken successful and saved for phone: %s`, jwtLogPrefix, phone))

	return token, nil
}
