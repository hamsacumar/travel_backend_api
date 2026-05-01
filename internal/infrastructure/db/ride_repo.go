package db

import (
	"database/sql"

	"github.com/hamsacumar/travel_backend_api/internal/domain/entity"
)

type RideRepo struct {
	DB *sql.DB
}

func NewRideRepo(db *sql.DB) *RideRepo {
	return &RideRepo{DB: db}
}

func (r *RideRepo) AddRide(ride *entity.Ride) error {
	_, err := r.DB.Exec(`INSERT INTO ride (ride_id, driver_id, start_lat, start_lon, end_lat, end_lon, date_of_journey, start_time, scheduled_by, seat_count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		ride.RideID,
		ride.DriverID,
		ride.StartLocation.Lat,
		ride.StartLocation.Lon,
		ride.EndLocation.Lat,
		ride.EndLocation.Lon,
		ride.DateOfJourney,
		ride.StartTime,
		ride.ScheduledBy,
		ride.SeatCount,
	)
	return err
}

func (r *RideRepo) FindByID(rideID string) (*entity.Ride, error) {
	row := r.DB.QueryRow(`SELECT ride_id, driver_id, start_lat, start_lon, end_lat, end_lon, date_of_journey, start_time, scheduled_by, seat_count FROM ride WHERE ride_id=$1`, rideID)
	var ride entity.Ride
	err := row.Scan(
		&ride.RideID,
		&ride.DriverID,
		&ride.StartLocation.Lat,
		&ride.StartLocation.Lon,
		&ride.EndLocation.Lat,
		&ride.EndLocation.Lon,
		&ride.DateOfJourney,
		&ride.StartTime,
		&ride.ScheduledBy,
		&ride.SeatCount,
	)
	if err != nil {
		return nil, err
	}
	return &ride, nil
}

// MoveRideToLog moves a ride from ride to ride_log in an idempotent transaction.
func (r *RideRepo) MoveRideToLog(rideID string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// Insert into ride_log only if not already present
	// Assumes ride_log has the same columns as ride
	insertSQL := `
        INSERT INTO ride_log (ride_id, driver_id, start_lat, start_lon, end_lat, end_lon, date_of_journey, start_time, scheduled_by, seat_count)
        SELECT r.ride_id, r.driver_id, r.start_lat, r.start_lon, r.end_lat, r.end_lon, r.date_of_journey, r.start_time, r.scheduled_by, r.seat_count
        FROM ride r
        WHERE r.ride_id = $1 AND NOT EXISTS (SELECT 1 FROM ride_log l WHERE l.ride_id = $1)
    `
	if _, err := tx.Exec(insertSQL, rideID); err != nil {
		return err
	}

	// Delete from ride table (safe even if 0 rows affected)
	if _, err := tx.Exec(`DELETE FROM ride WHERE ride_id = $1`, rideID); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
