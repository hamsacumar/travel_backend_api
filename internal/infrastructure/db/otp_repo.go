package db

import (
	"database/sql"

	"github.com/hamsacumar/travel_backend_api/internal/domain/entity"
)

type OTPRepo struct {
	DB *sql.DB
}

func NewOTPRepo(db *sql.DB) *OTPRepo {
	return &OTPRepo{DB: db}
}

func (r *OTPRepo) Save(o entity.OTP) error {
	deleteQuery := `DELETE FROM otps WHERE phone = $1`
	_, err := r.DB.Exec(deleteQuery, o.Phone)
	if err != nil {
		return err
	}

	insertQuery := `
        INSERT INTO otps (phone, code, expires_at, verified)
        VALUES ($1, $2, $3, $4)
    `
	_, err = r.DB.Exec(insertQuery,
		o.Phone, o.Code, o.ExpiresAt, o.Verified,
	)
	return err
}

func (r *OTPRepo) Find(phone, code string) (*entity.OTP, error) {
	query := `
		SELECT id, phone, code, expires_at, verified
		FROM otps WHERE phone = $1 AND code = $2
	`
	row := r.DB.QueryRow(query, phone, code)

	var o entity.OTP
	err := row.Scan(&o.ID, &o.Phone, &o.Code, &o.ExpiresAt, &o.Verified)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &o, nil
}
