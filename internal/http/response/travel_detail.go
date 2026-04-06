package response

type TravelDetailResponse struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Phone string   `json:"phone"`
	Buses []string `json:"buses_numbers"`
}
