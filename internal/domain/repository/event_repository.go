package repository

import "github.com/hamsacumar/travel_backend_api/internal/domain/entity"

type EventRepository interface {
	Save(e *entity.Event) error
	LatestForRide(rideID string) (*entity.Event, error)
}
