package repository

import "github.com/hamsacumar/travel_backend_api/internal/domain/entity"

type PassengerRepository interface {
	Create(p entity.Passenger) error
	FindByPhone(phone string) (*entity.Passenger, error)
	Verify(phone string) error
}
