package db

import (
	"database/sql"

	"github.com/hamsacumar/travel_backend_api/internal/domain/entity"
)

type TravelsRepo struct {
	DB *sql.DB
}

func NewTravelsRepo(db *sql.DB) *TravelsRepo {
	return &TravelsRepo{DB: db}
}

func (r *TravelsRepo) Create(t entity.Travels) error {
	query := `
		INSERT INTO travels (id, name, phone, buses_numbers, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`
	_, err := r.DB.Exec(query,
		t.ID, t.Name, t.Phone, t.IsVerified,
	)
	return err
}

func (r *TravelsRepo) FindByPhone(phone string) (*entity.Travels, error) {
	query := `
		SELECT id, name, phone, buses_numbers, is_verified, created_at, updated_at
		FROM travels WHERE phone = $1
	`
	row := r.DB.QueryRow(query, phone)

	var t entity.Travels
	err := row.Scan(
		&t.ID, &t.Name, &t.Phone, &t.IsVerified,
		&t.CreatedAt, &t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TravelsRepo) Verify(phone string) error {
	query := `UPDATE travels SET is_verified = true, updated_at = NOW() WHERE phone = $1`
	_, err := r.DB.Exec(query, phone)
	return err
}
