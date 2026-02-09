package dto

type ParsedQuote struct {
	// Plain text of the quote
	Text string `json:"text" example:"You're insane!! Stop it already, damn iiiiiiiiit!!"`
	// HTML-formatted text with truth coloring
	TextHtml string `json:"textHtml" example:"You&#39;re insane!! Stop it already, damn iiiiiiiiit!!"`
	// Numeric character identifier
	CharacterID string `json:"characterId" example:"10"`
	// Display name of the character
	Character string `json:"character" example:"Ushiromiya Battler"`
	// Comma-separated audio file IDs
	AudioID string `json:"audioId" example:"30101088, 30101089"`
	// Maps each audio ID to its character ID
	AudioCharMap map[string]string `json:"audioCharMap,omitempty" example:"30101088:10,30101089:10"`
	// Maps each audio ID to its spoken text
	AudioTextMap map[string]string `json:"audioTextMap,omitempty" example:"30101088:You're insane!!,30101089:Stop it already"`
	// Episode number (1-8)
	Episode int `json:"episode" example:"3"`
	// Content type marker
	ContentType string `json:"contentType" example:""`
	// Whether the quote contains red truth
	HasRedTruth bool `json:"hasRedTruth,omitempty" example:"false"`
	// Whether the quote contains blue truth
	HasBlueTruth bool `json:"hasBlueTruth,omitempty" example:"false"`
}
