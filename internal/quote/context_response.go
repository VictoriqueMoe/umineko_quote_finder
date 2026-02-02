package quote

type ContextResponse struct {
	Before []ParsedQuote `json:"before"`
	Quote  ParsedQuote   `json:"quote"`
	After  []ParsedQuote `json:"after"`
}
