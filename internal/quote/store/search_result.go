package store

import "umineko_quote/internal/dto"

func NewSearchResult(quote any, score int) dto.SearchResult {
	return dto.SearchResult{
		Quote: quote,
		Score: score,
	}
}
