package response

// ErrorResponse is returned when a request fails.
type ErrorResponse struct {
	Message string `json:"message"`
}

// MessageResponse is returned for simple success messages.
type MessageResponse struct {
	Message string `json:"message"`
}

// TokenResponse is returned after a successful login.
type TokenResponse struct {
	Token string `json:"token"`
}

// HealthResponse is returned by the health check endpoint.
type HealthResponse struct {
	Status string `json:"status"`
}
