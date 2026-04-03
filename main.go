package main

import (
    "log"
    "net/http"
    "os"

    "github.com/hamsacumar/travel_backend_api/adapter/repository"
    "github.com/hamsacumar/travel_backend_api/internal/http/handler"
    "github.com/hamsacumar/travel_backend_api/internal/http/router"
    dbRepo "github.com/hamsacumar/travel_backend_api/internal/infrastructure/db"
    "github.com/hamsacumar/travel_backend_api/internal/infrastructure/service/auth"
    "github.com/hamsacumar/travel_backend_api/internal/usecase"
)

func main() {

	//environment variables
	//DB_HOST=dpg-d6vavdc50q8c739hmo2g-a.singapore-postgres.render.com
	//DB_PORT=5432
	//DB_USER=eztaxi_user
	//DB_PASS=kYFA1VbkpUkOHexLyOoUH1yZr8sqpqg2
	//DB_NAME=eztaxi
	//PORT=8080

	// Initialize global DB
	repository.Connect()

 // Wire up repositories
 passengerRepo := dbRepo.NewPassengerRepo(repository.DB)
 driverRepo := dbRepo.NewDriverRepo(repository.DB)
 travelsRepo := dbRepo.NewTravelsRepo(repository.DB)
 otpRepo := dbRepo.NewOTPRepo(repository.DB)
 tokenRepo := dbRepo.NewTokenRepo(repository.DB)

 // Wire up services
 jwtService := auth.NewJWTService(tokenRepo)

 // Wire up usecase
 authUsecase := usecase.NewAuthUsecase(passengerRepo, driverRepo, travelsRepo, otpRepo, jwtService)

 // Wire up handler
 h := &handler.Handler{AuthUsecase: authUsecase, TokenRepo: tokenRepo}

	// Setup router
	r := router.SetupRouter(h)

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
