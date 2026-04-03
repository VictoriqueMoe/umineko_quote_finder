package dto

type ContextResponse struct {
	// Dialogue lines before the target quote
	Before any `json:"before" swaggertype:"array,object"`
	// The target quote
	Quote any `json:"quote" swaggertype:"object"`
	// Dialogue lines after the target quote
	After any `json:"after" swaggertype:"array,object"`
}
