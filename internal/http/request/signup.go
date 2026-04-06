package request

type SignUpInput struct {
	Role string `json:"role"` // "passenger" | "driver"

	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`

	// Driver-specific fields (optional if passenger)
	BusName    string `json:"bus_name"`
	BusNumbers string `json:"bus_number"`
	BusType    string `json:"bus_type"`
	SeatType   string `json:"seat_type"`

	// Travel-specific fields (optional if passenger or driver because driver wants travels only he use it)
	TravelsName   *string   `json:"travels_name"`
	TravelsNumber *string   `json:"travels_number"`
	BusesNumbers  []*string `json:"BusesNumbers"`
}
