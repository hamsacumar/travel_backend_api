package repository

import "github.com/hamsacumar/travel_backend_api/internal/domain/entity"

type TokenRepository interface {
    Save(t entity.Token) error
    FindByToken(token string) (*entity.Token, error)
    Revoke(token string) error
}
