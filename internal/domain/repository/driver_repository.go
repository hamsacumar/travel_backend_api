package repository

import "github.com/hamsacumar/travel_backend_api/internal/domain/entity"

type DriverRepository interface {
	Create(d entity.Driver) error
	FindByPhone(phone string) (*entity.Driver, error)
	Verify(phone string) error
}
