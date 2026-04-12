package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hamsacumar/travel_backend_api/internal/http/handler"
	"github.com/hamsacumar/travel_backend_api/internal/http/middleware"
)

func SetupRouter(
	h *handler.Handler,
	rideHandler *handler.RideHandler,
	travelHandler *handler.TravelRideHandler,
	detailHandler *handler.DetailHandler,
	adminDetailHandler *handler.AdminDetailHandler, // fixed type
) *mux.Router {

	r := mux.NewRouter()

	r.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

	//delete revoke token
	//delete expired otp
	//validation //validation for only one number for a role
	//bus number validation

	//ngnix - host
	//vendor delete before host

	//--------------------------------------Auth routes------------------------------------------
	r.HandleFunc("/register", h.Register).Methods(http.MethodPost) //bus register/ bus type / seat type
	r.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	r.HandleFunc("/verify", h.VerifyOTP).Methods(http.MethodPost)
	r.HandleFunc("/send-otp", h.SendOTP).Methods(http.MethodPost)
	r.HandleFunc("/logout", h.Logout).Methods(http.MethodPost)

	//---------------------------------------Ride routes------------------------------------------
	//need to add metadata
	r.Handle("/rides/driver", middleware.RoleAuthMiddleware([]string{"driver"}, http.HandlerFunc(rideHandler.AddRide))).Methods(http.MethodPost)
	r.Handle("/rides/travel", middleware.RoleAuthMiddleware([]string{"travel"}, http.HandlerFunc(travelHandler.TravelAddRide))).Methods(http.MethodPost)

	//---------------------------------------details routes----------------------------------------
	r.Handle("/driver_details", middleware.RoleAuthMiddleware([]string{"travel"}, http.HandlerFunc(detailHandler.GetDriverDetail))).Methods(http.MethodGet)
	r.Handle("/admin/driver_details", middleware.RoleAuthMiddleware([]string{"admin"}, http.HandlerFunc(adminDetailHandler.GetDriverDetailByAdmin))).Methods(http.MethodGet)
	r.Handle("/travel_details", middleware.RoleAuthMiddleware([]string{"admin"}, http.HandlerFunc(detailHandler.GetTravelDetail))).Methods(http.MethodGet)

	return r
}
