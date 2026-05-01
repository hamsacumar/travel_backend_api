package entity

import (
	"time"
)

type Driver struct {
	// ID is the primary key: a 6-digit zero-padded string (e.g., "042761")
	ID       string
	Username string
	Phone    string
	Email    string

	BusName    string
	BusNumbers string
	BusType    string
	SeatType   string

	TravelsID *string

	IsVerified bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
