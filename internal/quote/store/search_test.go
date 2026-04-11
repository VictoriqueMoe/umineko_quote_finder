package store

import (
	"testing"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"

	"umineko_quote/internal/dto"
)

func TestConcurrentExactSearch_EmptyIndices(t *testing.T) {
	results := concurrentExactSearchGeneric(
		[]int{},
		[]string{},
		[]dto.UminekoQuote{},
		"test",
		false,
		func(idx int) bool { return true },
	)

	if results != nil {
		t.Errorf("empty indices: got %v, want nil", results)
	}
}

func TestConcurrentExactSearch_FindsMatches(t *testing.T) {
	quotes := []dto.UminekoQuote{
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Hello World"}}},
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Goodbye World"}}},
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Hello Again"}}},
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Something Else"}}},
	}
	lowerTexts := []string{
		"hello world",
		"goodbye world",
		"hello again",
		"something else",
	}
	indices := []int{0, 1, 2, 3}

	results := concurrentExactSearchGeneric(
		indices,
		lowerTexts,
		quotes,
		"hello",
		false,
		func(idx int) bool { return true },
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
	quotes := []dto.UminekoQuote{
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Hello World", CharacterID: "10"}}},
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Hello Again", CharacterID: "27"}}},
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Hello There", CharacterID: "10"}}},
	}
	lowerTexts := []string{
		"hello world",
		"hello again",
		"hello there",
	}
	indices := []int{0, 1, 2}

	results := concurrentExactSearchGeneric(
		indices,
		lowerTexts,
		quotes,
		"hello",
		false,
		func(idx int) bool { return quotes[idx].CharacterID == "10" },
	)

	if len(results) != 2 {
		t.Fatalf("expected 2 filtered matches, got %d", len(results))
	}
	for i := 0; i < len(results); i++ {
		q := results[i].Quote.(dto.UminekoQuote)
		if q.CharacterID != "10" {
			t.Errorf("result %d CharacterID: got %q, want %q", i, q.CharacterID, "10")
		}
	}
}

func TestConcurrentExactSearch_NoMatches(t *testing.T) {
	quotes := []dto.UminekoQuote{
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Hello World"}}},
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Goodbye World"}}},
	}
	lowerTexts := []string{
		"hello world",
		"goodbye world",
	}
	indices := []int{0, 1}

	results := concurrentExactSearchGeneric(
		indices,
		lowerTexts,
		quotes,
		"beatrice",
		false,
		func(idx int) bool { return true },
	)

	if len(results) != 0 {
		t.Errorf("expected 0 matches, got %d", len(results))
	}
}

func TestConcurrentExactSearch_SubsetIndices(t *testing.T) {
	quotes := []dto.UminekoQuote{
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Hello World"}}},
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Hello Again"}}},
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Hello There"}}},
	}
	lowerTexts := []string{
		"hello world",
		"hello again",
		"hello there",
	}
	indices := []int{0, 2}

	results := concurrentExactSearchGeneric(
		indices,
		lowerTexts,
		quotes,
		"hello",
		false,
		func(idx int) bool { return true },
	)

	if len(results) != 2 {
		t.Fatalf("expected 2 matches from subset, got %d", len(results))
	}
}

func TestConcurrentExactSearch_CaseInsensitive(t *testing.T) {
	quotes := []dto.UminekoQuote{
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Hello WORLD"}}},
	}
	lowerTexts := []string{
		"hello world",
	}
	indices := []int{0}

	results := concurrentExactSearchGeneric(
		indices,
		lowerTexts,
		quotes,
		"world",
		false,
		func(idx int) bool { return true },
	)

	if len(results) != 1 {
		t.Fatalf("expected 1 case-insensitive match, got %d", len(results))
	}
}

func TestContainsWord(t *testing.T) {
	cases := []struct {
		text string
		word string
		want bool
	}{
		{"i have a brother", "broth", false},
		{"hot broth on the table", "broth", true},
		{"broth is good", "broth", true},
		{"all broth", "broth", true},
		{"brothers in arms", "broth", false},
		{"un-broth-like", "broth", true},
		{"café au lait", "café", true},
		{"cafés", "café", false},
		{"真実は一つ", "真実", false},
	}
	for _, c := range cases {
		got := containsWord(c.text, c.word)
		if got != c.want {
			t.Errorf("containsWord(%q, %q) = %v, want %v", c.text, c.word, got, c.want)
		}
	}
}

func TestConcurrentExactSearch_ExactMode(t *testing.T) {
	quotes := []dto.UminekoQuote{
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "I have a brother"}}},
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "hot broth on the table"}}},
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Brothers in arms"}}},
	}
	lowerTexts := []string{
		"i have a brother",
		"hot broth on the table",
		"brothers in arms",
	}
	indices := []int{0, 1, 2}

	loose := concurrentExactSearchGeneric(
		indices,
		lowerTexts,
		quotes,
		"broth",
		false,
		func(idx int) bool { return true },
	)
	if len(loose) != 3 {
		t.Fatalf("loose mode: expected 3 matches, got %d", len(loose))
	}

	exact := concurrentExactSearchGeneric(
		indices,
		lowerTexts,
		quotes,
		"broth",
		true,
		func(idx int) bool { return true },
	)
	if len(exact) != 1 {
		t.Fatalf("exact mode: expected 1 match, got %d", len(exact))
	}
}
