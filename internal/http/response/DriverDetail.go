package response

type DriverDetailResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone_number"`
	BusName   string `json:"bus_name"`
	BusNumber string `json:"bus_number"`
	BusType   string `json:"bus_type"`
	SeatType  string `json:"bus_seat_type"`
}
