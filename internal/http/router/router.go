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

	// CORS handling: allow preflight and set headers
	r.Use(middleware.CORSMiddleware())
	// Advertise allowed methods for matched routes (for CORS)
	r.Use(mux.CORSMethodMiddleware(r))

	r.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

	//delete revoke token
	//delete expired otp
	//validation //validation for only one number for a role
	//bus number validation

	//ngnix - host
	//vendor delete before host
	//ip address - configure

	//--------------------------------------Auth routes------------------------------------------
	r.HandleFunc("/register", h.Register).Methods(http.MethodPost) //travel things
	r.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	r.HandleFunc("/verify", h.VerifyOTP).Methods(http.MethodPost)
	r.HandleFunc("/send-otp", h.SendOTP).Methods(http.MethodPost)
	r.HandleFunc("/logout", h.Logout).Methods(http.MethodPost)

	// Global OPTIONS handler (preflight). Must be after route registrations.
	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	//---------------------------------------Ride routes------------------------------------------
	//need to add metadata
	//need to validate driver and travel can't put same ride at a time
	r.Handle("/rides/driver", middleware.RoleAuthMiddleware([]string{"driver"}, http.HandlerFunc(rideHandler.AddRide))).Methods(http.MethodPost)
	r.Handle("/rides/travel", middleware.RoleAuthMiddleware([]string{"travel"}, http.HandlerFunc(travelHandler.TravelAddRide))).Methods(http.MethodPost)

	//---------------------------------------details routes----------------------------------------
	r.Handle("/driver_details", middleware.RoleAuthMiddleware([]string{"travel"}, http.HandlerFunc(detailHandler.GetDriverDetail))).Methods(http.MethodGet)
	r.Handle("/admin/driver_details", middleware.RoleAuthMiddleware([]string{"admin"}, http.HandlerFunc(adminDetailHandler.GetDriverDetailByAdmin))).Methods(http.MethodGet)
	r.Handle("/travel_details", middleware.RoleAuthMiddleware([]string{"admin"}, http.HandlerFunc(detailHandler.GetTravelDetail))).Methods(http.MethodGet)

	return r
}
