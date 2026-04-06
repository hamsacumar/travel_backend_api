package entity

import (
	"time"

	"github.com/google/uuid"
)

// Token represents a persisted JWT token issued to a user
type Token struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Role      string
	Token     string
	ExpiresAt *time.Time
	Revoked   bool
	CreatedAt time.Time
}
