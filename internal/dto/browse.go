package dto

type CharacterResponse struct {
	// Numeric character identifier
	CharacterID string `json:"characterId" example:"10"`
	// Display name of the character
	Character string `json:"character" example:"Ushiromiya Battler"`
	// Paginated list of quotes
	Quotes any `json:"quotes" swaggertype:"array,object"`
	// Total number of matching quotes
	Total int `json:"total" example:"4304"`
	// Maximum results per page
	Limit int `json:"limit" example:"50"`
	// Pagination offset
	Offset int `json:"offset" example:"0"`
}
