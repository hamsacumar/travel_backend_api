package request

type SetSeatStatusRequest struct {
	Status string                 `json:"status"`
	Info   map[string]interface{} `json:"info,omitempty"`
}
