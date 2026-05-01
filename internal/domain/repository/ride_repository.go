package repository

import "github.com/hamsacumar/travel_backend_api/internal/domain/entity"

type RideRepository interface {
	AddRide(ride *entity.Ride) error
	FindByID(rideID string) (*entity.Ride, error)
	// MoveRideToLog moves a ride row from ride to ride_log in an idempotent way
	// and deletes it from ride table. Should not error if ride is already moved or missing.
	MoveRideToLog(rideID string) error
}
