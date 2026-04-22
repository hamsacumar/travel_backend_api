package entity

type Location struct {
	Lat float64
	Lon float64
}

type Ride struct {
	RideID        string
	DriverID      string
	StartLocation Location
	EndLocation   Location
	DateOfJourney string // Should be in YYYY-MM-DD format
	StartTime     string // Should be in HH:MM:SS format
	TicketPrice   float64
	Scheduled     string
	ScheduledBy   string // driver/travel/admin
}
