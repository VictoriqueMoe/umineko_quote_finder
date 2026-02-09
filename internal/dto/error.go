package dto

type ErrorResponse struct {
	// Error message
	Error string `json:"error" example:"query parameter 'q' is required"`
}
