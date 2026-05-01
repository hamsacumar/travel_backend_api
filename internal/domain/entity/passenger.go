package entity

import (
	"time"
)

type Passenger struct {
	// ID is the primary key: a 6-digit zero-padded string (e.g., "042761")
	ID       string
	Username string
	Phone    string
	Email    string

	IsVerified bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
