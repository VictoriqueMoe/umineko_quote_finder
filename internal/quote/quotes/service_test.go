package quotes

import (
	"fmt"
	"testing"

	"umineko_quote/internal/dao/memory"
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/game"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
)

var testNames = map[string]string{"10": "Ushiromiya Battler", "27": "Beatrice"}

func resolveTestName(id string) string {
	return testNames[id]
}

func testQuote(idx int, text string, characterID string, audioID string) dto.HigurashiQuote {
	return dto.HigurashiQuote{
		HigurashiQuote: scriptdto.HigurashiQuote{
			BaseQuote: scriptdto.BaseQuote{
				Text:        text,
				CharacterID: characterID,
				AudioID:     audioID,
				Episode:     1,
			},
		},
		Index: idx,
	}
}

func newTestService(quotes map[language.Language][]dto.HigurashiQuote) Service[dto.HigurashiQuote] {
	return NewService[dto.HigurashiQuote](memory.NewQuotes(game.Higurashi{}, quotes), resolveTestName)
}

func repeatedQuotes(n int, characterID string) map[language.Language][]dto.HigurashiQuote {
	en := make([]dto.HigurashiQuote, n)
	for i := range en {
		en[i] = testQuote(i, "quote", characterID, "")
	}
	return map[language.Language][]dto.HigurashiQuote{language.English: en}
}

func contextQuotes() map[language.Language][]dto.HigurashiQuote {
	en := make([]dto.HigurashiQuote, 10)
	for i := range en {
		en[i] = testQuote(i, fmt.Sprintf("line %d", i), "10", fmt.Sprintf("a%d", i))
	}
	return map[language.Language][]dto.HigurashiQuote{language.English: en}
}

func search(t *testing.T, svc Service[dto.HigurashiQuote], p params.SearchParams) dto.SearchResponse {
	t.Helper()

	resp, err := svc.Search(p, params.QuoteScope{})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	return resp
}

func browse(t *testing.T, svc Service[dto.HigurashiQuote], p params.BrowseParams) dto.CharacterResponse {
	t.Helper()

	resp, err := svc.Browse(p, params.QuoteScope{})
	if err != nil {
		t.Fatalf("Browse: %v", err)
	}
	return resp
}

func TestSearch_NoMatches(t *testing.T) {
	resp := search(t, newTestService(repeatedQuotes(5, "10")), params.SearchParams{Query: "xyzzy", Lang: language.English})

	if resp.Results == nil || len(resp.Results) != 0 {
		t.Errorf("Results: got %v, want empty non-nil", resp.Results)
	}
	if resp.Total != 0 {
		t.Errorf("Total: got %d, want 0", resp.Total)
	}
	if resp.Limit != defaultSearchLimit {
		t.Errorf("Limit: got %d, want %d", resp.Limit, defaultSearchLimit)
	}
}

func TestSearch_Pagination(t *testing.T) {
	cases := []struct {
		name        string
		quotes      int
		limit       int
		offset      int
		wantResults int
	}{
		{"first page", 25, 10, 0, 10},
		{"second page", 25, 10, 10, 10},
		{"partial last page", 25, 10, 20, 5},
		{"offset beyond total", 5, 10, 100, 0},
		{"limit larger than total", 3, 100, 0, 3},
	}

	for _, c := range cases {
		resp := search(t, newTestService(repeatedQuotes(c.quotes, "10")), params.SearchParams{Query: "quote", Lang: language.English, Limit: c.limit, Offset: c.offset})

		if len(resp.Results) != c.wantResults {
			t.Errorf("%s: Results length: got %d, want %d", c.name, len(resp.Results), c.wantResults)
		}
		if resp.Total != c.quotes {
			t.Errorf("%s: Total: got %d, want %d", c.name, resp.Total, c.quotes)
		}
		if resp.Limit != c.limit || resp.Offset != c.offset {
			t.Errorf("%s: Limit/Offset: got %d/%d, want %d/%d", c.name, resp.Limit, resp.Offset, c.limit, c.offset)
		}
	}
}

func TestSearch_AutoPicksFirstLanguageWithResults(t *testing.T) {
	quotes := map[language.Language][]dto.HigurashiQuote{
		language.English:  {testQuote(0, "hello", "10", "")},
		language.Japanese: {testQuote(0, "ひぐらし", "10", "")},
	}

	resp := search(t, newTestService(quotes), params.SearchParams{Query: "ひぐらし", Lang: language.Auto})

	if resp.Lang != string(language.Japanese) || resp.Total != 1 {
		t.Errorf("auto search: got lang %q total %d, want ja 1", resp.Lang, resp.Total)
	}
}

func TestSearch_AutoFallsBackToEnglish(t *testing.T) {
	resp := search(t, newTestService(repeatedQuotes(3, "10")), params.SearchParams{Query: "xyzzy", Lang: language.Auto})

	if resp.Total != 0 || resp.Lang != "" {
		t.Errorf("auto fallback: got lang %q total %d, want no lang and 0", resp.Lang, resp.Total)
	}
	if resp.Limit != defaultSearchLimit {
		t.Errorf("auto fallback Limit: got %d, want %d", resp.Limit, defaultSearchLimit)
	}
}

func TestBrowse_ResolvesCharacterName(t *testing.T) {
	resp := browse(t, newTestService(repeatedQuotes(25, "27")), params.BrowseParams{Lang: language.English, CharacterID: "27", Limit: 10})

	if resp.Character != "Beatrice" {
		t.Errorf("Character: got %q, want Beatrice", resp.Character)
	}
	if got := len(resp.Quotes.([]dto.HigurashiQuote)); got != 10 {
		t.Errorf("Quotes length: got %d, want 10", got)
	}
	if resp.Total != 25 {
		t.Errorf("Total: got %d, want 25", resp.Total)
	}
}

func TestBrowse_EmptyCharacterID(t *testing.T) {
	resp := browse(t, newTestService(repeatedQuotes(5, "10")), params.BrowseParams{Lang: language.English})

	if resp.CharacterID != "" || resp.Character != "" {
		t.Errorf("CharacterID/Character: got %q/%q, want empty", resp.CharacterID, resp.Character)
	}
	if resp.Limit != defaultBrowseLimit {
		t.Errorf("Limit: got %d, want %d", resp.Limit, defaultBrowseLimit)
	}
}

func TestBrowse_UnknownCharacterReturnsNothing(t *testing.T) {
	resp := browse(t, newTestService(repeatedQuotes(5, "10")), params.BrowseParams{Lang: language.English, CharacterID: "99"})

	quotes := resp.Quotes.([]dto.HigurashiQuote)
	if quotes == nil || len(quotes) != 0 || resp.Total != 0 {
		t.Errorf("unknown character: got %d quotes, total %d, want empty non-nil and 0", len(quotes), resp.Total)
	}
}

func TestBrowse_OffsetBeyondTotal(t *testing.T) {
	resp := browse(t, newTestService(repeatedQuotes(5, "10")), params.BrowseParams{Lang: language.English, CharacterID: "10", Limit: 10, Offset: 100})

	if got := len(resp.Quotes.([]dto.HigurashiQuote)); got != 0 {
		t.Errorf("Quotes length: got %d, want 0", got)
	}
	if resp.Total != 5 {
		t.Errorf("Total: got %d, want 5", resp.Total)
	}
}

func TestRandom_UnknownCharacterReturnsNil(t *testing.T) {
	q, err := newTestService(repeatedQuotes(5, "10")).Random(language.English, "99", 0, params.QuoteScope{})
	if err != nil {
		t.Fatalf("Random: %v", err)
	}
	if q != nil {
		t.Errorf("Random with unknown character: got %+v, want nil", q)
	}
}

func TestRandom_NoFiltersExcludesNarrator(t *testing.T) {
	quotes := map[language.Language][]dto.HigurashiQuote{
		language.English: {testQuote(0, "a", game.NarratorID, ""), testQuote(1, "b", "10", "")},
	}
	svc := newTestService(quotes)

	for range 20 {
		q, err := svc.Random(language.English, "", 0, params.QuoteScope{})
		if err != nil {
			t.Fatalf("Random: %v", err)
		}
		if q == nil || q.CharacterID == game.NarratorID {
			t.Fatalf("Random with no filters: got %+v, want a non-narrator line", q)
		}
	}
}

func TestGetContext(t *testing.T) {
	ctx, err := newTestService(contextQuotes()).GetContext(language.English, "a5", 2)
	if err != nil {
		t.Fatalf("GetContext: %v", err)
	}

	if ctx.Quote.(dto.HigurashiQuote).Index != 5 {
		t.Errorf("Quote: got %+v, want index 5", ctx.Quote)
	}
	if got := ctx.Before.([]dto.HigurashiQuote); len(got) != 2 || got[0].Index != 3 || got[1].Index != 4 {
		t.Errorf("Before: got %+v, want indexes 3 4", got)
	}
	if got := ctx.After.([]dto.HigurashiQuote); len(got) != 2 || got[0].Index != 6 || got[1].Index != 7 {
		t.Errorf("After: got %+v, want indexes 6 7", got)
	}
}

func TestGetContext_ClipsAtEdges(t *testing.T) {
	svc := newTestService(contextQuotes())

	start, err := svc.GetContext(language.English, "a0", 3)
	if err != nil {
		t.Fatalf("GetContext: %v", err)
	}
	if len(start.Before.([]dto.HigurashiQuote)) != 0 || len(start.After.([]dto.HigurashiQuote)) != 3 {
		t.Errorf("start: got %d before, %d after, want 0 and 3", len(start.Before.([]dto.HigurashiQuote)), len(start.After.([]dto.HigurashiQuote)))
	}

	end, err := svc.GetContext(language.English, "a9", 3)
	if err != nil {
		t.Fatalf("GetContext: %v", err)
	}
	if len(end.Before.([]dto.HigurashiQuote)) != 3 || len(end.After.([]dto.HigurashiQuote)) != 0 {
		t.Errorf("end: got %d before, %d after, want 3 and 0", len(end.Before.([]dto.HigurashiQuote)), len(end.After.([]dto.HigurashiQuote)))
	}
}

func TestGetContext_LinesClamped(t *testing.T) {
	en := make([]dto.HigurashiQuote, 60)
	for i := range en {
		en[i] = testQuote(i, "line", "10", fmt.Sprintf("a%d", i))
	}

	ctx, err := newTestService(map[language.Language][]dto.HigurashiQuote{language.English: en}).GetContext(language.English, "a30", 100)
	if err != nil {
		t.Fatalf("GetContext: %v", err)
	}
	if got := len(ctx.Before.([]dto.HigurashiQuote)); got != maxContextLines {
		t.Errorf("Before: got %d lines, want %d", got, maxContextLines)
	}
}

func TestGetContext_UnknownAudioID(t *testing.T) {
	ctx, err := newTestService(contextQuotes()).GetContext(language.English, "nope", 2)
	if err != nil {
		t.Fatalf("GetContext: %v", err)
	}
	if ctx != nil {
		t.Errorf("GetContext with unknown audio id: got %+v, want nil", ctx)
	}
}

func TestNearestVoicedAudioID_TakesFirstOfCompositeID(t *testing.T) {
	quotes := map[language.Language][]dto.HigurashiQuote{
		language.English: {
			testQuote(0, "a", "10", "a0"),
			testQuote(1, "b", "10", ""),
			testQuote(2, "c", "27", "b1, b2"),
		},
	}

	next, err := newTestService(quotes).NearestVoicedAudioID(language.English, "a0", "next")
	if err != nil {
		t.Fatalf("NearestVoicedAudioID: %v", err)
	}
	if next != "b1" {
		t.Errorf("next: got %q, want b1", next)
	}
}
