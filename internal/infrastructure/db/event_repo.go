package db

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/hamsacumar/travel_backend_api/internal/domain/entity"
)

type EventRepo struct {
	DB *sql.DB
}

func NewEventRepo(db *sql.DB) *EventRepo {
	return &EventRepo{DB: db}
}

func (r *EventRepo) Save(e *entity.Event) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now().UTC()
	}
	var infoJSON []byte
	var err error
	if e.Info != nil {
		infoJSON, err = json.Marshal(e.Info)
		if err != nil {
			return err
		}
	}
	_, err = r.DB.Exec(`
        INSERT INTO ride_event (ride_id, type, created_at, info)
        VALUES ($1, $2, $3, $4, $5)
    `, e.RideID, e.Type, e.CreatedAt, infoJSON)
	return err
}

func (r *EventRepo) LatestForRide(rideID string) (*entity.Event, error) {
	row := r.DB.QueryRow(`
        SELECT event_id, ride_id, type, created_at, info
        FROM ride_event
        WHERE ride_id = $1
        ORDER BY created_at DESC
        LIMIT 1
    `, rideID)

	var e entity.Event
	var infoJSON sql.NullString
	err := row.Scan(&e.EventID, &e.RideID, &e.Type, &e.CreatedAt, &infoJSON)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if infoJSON.Valid {
		var m map[string]interface{}
		_ = json.Unmarshal([]byte(infoJSON.String), &m)
		e.Info = m
	}
	return &e, nil
}
