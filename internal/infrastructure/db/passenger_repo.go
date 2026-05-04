package db

import (
	"database/sql"

	"github.com/hamsacumar/travel_backend_api/internal/domain/entity"
)

type PassengerRepo struct {
	DB *sql.DB
}

func NewPassengerRepo(db *sql.DB) *PassengerRepo {
	return &PassengerRepo{DB: db}
}

func (r *PassengerRepo) Create(p entity.Passenger) error {
	query := `
        INSERT INTO passengers (username, phone, email, is_verified, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
    `
	_, err := r.DB.Exec(query,
		p.Username, p.Phone,
		p.Email, p.IsVerified,
	)
	return err
}

func (r *PassengerRepo) FindByPhone(phone string) (*entity.Passenger, error) {
	query := `
        SELECT id, username, phone, email, is_verified, created_at, updated_at
        FROM passengers WHERE phone = $1
    `
	row := r.DB.QueryRow(query, phone)

	var p entity.Passenger
	err := row.Scan(
		&p.ID, &p.Username, &p.Phone, &p.Email, &p.IsVerified,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PassengerRepo) Verify(phone string) error {
	query := `UPDATE passengers SET is_verified = true, updated_at = NOW() WHERE phone = $1`
	_, err := r.DB.Exec(query, phone)
	return err
}
