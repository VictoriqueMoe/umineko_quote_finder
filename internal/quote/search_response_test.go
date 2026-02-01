package quote

import "testing"

func TestNewSearchResponse_NilResults(t *testing.T) {
	resp := NewSearchResponse(nil, 10, 0)

	if len(resp.Results) != 0 {
		t.Errorf("Results length: got %d, want 0", len(resp.Results))
	}
	if resp.Total != 0 {
		t.Errorf("Total: got %d, want 0", resp.Total)
	}
	if resp.Limit != 10 {
		t.Errorf("Limit: got %d, want 10", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Errorf("Offset: got %d, want 0", resp.Offset)
	}
}

func TestNewSearchResponse_Pagination(t *testing.T) {
	results := make([]SearchResult, 25)
	for i := 0; i < 25; i++ {
		results[i] = NewSearchResult(ParsedQuote{Text: "quote"}, 100)
	}

	resp := NewSearchResponse(results, 10, 0)

	if len(resp.Results) != 10 {
		t.Errorf("Results length: got %d, want 10", len(resp.Results))
	}
	if resp.Total != 25 {
		t.Errorf("Total: got %d, want 25", resp.Total)
	}
	if resp.Offset != 0 {
		t.Errorf("Offset: got %d, want 0", resp.Offset)
	}
}

func TestNewSearchResponse_PaginationSecondPage(t *testing.T) {
	results := make([]SearchResult, 25)
	for i := 0; i < 25; i++ {
		results[i] = NewSearchResult(ParsedQuote{Text: "quote"}, 100)
	}

	resp := NewSearchResponse(results, 10, 10)

	if len(resp.Results) != 10 {
		t.Errorf("Results length: got %d, want 10", len(resp.Results))
	}
	if resp.Total != 25 {
		t.Errorf("Total: got %d, want 25", resp.Total)
	}
	if resp.Offset != 10 {
		t.Errorf("Offset: got %d, want 10", resp.Offset)
	}
}

func TestNewSearchResponse_OffsetBeyondTotal(t *testing.T) {
	results := make([]SearchResult, 5)
	for i := 0; i < 5; i++ {
		results[i] = NewSearchResult(ParsedQuote{Text: "quote"}, 100)
	}

	resp := NewSearchResponse(results, 10, 100)

	if len(resp.Results) != 0 {
		t.Errorf("Results length: got %d, want 0", len(resp.Results))
	}
	if resp.Total != 5 {
		t.Errorf("Total: got %d, want 5", resp.Total)
	}
	if resp.Offset != 100 {
		t.Errorf("Offset: got %d, want 100", resp.Offset)
	}
}

func TestNewSearchResponse_PartialLastPage(t *testing.T) {
	results := make([]SearchResult, 25)
	for i := 0; i < 25; i++ {
		results[i] = NewSearchResult(ParsedQuote{Text: "quote"}, 100)
	}

	resp := NewSearchResponse(results, 10, 20)

	if len(resp.Results) != 5 {
		t.Errorf("Results length: got %d, want 5", len(resp.Results))
	}
	if resp.Total != 25 {
		t.Errorf("Total: got %d, want 25", resp.Total)
	}
}

func TestNewSearchResponse_LimitLargerThanTotal(t *testing.T) {
	results := make([]SearchResult, 3)
	for i := 0; i < 3; i++ {
		results[i] = NewSearchResult(ParsedQuote{Text: "quote"}, 100)
	}

	resp := NewSearchResponse(results, 100, 0)

	if len(resp.Results) != 3 {
		t.Errorf("Results length: got %d, want 3", len(resp.Results))
	}
	if resp.Total != 3 {
		t.Errorf("Total: got %d, want 3", resp.Total)
	}
}
