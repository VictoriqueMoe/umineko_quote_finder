package quote

import (
	"testing"

	"umineko_quote/internal/dto"
)

func TestConcurrentExactSearch_EmptyIndices(t *testing.T) {
	results := concurrentExactSearch(
		[]int{},
		[]string{},
		[]dto.ParsedQuote{},
		"test",
		func(q dto.ParsedQuote) bool { return true },
	)

	if results != nil {
		t.Errorf("empty indices: got %v, want nil", results)
	}
}

func TestConcurrentExactSearch_FindsMatches(t *testing.T) {
	quotes := []dto.ParsedQuote{
		{ScriptParsedQuote: dto.ScriptParsedQuote{Text: "Hello World"}},
		{ScriptParsedQuote: dto.ScriptParsedQuote{Text: "Goodbye World"}},
		{ScriptParsedQuote: dto.ScriptParsedQuote{Text: "Hello Again"}},
		{ScriptParsedQuote: dto.ScriptParsedQuote{Text: "Something Else"}},
	}
	lowerTexts := []string{
		"hello world",
		"goodbye world",
		"hello again",
		"something else",
	}
	indices := []int{0, 1, 2, 3}

	results := concurrentExactSearch(
		indices,
		lowerTexts,
		quotes,
		"hello",
		func(q dto.ParsedQuote) bool { return true },
	)

	if len(results) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(results))
	}
	for i := 0; i < len(results); i++ {
		if results[i].Score != 100 {
			t.Errorf("result %d score: got %d, want 100", i, results[i].Score)
		}
	}
}

func TestConcurrentExactSearch_RespectsFilter(t *testing.T) {
	quotes := []dto.ParsedQuote{
		{ScriptParsedQuote: dto.ScriptParsedQuote{Text: "Hello World", CharacterID: "10"}},
		{ScriptParsedQuote: dto.ScriptParsedQuote{Text: "Hello Again", CharacterID: "27"}},
		{ScriptParsedQuote: dto.ScriptParsedQuote{Text: "Hello There", CharacterID: "10"}},
	}
	lowerTexts := []string{
		"hello world",
		"hello again",
		"hello there",
	}
	indices := []int{0, 1, 2}

	results := concurrentExactSearch(
		indices,
		lowerTexts,
		quotes,
		"hello",
		func(q dto.ParsedQuote) bool { return q.CharacterID == "10" },
	)

	if len(results) != 2 {
		t.Fatalf("expected 2 filtered matches, got %d", len(results))
	}
	for i := 0; i < len(results); i++ {
		if results[i].Quote.CharacterID != "10" {
			t.Errorf("result %d CharacterID: got %q, want %q", i, results[i].Quote.CharacterID, "10")
		}
	}
}

func TestConcurrentExactSearch_NoMatches(t *testing.T) {
	quotes := []dto.ParsedQuote{
		{ScriptParsedQuote: dto.ScriptParsedQuote{Text: "Hello World"}},
		{ScriptParsedQuote: dto.ScriptParsedQuote{Text: "Goodbye World"}},
	}
	lowerTexts := []string{
		"hello world",
		"goodbye world",
	}
	indices := []int{0, 1}

	results := concurrentExactSearch(
		indices,
		lowerTexts,
		quotes,
		"beatrice",
		func(q dto.ParsedQuote) bool { return true },
	)

	if len(results) != 0 {
		t.Errorf("expected 0 matches, got %d", len(results))
	}
}

func TestConcurrentExactSearch_SubsetIndices(t *testing.T) {
	quotes := []dto.ParsedQuote{
		{ScriptParsedQuote: dto.ScriptParsedQuote{Text: "Hello World"}},
		{ScriptParsedQuote: dto.ScriptParsedQuote{Text: "Hello Again"}},
		{ScriptParsedQuote: dto.ScriptParsedQuote{Text: "Hello There"}},
	}
	lowerTexts := []string{
		"hello world",
		"hello again",
		"hello there",
	}
	indices := []int{0, 2}

	results := concurrentExactSearch(
		indices,
		lowerTexts,
		quotes,
		"hello",
		func(q dto.ParsedQuote) bool { return true },
	)

	if len(results) != 2 {
		t.Fatalf("expected 2 matches from subset, got %d", len(results))
	}
}

func TestConcurrentExactSearch_CaseInsensitive(t *testing.T) {
	quotes := []dto.ParsedQuote{
		{ScriptParsedQuote: dto.ScriptParsedQuote{Text: "Hello WORLD"}},
	}
	lowerTexts := []string{
		"hello world",
	}
	indices := []int{0}

	results := concurrentExactSearch(
		indices,
		lowerTexts,
		quotes,
		"world",
		func(q dto.ParsedQuote) bool { return true },
	)

	if len(results) != 1 {
		t.Fatalf("expected 1 case-insensitive match, got %d", len(results))
	}
}
