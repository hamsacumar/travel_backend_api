package repository

import "github.com/hamsacumar/travel_backend_api/internal/domain/entity"

type OTPRepository interface {
	Save(o entity.OTP) error
	Find(phone, code string) (*entity.OTP, error)
}
