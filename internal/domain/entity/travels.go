package entity

import (
	"time"

	"github.com/google/uuid"
)

type Travels struct {
	ID           uuid.UUID
	Name         string
	Phone        string
	BusesNumbers []string

	IsVerified bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
