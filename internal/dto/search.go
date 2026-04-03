package dto

type (
	SearchResult struct {
		// The matched quote
		Quote any `json:"quote" swaggertype:"object"`
		// Relevance score (100 = exact match)
		Score int `json:"score" example:"100"`
	}

	SearchResponse struct {
		Results []SearchResult `json:"results"`
		Total   int            `json:"total"`
		Limit   int            `json:"limit"`
		Offset  int            `json:"offset"`
		Lang    string         `json:"lang,omitempty"`
	}

	SearchAPIResponse struct {
		// The search query that was used
		Query string `json:"query" example:"without love"`
		// Matching quotes with relevance scores
		Results []SearchResult `json:"results"`
		// Total number of matches
		Total int `json:"total" example:"25"`
		// Maximum results per page
		Limit int `json:"limit" example:"30"`
		// Pagination offset
		Offset int `json:"offset" example:"0"`
	}
)
