package entity

import "time"

// EventType defines all allowed event types for a ride lifecycle
const (
	EventRideCreated   = "ride_created"
	EventTripCreated   = "trip_created"
	EventTripStarted   = "trip_started"
	EventTripCompleted = "trip_completed"
	EventRideCancelled = "ride_cancelled"
	EventRideDeleted   = "ride_deleted"
)

// Event represents a domain event for a ride
type Event struct {
	EventID   string                 // uuid
	RideID    string                 // associated ride id
	Type      string                 // one of the constants above
	CreatedAt time.Time              // event time
	Info      map[string]interface{} // optional metadata (actor, reason, etc.)
}
