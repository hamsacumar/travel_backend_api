package entity

import (
	"time"
)

// Token represents a persisted JWT token issued to a user
type Token struct {
	// Phone is the unique identifier for the user who owns this token
	Phone     string
	Role      string
	Token     string
	ExpiresAt *time.Time
	Revoked   bool
	CreatedAt time.Time
}
