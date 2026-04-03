package entity

import (
	"time"

	"github.com/google/uuid"
)

type Driver struct {
	ID       uuid.UUID
	Username string
	Phone    string
	Email    string

	BusName    string
	BusNumbers string
	BusType    string
	SeatType   string

	TravelsID string

	IsVerified bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
