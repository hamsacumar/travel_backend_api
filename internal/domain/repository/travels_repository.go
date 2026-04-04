package repository

import "github.com/hamsacumar/travel_backend_api/internal/domain/entity"

type TravelsRepository interface {
	Create(t entity.Travels) error
	FindByPhone(phone string) (*entity.Travels, error)
	Verify(phone string) error
}
