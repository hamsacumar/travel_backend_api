package entity

import (
	"time"
)

type Travels struct {
	// ID is the primary key: a 6-digit zero-padded string (e.g., "042761")
	ID    string
	Email string
	Name  string
	Phone string

	IsVerified bool

	CreatedAt time.Time
	UpdatedAt time.Time
}
