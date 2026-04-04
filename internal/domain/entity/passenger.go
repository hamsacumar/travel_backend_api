package entity

import (
	"time"

	"github.com/google/uuid"
)

type Passenger struct {
	ID       uuid.UUID
	Username string
	Phone    string
	Email    string

	IsVerified bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
