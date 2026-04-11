package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hamsacumar/travel_backend_api/internal/domain/entity"
	"github.com/hamsacumar/travel_backend_api/internal/domain/repository"
	"github.com/hamsacumar/travel_backend_api/internal/http/request"
)

const rideusecaseLogPrefix = `travels-api.internal.usecase.ride_usecase`

type RideUsecase struct {
	RideRepo   repository.RideRepository
	DriverRepo repository.DriverRepository // Add DriverRepo for validation
}

//func NewRideUsecase(rideRepo repository.RideRepository) *RideUsecase {
//	return &RideUsecase{RideRepo: rideRepo}
//}

func (u *RideUsecase) AddRide(req request.AddRideRequest, ctx context.Context) (*entity.Ride, error) {

	driverID, ok := ctx.Value("user_id").(string)
	if !ok || driverID == "" {
		log.Printf(fmt.Sprintf(`[%s] missing user_id in context`, rideusecaseLogPrefix))
		return nil, fmt.Errorf("missing user id")
	}

	ride := &entity.Ride{
		RideID:        uuid.NewString(),
		DriverID:      driverID,
		StartLocation: entity.Location{Lat: req.StartLocation.Lat, Lon: req.StartLocation.Lon},
		EndLocation:   entity.Location{Lat: req.EndLocation.Lat, Lon: req.EndLocation.Lon},
		DateOfJourney: req.DateOfJourney,
		StartTime:     req.StartTime,
		TicketPrice:   req.TicketPrice,
		Scheduled:     req.Scheduled,
		ScheduledBy:   "driver",
	}
	if err := u.RideRepo.AddRide(ride); err != nil {
		return nil, err
	}
	return ride, nil
}

func (u *RideUsecase) TravelAddRide(req request.TravelRideRequest, ctx context.Context) (*entity.Ride, error) {
	if req.DriverID == "" {
		log.Printf(fmt.Sprintf(`[%s] missing driver_id in request`, rideusecaseLogPrefix))
		return nil, fmt.Errorf("missing driver id in request")
	}

	//validate that driver belongs to the travel
	travelID, ok := ctx.Value("travel_id").(string)
	if !ok || travelID == "" {
		log.Printf(fmt.Sprintf(`[%s] missing travel_id in context`, rideusecaseLogPrefix))
		return nil, fmt.Errorf("missing travel id in context")
	}

	exists, err := u.DriverRepo.ExistsForTravel(req.DriverID, travelID)
	if err != nil {
		log.Printf(fmt.Sprintf(`[%s] error checking driver-travel association: %v`, rideusecaseLogPrefix, err))
		return nil, fmt.Errorf("error checking driver-travel association")
	}
	if !exists {
		log.Printf(fmt.Sprintf(`[%s] driver %s does not belong to travel %s`, rideusecaseLogPrefix, req.DriverID, travelID))
		return nil, fmt.Errorf("driver does not belong to this travel")
	}

	ride := &entity.Ride{
		RideID:        uuid.NewString(),
		DriverID:      req.DriverID,
		StartLocation: entity.Location{Lat: req.RideData.StartLocation.Lat, Lon: req.RideData.StartLocation.Lon},
		EndLocation:   entity.Location{Lat: req.RideData.EndLocation.Lat, Lon: req.RideData.EndLocation.Lon},
		DateOfJourney: req.RideData.DateOfJourney,
		StartTime:     req.RideData.StartTime,
		TicketPrice:   req.RideData.TicketPrice,
		Scheduled:     req.RideData.Scheduled,
		ScheduledBy:   "travel",
	}
	if err := u.RideRepo.AddRide(ride); err != nil {
		return nil, err
	}
	return ride, nil
}
