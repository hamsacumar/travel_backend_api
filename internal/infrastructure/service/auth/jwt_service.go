package auth

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (j *JWTService) GenerateToken(userID string, role string) (token string, err error) {
	log.Printf(fmt.Sprintf(`[%s] GenerateToken started for userID: %s role: %s`, jwtLogPrefix, userID, role))

	if userID == "" {
		log.Printf(fmt.Sprintf(`[%s] GenerateToken error: userID is empty`, jwtLogPrefix))
		return "", fmt.Errorf("userID is required")
	}

	if role == "" {
		log.Printf(fmt.Sprintf(`[%s] GenerateToken error: role is empty`, jwtLogPrefix))
		return "", fmt.Errorf("role is required")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * 180 * time.Hour).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = t.SignedString([]byte(j.Secret))
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] GenerateToken sign error: %v`, jwtLogPrefix, err))
		return "", err
	}

	// Persist the token to DB
	// Parse userID string to UUID
	uid, parseErr := uuid.Parse(userID)
	if parseErr != nil {
		log.Printf(fmt.Sprintf(`[%s] GenerateToken parse userID error: %v`, jwtLogPrefix, parseErr))
		return "", parseErr
	}

	tok := entity.Token{
		ID:        uuid.New(),
		UserID:    uid,
		Role:      role,
		Token:     token,
		ExpiresAt: nil,
		Revoked:   false, // forcefully revoke (expired) the token
		CreatedAt: time.Now(),
	}

	if err := j.tokenRepo.Save(tok); err != nil {
		log.Printf(fmt.Sprintf(`[%s] GenerateToken save error: %v`, jwtLogPrefix, err))
		return "", err
	}

	log.Printf(fmt.Sprintf(`[%s] GenerateToken successful and saved for userID: %s`, jwtLogPrefix, userID))
	return token, nil
}
