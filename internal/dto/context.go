package dto

type ContextResponse struct {
	// Dialogue lines before the target quote
	Before []ParsedQuote `json:"before"`
	// The target quote
	Quote ParsedQuote `json:"quote"`
	// Dialogue lines after the target quote
	After []ParsedQuote `json:"after"`
}
