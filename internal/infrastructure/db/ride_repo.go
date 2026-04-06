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
	_, err := r.DB.Exec(`INSERT INTO ride (ride_id, driver_id, start_lat, start_lon, end_lat, end_lon, date_of_journey, start_time, scheduled_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		ride.RideID,
		ride.DriverID,
		ride.StartLocation.Lat,
		ride.StartLocation.Lon,
		ride.EndLocation.Lat,
		ride.EndLocation.Lon,
		ride.DateOfJourney,
		ride.StartTime,
		ride.ScheduledBy,
	)
	return err
}
