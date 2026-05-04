package event

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hamsacumar/travel_backend_api/internal/domain/entity"
	"github.com/hamsacumar/travel_backend_api/internal/domain/repository"
)

const rideeventusecaseLogPrefix = `travels-api.internal.infra.ride_event_usecase`

type EventUsecase struct {
	EventRepo repository.EventRepository
	RideRepo  repository.RideRepository
}

func NewEventUsecase(eventRepo repository.EventRepository, rideRepo repository.RideRepository) *EventUsecase {
	return &EventUsecase{EventRepo: eventRepo, RideRepo: rideRepo}
}

// failAndDelete attempts to publish a ride_deleted event with details about a failed event
// and the underlying error. This method must never panic; any internal error is ignored.
func (u *EventUsecase) failAndDelete(ctx context.Context, rideID, failedEvent string, cause error, extra map[string]interface{}) {
	info := map[string]interface{}{
		"auto":         true,
		"reason":       "failed_event",
		"failed_event": failedEvent,
	}

	if cause != nil {
		info["error"] = cause.Error()
	}

	for k, v := range extra {
		info[k] = v
	}

	log.Printf("CRITICAL FAILURE → deleting ride %s, reason: %v", rideID, info)

	_ = u.PublishRideDeleted(ctx, rideID, info)
}

// PublishRideCreated should be called after a ride is added
func (u *EventUsecase) PublishRideCreated(ctx context.Context, rideID string, info map[string]interface{}) error {
	// Persist ride_created
	if err := u.publish(ctx, rideID, entity.EventRideCreated, info); err != nil {
		log.Printf(fmt.Sprintf(`failed to publish ride_created: %s`, rideeventusecaseLogPrefix, err))
		// On any failure, publish ride_deleted with error info instead of returning error
		u.failAndDelete(context.Background(), rideID, entity.EventRideCreated, err, info)
		return err
	}

	// Automatically schedule trip_created at the ride's configured StartTime (DateOfJourney + StartTime)
	ride, err := u.RideRepo.FindByID(rideID)
	if err != nil {
		u.failAndDelete(context.Background(), rideID, entity.EventTripCreated, err, map[string]interface{}{"phase": "fetch_ride_for_schedule"})
		return err
	}

	layoutDate := "2006-01-02"
	layoutTime := "15:04:05"
	d, err := time.Parse(layoutDate, ride.DateOfJourney)
	if err != nil {
		log.Printf("invalid date_of_journey: %v", ride.DateOfJourney)
		return err
	}
	t, err := time.Parse(layoutTime, ride.StartTime)
	if err != nil {
		log.Printf("invalid start_time: %v", ride.StartTime)
		return err
	}
	triggerAt := time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)

	delay := time.Until(triggerAt)
	if delay <= 0 {
		log.Printf("start time already passed for ride %s", rideID)
		return err
	}
	go func() {
		timer := time.NewTimer(delay)
		<-timer.C
		if err := u.PublishTripCreated(context.Background(), rideID, map[string]interface{}{"auto": true}); err != nil {
			u.failAndDelete(context.Background(), rideID, entity.EventTripCreated, err, map[string]interface{}{"auto": true})
		}
	}()
	return err
}

func (u *EventUsecase) PublishTripCreated(ctx context.Context, rideID string, info map[string]interface{}) error {
	// Disallow duplicate or out-of-order events
	if err := u.validateTransition(rideID, entity.EventTripCreated); err != nil {
		u.failAndDelete(context.Background(), rideID, entity.EventTripCreated, err, info)
		return err
	}
	// Just publish trip_created; trip_started is manual per latest rule
	if err := u.publish(ctx, rideID, entity.EventTripCreated, info); err != nil {
		u.failAndDelete(context.Background(), rideID, entity.EventTripCreated, err, info)
		return err
	}
	return nil
}

func (u *EventUsecase) PublishTripCompleted(ctx context.Context, rideID string, info map[string]interface{}) error {
	if err := u.transitionAndPublish(ctx, rideID, entity.EventTripCompleted, info); err != nil {
		u.failAndDelete(context.Background(), rideID, entity.EventTripCompleted, err, info)
		return err
	}
	// move the ride to ride_log after completion
	if u.RideRepo != nil {
		if err := u.RideRepo.MoveRideToLog(rideID, entity.EventRideDeleted); err != nil {
			log.Printf("FAILED to move ride to log: %v", err)
			return err
		}
		log.Printf(fmt.Sprintf(`move to ride to ride_log`, rideeventusecaseLogPrefix))
	}
	return nil
}

func (u *EventUsecase) PublishTripStarted(ctx context.Context, rideID string, info map[string]interface{}) error {
	if err := u.transitionAndPublish(ctx, rideID, entity.EventTripStarted, info); err != nil {
		u.failAndDelete(context.Background(), rideID, entity.EventTripStarted, err, info)
		return err
	}
	return nil
}

func (u *EventUsecase) PublishRideCancelled(ctx context.Context, rideID string, info map[string]interface{}) error {
	// Cancel is not allowed after trip_created
	latest, err := u.EventRepo.LatestForRide(rideID)
	if err != nil {
		u.failAndDelete(context.Background(), rideID, entity.EventRideCancelled, err, info)
		return err
	}
	if latest != nil && (latest.Type == entity.EventTripCreated || latest.Type == entity.EventTripStarted || latest.Type == entity.EventTripCompleted) {
		// Business rule violation → delete ride instead of returning error
		u.failAndDelete(context.Background(), rideID, entity.EventRideCancelled, fmt.Errorf("ride_cancelled not allowed after trip_created"), info)
		return err
	}
	if err := u.transitionAndPublish(ctx, rideID, entity.EventRideCancelled, info); err != nil {
		u.failAndDelete(context.Background(), rideID, entity.EventRideCancelled, err, info)
		return err
	}
	// move the ride to ride_log after cancellation
	if u.RideRepo != nil {
		if err := u.RideRepo.MoveRideToLog(rideID, entity.EventRideDeleted); err != nil {
			log.Printf("FAILED to move ride to log: %v", err)
			return err
		}
		log.Printf(fmt.Sprintf(`move to ride to ride_log`, rideeventusecaseLogPrefix))
		return nil
	}
	return nil
}

func (u *EventUsecase) PublishRideDeleted(ctx context.Context, rideID string, info map[string]interface{}) error {
	// ride_deleted can happen anytime
	if err := u.publish(ctx, rideID, entity.EventRideDeleted, info); err != nil {
		// Final fallback: if even deletion cannot be recorded, swallow error
		return err
	}
	// move the ride to ride_log after deletion
	if u.RideRepo != nil {
		log.Printf(fmt.Sprintf(`move to ride to ride_log`, rideeventusecaseLogPrefix))
		if err := u.RideRepo.MoveRideToLog(rideID, entity.EventRideDeleted); err != nil {
			log.Printf("FAILED to move ride to log: %v", err)
			return err
		}
	}
	return nil
}

// Helpers
func (u *EventUsecase) transitionAndPublish(ctx context.Context, rideID string, next string, info map[string]interface{}) error {
	if err := u.validateTransition(rideID, next); err != nil {
		return err
	}
	return u.publish(ctx, rideID, next, info)
}

func (u *EventUsecase) publish(_ context.Context, rideID string, eventType string, info map[string]interface{}) error {
	e := &entity.Event{
		RideID:    rideID,
		Type:      eventType,
		CreatedAt: time.Now().UTC(),
		Info:      info,
	}
	log.Printf(fmt.Sprintf(`event: %s`, rideeventusecaseLogPrefix, e))
	return u.EventRepo.Save(e)
}

func (u *EventUsecase) validateTransition(rideID string, next string) error {
	latest, err := u.EventRepo.LatestForRide(rideID)
	if err != nil {
		return err
	}
	if latest == nil {
		// Only ride_created is valid as first event
		if next != entity.EventRideCreated && next != entity.EventRideDeleted {
			return fmt.Errorf("first event must be ride_created or ride_deleted")
		}
		return nil
	}

	// Define allowed sequence
	order := map[string]int{
		entity.EventRideCreated:   1,
		entity.EventTripCreated:   2,
		entity.EventTripStarted:   3,
		entity.EventTripCompleted: 4,
		entity.EventRideCancelled: 5, // special; only allowed before trip_created by higher-level checks
		entity.EventRideDeleted:   6, // can happen anytime; skip order check
	}

	if next == entity.EventRideDeleted {
		return nil
	}

	ln, okL := order[latest.Type]
	nn, okN := order[next]
	if !okL || !okN {
		return fmt.Errorf("unknown event type")
	}
	if nn <= ln {
		return fmt.Errorf("invalid event order: %s -> %s", latest.Type, next)
	}
	return nil
}
