package request

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type AddRideRequest struct {
	StartLocation Location `json:"start_location"`
	EndLocation   Location `json:"end_location"`
	DateOfJourney string   `json:"date_of_journey"` // Should be in YYYY-MM-DD format
	StartTime     string   `json:"start_time"`      // Should be in HH:MM:SS format
}

type TravelRideRequest struct {
	DriverID string
	RideData AddRideRequest
}
