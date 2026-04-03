package store

import "umineko_quote/internal/dto"

func NewSearchResponse(results []dto.SearchResult, limit int, offset int) dto.SearchResponse {
	if results == nil {
		results = []dto.SearchResult{}
	}

	total := len(results)

	if offset >= total {
		return dto.SearchResponse{
			Results: []dto.SearchResult{},
			Total:   total,
			Limit:   limit,
			Offset:  offset,
		}
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return dto.SearchResponse{
		Results: results[offset:end],
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}
}
