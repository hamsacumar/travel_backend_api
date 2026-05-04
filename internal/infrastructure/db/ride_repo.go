package db

import (
	"database/sql"
	"log"

	"github.com/hamsacumar/travel_backend_api/internal/domain/entity"
)

type RideRepo struct {
	DB *sql.DB
}

func NewRideRepo(db *sql.DB) *RideRepo {
	return &RideRepo{DB: db}
}

func (r *RideRepo) AddRide(ride *entity.Ride) error {
	_, err := r.DB.Exec(`INSERT INTO ride (driver_id, start_lat, start_lon, end_lat, end_lon, date_of_journey, start_time,ticket_price,scheduled, scheduled_by,status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,$11)`,
		ride.DriverID,
		ride.StartLocation.Lat,
		ride.StartLocation.Lon,
		ride.EndLocation.Lat,
		ride.EndLocation.Lon,
		ride.DateOfJourney,
		ride.StartTime,
		ride.TicketPrice,
		ride.Scheduled,
		ride.ScheduledBy,
		entity.EventRideCreated,
	)
	return err
}

func (r *RideRepo) FindByID(rideID string) (*entity.Ride, error) {
	row := r.DB.QueryRow(`SELECT ride_id, driver_id, start_lat, start_lon, end_lat, end_lon, date_of_journey, start_time, scheduled_by FROM ride WHERE ride_id=$1`, rideID)
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
	)
	if err != nil {
		return nil, err
	}
	return &ride, nil
}

// MoveRideToLog moves a ride from ride to ride_log in an idempotent transaction.
func (r *RideRepo) MoveRideToLog(rideID string, status string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Step 1: Insert into ride_log (only if ride exists)
	insertQuery := `
	INSERT INTO ride_log (
		ride_id,
		driver_id,
		start_lat,
		start_lon,
		end_lat,
		end_lon,
		date_of_journey,
		start_time,
		ticket_price,
		scheduled,
		scheduled_by,
		created_at,
		status
	)
	SELECT
		ride_id,
		driver_id,
		start_lat,
		start_lon,
		end_lat,
		end_lon,
		date_of_journey,
		start_time,
		ticket_price,
		scheduled,
		scheduled_by,
		NOW(),
		$2
	FROM ride
	WHERE ride_id = $1;
	`

	res, err := tx.Exec(insertQuery, rideID, status)
	if err != nil {
		log.Printf("❌ insert into ride_log failed: %v", err)
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		log.Printf("⚠️ ride not found (already deleted?): %s", rideID)
		return tx.Commit() // NOT an error anymore
	}

	// Step 2: Delete from ride
	deleteQuery := `DELETE FROM ride WHERE ride_id = $1`
	_, err = tx.Exec(deleteQuery, rideID)
	if err != nil {
		log.Printf("❌ delete from ride failed: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Printf("✅ ride moved to log successfully: %s", rideID)
	return nil
}
