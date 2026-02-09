package quote

import "umineko_quote/internal/dto"

func NewSearchResult(quote dto.ParsedQuote, score int) dto.SearchResult {
	return dto.SearchResult{
		Quote: quote,
		Score: score,
	}
}
