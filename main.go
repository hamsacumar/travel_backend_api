package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hamsacumar/travel_backend_api/adapter/repository"
	"github.com/hamsacumar/travel_backend_api/internal/http/handler"
	"github.com/hamsacumar/travel_backend_api/internal/http/router"
	dbRepo "github.com/hamsacumar/travel_backend_api/internal/infrastructure/db"
	"github.com/hamsacumar/travel_backend_api/internal/infrastructure/event"
	"github.com/hamsacumar/travel_backend_api/internal/infrastructure/service/auth"
	"github.com/hamsacumar/travel_backend_api/internal/usecase"
)

func main() {
	log.Println("Starting main()...")

	//environment variables
	//DB_HOST=dpg-d6vavdc50q8c739hmo2g-a.singapore-postgres.render.com
	//DB_PORT=5432
	//DB_USER=eztaxi_user
	//DB_PASS=kYFA1VbkpUkOHexLyOoUH1yZr8sqpqg2
	//DB_NAME=eztaxi
	//PORT=8080

	log.Println("Connecting to database...")
	// Initialize global DB
	repository.Connect()
	log.Println("Database connected.")

	// Wire up repositories
	passengerRepo := dbRepo.NewPassengerRepo(repository.DB)
	driverRepo := dbRepo.NewDriverRepo(repository.DB)
	travelsRepo := dbRepo.NewTravelsRepo(repository.DB)
	otpRepo := dbRepo.NewOTPRepo(repository.DB)
	tokenRepo := dbRepo.NewTokenRepo(repository.DB)
	log.Println("Repositories wired up.")

	// Wire up services
	jwtService := auth.NewJWTService(tokenRepo)
	log.Println("JWT service wired up.")

	// Wire up usecase
	authUsecase := usecase.NewAuthUsecase(passengerRepo, driverRepo, travelsRepo, otpRepo, jwtService)
	rideRepo := dbRepo.NewRideRepo(repository.DB)
	eventRepo := dbRepo.NewEventRepo(repository.DB)
	eventUC := event.NewEventUsecase(eventRepo, rideRepo)
	rideUsecase := &usecase.RideUsecase{RideRepo: rideRepo, DriverRepo: driverRepo, EventUC: eventUC}
	detailUsecase := &usecase.DetailUsecase{DriverRepo: driverRepo, TravelRepo: travelsRepo}
	log.Println("Usecases wired up.")

	// Wire up handler
	h := &handler.Handler{AuthUsecase: authUsecase, TokenRepo: tokenRepo}
	rideHandler := &handler.RideHandler{RideUsecase: rideUsecase}
	travelHandler := &handler.TravelRideHandler{RideUsecase: rideUsecase}
	eventHandler := &handler.EventHandler{EventUC: eventUC}
	detailHandler := &handler.DetailHandler{DetailUsecase: detailUsecase}
	log.Println("Handlers wired up.")

	// Setup router
	r := router.SetupRouter(h, rideHandler, travelHandler, eventHandler, detailHandler)
	//adminDetailHandler

	log.Println("Router set up.")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal(err)
	}
}
