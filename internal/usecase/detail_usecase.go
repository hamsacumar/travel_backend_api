package usecase

import (
	"context"
	"fmt"

	"github.com/hamsacumar/travel_backend_api/internal/domain/entity"
	"github.com/hamsacumar/travel_backend_api/internal/domain/repository"
	"github.com/hamsacumar/travel_backend_api/internal/http/request"
	"github.com/hamsacumar/travel_backend_api/internal/http/response"
)

type DetailUsecase struct {
	DriverRepo repository.DriverRepository
	TravelRepo repository.TravelsRepository
}

func (u *DetailUsecase) GetDriverDetailByTravel(ctx context.Context, req request.DriverDetailRequest) (response.DriverDetailResponse, error) {
	travelID := ctx.Value("travel_id")
	if travelID == nil {
		return response.DriverDetailResponse{}, fmt.Errorf("unauthorized: travel_id missing")
	}
	var driver *entity.Driver
	var err error
	if req.Phone != "" {
		driver, err = u.DriverRepo.FindByPhoneAndTravel(req.Phone, travelID.(string))
	} else {
		driver, err = u.DriverRepo.FindByBusNumberAndTravel(req.BusNumber, travelID.(string))
	}
	if err != nil {
		return response.DriverDetailResponse{}, fmt.Errorf("error fetching driver details")
	}
	if driver == nil {
		return response.DriverDetailResponse{}, fmt.Errorf("driver not found")
	}
	resp := response.DriverDetailResponse{
		ID:        driver.ID.String(),
		Username:  driver.Username,
		Email:     driver.Email,
		Phone:     driver.Phone,
		BusName:   driver.BusName,
		BusNumber: driver.BusNumbers,
		BusType:   driver.BusType,
		SeatType:  driver.SeatType,
	}
	return resp, nil
}

func (u *DetailUsecase) GetTravelDetailByPhone(ctx context.Context, phone string) (response.TravelDetailResponse, error) {
	travelRepo, ok := u.TravelRepo.(repository.TravelsRepository)
	if !ok {
		return response.TravelDetailResponse{}, fmt.Errorf("travel repository not available")
	}
	travel, err := travelRepo.FindByPhone(phone)
	if err != nil {
		return response.TravelDetailResponse{}, fmt.Errorf("error fetching travel details")
	}
	if travel == nil {
		return response.TravelDetailResponse{}, fmt.Errorf("travel not found")
	}
	resp := response.TravelDetailResponse{
		ID:    travel.ID.String(),
		Name:  travel.Name,
		Phone: travel.Phone,
	}
	return resp, nil
}

func (u *DetailUsecase) GetDriverDetailByAdmin(req request.DriverDetailRequest) (response.DriverDetailResponse, error) {
	var driver *entity.Driver
	var err error
	if req.Phone != "" {
		driver, err = u.DriverRepo.FindByPhone(req.Phone)
	} else {
		driver, err = u.DriverRepo.FindByBusNumber(req.BusNumber)
	}
	if err != nil {
		return response.DriverDetailResponse{}, fmt.Errorf("error fetching driver details")
	}
	if driver == nil {
		return response.DriverDetailResponse{}, fmt.Errorf("driver not found")
	}
	resp := response.DriverDetailResponse{
		ID:        driver.ID.String(),
		Username:  driver.Username,
		Email:     driver.Email,
		Phone:     driver.Phone,
		BusName:   driver.BusName,
		BusNumber: driver.BusNumbers,
		BusType:   driver.BusType,
		SeatType:  driver.SeatType,
	}
	return resp, nil
}
