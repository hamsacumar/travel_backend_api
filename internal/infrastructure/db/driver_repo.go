package db

import (
	"database/sql"

	"github.com/hamsacumar/travel_backend_api/internal/domain/entity"
)

type DriverRepo struct {
	DB *sql.DB
}

func NewDriverRepo(db *sql.DB) *DriverRepo {
	return &DriverRepo{DB: db}
}

func (r *DriverRepo) Create(d entity.Driver) error {
	query := `
        INSERT INTO drivers (
            username, phone, email,
            bus_name, bus_numbers, bus_type, seat_count,seat_type,
            travels_id, is_verified, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
    `
	_, err := r.DB.Exec(query,
		d.Username, d.Phone,
		d.Email,
		d.BusName, d.BusNumbers, d.BusType, d.SeatCount, d.SeatType,
		d.TravelsID, d.IsVerified,
	)
	return err
}

func (r *DriverRepo) FindByPhone(phone string) (*entity.Driver, error) {
	query := `
        SELECT id, username, phone, email,
            bus_name, bus_numbers, bus_type, seat_count,seat_type,
            travels_id, is_verified, created_at, updated_at
        FROM drivers WHERE phone = $1
    `
	row := r.DB.QueryRow(query, phone)

	var d entity.Driver
	err := row.Scan(
		&d.ID, &d.Username, &d.Phone,
		&d.Email,
		&d.BusName, &d.BusNumbers, &d.BusType, &d.SeatCount, &d.SeatType,
		&d.TravelsID, &d.IsVerified,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DriverRepo) Verify(phone string) error {
	query := `UPDATE drivers SET is_verified = true, updated_at = NOW() WHERE phone = $1`
	_, err := r.DB.Exec(query, phone)
	return err
}

func (r *DriverRepo) FindByPhoneAndTravel(phone string, travelID string) (*entity.Driver, error) {
	query := `
		SELECT id, username, phone, email,
			bus_name, bus_numbers, bus_type, seat_count,seat_type,
			travels_id, is_verified, created_at, updated_at
		FROM drivers WHERE phone = $1 AND travels_id = $2
	`
	row := r.DB.QueryRow(query, phone, travelID)
	var d entity.Driver
	err := row.Scan(
		&d.ID, &d.Username, &d.Phone,
		&d.Email,
		&d.BusName, &d.BusNumbers, &d.BusType, &d.SeatCount, &d.SeatType,
		&d.TravelsID, &d.IsVerified,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DriverRepo) FindByBusNumberAndTravel(busNumber string, travelID string) (*entity.Driver, error) {
	query := `
		SELECT id, username, phone, email,
			bus_name, bus_numbers, bus_type, seat_count,seat_type,
			travels_id, is_verified, created_at, updated_at
		FROM drivers WHERE bus_numbers = $1 AND travels_id = $2
	`
	row := r.DB.QueryRow(query, busNumber, travelID)
	var d entity.Driver
	err := row.Scan(
		&d.ID, &d.Username, &d.Phone,
		&d.Email,
		&d.BusName, &d.BusNumbers, &d.BusType, &d.SeatCount, d.SeatType,
		&d.TravelsID, &d.IsVerified,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DriverRepo) FindByBusNumber(busNumber string) (*entity.Driver, error) {
	query := `
        SELECT id, username, phone, email,
            bus_name, bus_numbers, bus_type, seat_count,seat_type,
            travels_id, is_verified, created_at, updated_at
        FROM drivers WHERE bus_numbers = $1
    `
	row := r.DB.QueryRow(query, busNumber)
	var d entity.Driver
	err := row.Scan(
		&d.ID, &d.Username, &d.Phone,
		&d.Email,
		&d.BusName, &d.BusNumbers, &d.BusType, &d.SeatCount, d.SeatType,
		&d.TravelsID, &d.IsVerified,
		&d.CreatedAt, &d.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DriverRepo) ExistsForTravel(driverID, travelID string) (bool, error) {
	query := `SELECT COUNT(1) FROM drivers WHERE id = $1 AND travels_id = $2`
	var count int
	err := r.DB.QueryRow(query, driverID, travelID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
