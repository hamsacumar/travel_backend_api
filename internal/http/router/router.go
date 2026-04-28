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
	//ip address - configure

	//--------------------------------------Auth routes------------------------------------------ //proper id instead of user_id
	r.HandleFunc("/register", h.Register).Methods(http.MethodPost) //travel things
	r.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	r.HandleFunc("/verify", h.VerifyOTP).Methods(http.MethodPost)
	r.HandleFunc("/send-otp", h.SendOTP).Methods(http.MethodPost)
	r.HandleFunc("/logout", h.Logout).Methods(http.MethodPost)

	//---------------------------------------Ride routes------------------------------------------
	//need to add metadata
	//(only my app should be master/don't use other app)
	r.Handle("/rides/driver", middleware.RoleAuthMiddleware([]string{"driver"}, http.HandlerFunc(rideHandler.AddRide))).Methods(http.MethodPost)
	r.Handle("/rides/travel", middleware.RoleAuthMiddleware([]string{"travel"}, http.HandlerFunc(travelHandler.TravelAddRide))).Methods(http.MethodPost)

	//---------------------------------------details routes----------------------------------------
	r.Handle("/driver_details", middleware.RoleAuthMiddleware([]string{"travel"}, http.HandlerFunc(detailHandler.GetDriverDetail))).Methods(http.MethodGet)
	r.Handle("/admin/driver_details", middleware.RoleAuthMiddleware([]string{"admin"}, http.HandlerFunc(adminDetailHandler.GetDriverDetailByAdmin))).Methods(http.MethodGet)
	r.Handle("/travel_details", middleware.RoleAuthMiddleware([]string{"admin"}, http.HandlerFunc(detailHandler.GetTravelDetail))).Methods(http.MethodGet)

	//---------------------------------------seat routes----------------------------------------
	//all bus details page for passenger //lock machanism
	//view the seat details
	//passenger book for them get passengerid from token
	//payment (paid)
	// bus details for particular travel
	//view the seat details
	//travel book for passenger (get passenger details/if passenger their no nedd/if not add in passenger table)
	//payment

	return r
}
