package entity

import (
	"time"

	"github.com/google/uuid"
)

type Travels struct {
	ID    uuid.UUID
	Email string
	Name  string
	Phone string

	IsVerified bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
