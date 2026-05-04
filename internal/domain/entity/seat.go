package entity

import "time"

// Seat status constants
const (
	SeatOpen   = "seat_open"
	SeatClose  = "seat_close"
	SeatBook   = "seat_book"
	SeatPaid   = "seat_paid"
	SeatCheck  = "seat_check"
	SeatCancel = "seat_cancel"
)

// Seat represents a seat state for a given ride. Composite PK: (RideID, SeatNo)
type Seat struct {
	RideID    string
	SeatNo    int
	Status    string
	CreatedAt time.Time
	Info      map[string]interface{}
}
