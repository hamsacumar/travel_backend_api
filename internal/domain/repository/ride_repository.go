package repository

import "github.com/hamsacumar/travel_backend_api/internal/domain/entity"

type RideRepository interface {
	AddRide(ride *entity.Ride) error
}
