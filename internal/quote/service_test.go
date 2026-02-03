package quote

import (
	"strings"
	"testing"
)

var testService = NewService()

func TestService_Search_ExactMatch(t *testing.T) {
	svc := testService

	resp := svc.Search("Beatrice", "en", 10, 0, "", 0, TruthAll)

	if resp.Total == 0 {
		t.Fatal("expected search results for 'Beatrice'")
	}
	if len(resp.Results) == 0 {
		t.Fatal("expected non-empty results slice")
	}
	if resp.Limit != 10 {
		t.Errorf("Limit: got %d, want 10", resp.Limit)
	}
}

func TestService_Search_DefaultValues(t *testing.T) {
	svc := testService

	resp := svc.Search("witch", "", 0, -1, "", 0, TruthAll)

	if resp.Limit != 30 {
		t.Errorf("default limit: got %d, want 30", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Errorf("default offset: got %d, want 0", resp.Offset)
	}
}

func TestService_Search_WithCharacterFilter(t *testing.T) {
	svc := testService

	resp := svc.Search("witch", "en", 10, 0, "10", 0, TruthAll)

	for i := 0; i < len(resp.Results); i++ {
		if resp.Results[i].Quote.CharacterID != "10" {
			t.Errorf("result %d CharacterID: got %q, want %q", i, resp.Results[i].Quote.CharacterID, "10")
		}
	}
}

func TestService_Search_WithEpisodeFilter(t *testing.T) {
	svc := testService

	resp := svc.Search("witch", "en", 10, 0, "", 1, TruthAll)

	for i := 0; i < len(resp.Results); i++ {
		if resp.Results[i].Quote.Episode != 1 {
			t.Errorf("result %d Episode: got %d, want 1", i, resp.Results[i].Quote.Episode)
		}
	}
}

func TestService_Search_RedTruthFilter(t *testing.T) {
	svc := testService

	resp := svc.Search("truth", "en", 10, 0, "", 0, TruthRed)

	for i := 0; i < len(resp.Results); i++ {
		if !strings.Contains(resp.Results[i].Quote.TextHtml, "red-truth") {
			t.Errorf("result %d should contain red-truth in HTML", i)
		}
	}
}

func TestService_Search_NoResults(t *testing.T) {
	svc := testService

	resp := svc.Search("xyzzyxyzzyxyzzy", "en", 10, 0, "", 0, TruthAll)

	if resp.Total != 0 {
		t.Errorf("Total: got %d, want 0", resp.Total)
	}
	if len(resp.Results) != 0 {
		t.Errorf("Results: got %d, want 0", len(resp.Results))
	}
}

func TestService_Search_Japanese(t *testing.T) {
	svc := testService

	resp := svc.Search("ベアトリーチェ", "ja", 10, 0, "", 0, TruthAll)

	if resp.Total == 0 {
		t.Fatal("expected Japanese search results")
	}
}

func TestService_Search_UnknownLang(t *testing.T) {
	svc := testService

	resp := svc.Search("test", "fr", 10, 0, "", 0, TruthAll)

	if resp.Total != 0 {
		t.Errorf("Total for unknown lang: got %d, want 0", resp.Total)
	}
}

func TestService_Browse(t *testing.T) {
	svc := testService

	resp := svc.Browse("en", "10", 10, 0, 0, TruthAll)

	if resp.Total == 0 {
		t.Fatal("expected browse results for Battler")
	}
	if resp.CharacterID != "10" {
		t.Errorf("CharacterID: got %q, want %q", resp.CharacterID, "10")
	}
	if resp.Character != CharacterNames["10"] {
		t.Errorf("Character: got %q, want %q", resp.Character, CharacterNames["10"])
	}
	if len(resp.Quotes) > 10 {
		t.Errorf("Quotes length exceeds limit: got %d", len(resp.Quotes))
	}
}

func TestService_Browse_WithEpisode(t *testing.T) {
	svc := testService

	resp := svc.Browse("en", "10", 10, 0, 1, TruthAll)

	for i := 0; i < len(resp.Quotes); i++ {
		if resp.Quotes[i].Episode != 1 {
			t.Errorf("quote %d Episode: got %d, want 1", i, resp.Quotes[i].Episode)
		}
	}
}

func TestService_Browse_DefaultValues(t *testing.T) {
	svc := testService

	resp := svc.Browse("", "", 0, -1, 0, TruthAll)

	if resp.Limit != 50 {
		t.Errorf("default limit: got %d, want 50", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Errorf("default offset: got %d, want 0", resp.Offset)
	}
}

func TestService_Browse_UnknownLang(t *testing.T) {
	svc := testService

	resp := svc.Browse("fr", "10", 10, 0, 0, TruthAll)

	if resp.Total != 0 {
		t.Errorf("Total for unknown lang: got %d, want 0", resp.Total)
	}
}

func TestService_GetByCharacter(t *testing.T) {
	svc := testService

	resp := svc.GetByCharacter("en", "27", 10, 0, 0, TruthAll)

	if resp.Total == 0 {
		t.Fatal("expected results for Beatrice")
	}
	if resp.Character != CharacterNames["27"] {
		t.Errorf("Character: got %q, want %q", resp.Character, CharacterNames["27"])
	}
	for i := 0; i < len(resp.Quotes); i++ {
		if resp.Quotes[i].CharacterID != "27" {
			t.Errorf("quote %d CharacterID: got %q, want %q", i, resp.Quotes[i].CharacterID, "27")
		}
	}
}

func TestService_GetByCharacter_WithEpisode(t *testing.T) {
	svc := testService

	resp := svc.GetByCharacter("en", "10", 10, 0, 1, TruthAll)

	for i := 0; i < len(resp.Quotes); i++ {
		if resp.Quotes[i].Episode != 1 {
			t.Errorf("quote %d Episode: got %d, want 1", i, resp.Quotes[i].Episode)
		}
		if resp.Quotes[i].CharacterID != "10" {
			t.Errorf("quote %d CharacterID: got %q, want %q", i, resp.Quotes[i].CharacterID, "10")
		}
	}
}

func TestService_GetByCharacter_UnknownCharacter(t *testing.T) {
	svc := testService

	resp := svc.GetByCharacter("en", "999", 10, 0, 0, TruthAll)

	if resp.Total != 0 {
		t.Errorf("Total for unknown character: got %d, want 0", resp.Total)
	}
}

func TestService_GetByCharacter_DefaultValues(t *testing.T) {
	svc := testService

	resp := svc.GetByCharacter("", "10", 0, -1, 0, TruthAll)

	if resp.Limit != 50 {
		t.Errorf("default limit: got %d, want 50", resp.Limit)
	}
}

func TestService_GetByAudioID(t *testing.T) {
	svc := testService

	q := svc.GetByAudioID("en", "11900001")

	if q == nil {
		t.Fatal("expected to find quote by audio ID")
	}
	if q.CharacterID != "19" {
		t.Errorf("CharacterID: got %q, want %q", q.CharacterID, "19")
	}
}

func TestService_GetByAudioID_NotFound(t *testing.T) {
	svc := testService

	q := svc.GetByAudioID("en", "99999999")

	if q != nil {
		t.Errorf("expected nil for unknown audio ID, got %+v", q)
	}
}

func TestService_GetByAudioID_DefaultLang(t *testing.T) {
	svc := testService

	q := svc.GetByAudioID("", "11900001")

	if q == nil {
		t.Fatal("expected to find quote with empty lang (should default to en)")
	}
}

func TestService_GetByAudioID_UnknownLang(t *testing.T) {
	svc := testService

	q := svc.GetByAudioID("fr", "11900001")

	if q != nil {
		t.Errorf("expected nil for unknown lang, got %+v", q)
	}
}

func TestService_Random(t *testing.T) {
	svc := testService

	q := svc.Random("en", "", 0, TruthAll)

	if q == nil {
		t.Fatal("expected a random quote")
	}
	if q.CharacterID == "narrator" {
		t.Error("Random with no filters should exclude narrator")
	}
}

func TestService_Random_WithCharacter(t *testing.T) {
	svc := testService

	for i := 0; i < 10; i++ {
		q := svc.Random("en", "27", 0, TruthAll)
		if q == nil {
			t.Fatal("expected a random Beatrice quote")
		}
		if q.CharacterID != "27" {
			t.Errorf("CharacterID: got %q, want %q", q.CharacterID, "27")
		}
	}
}

func TestService_Random_WithEpisode(t *testing.T) {
	svc := testService

	for i := 0; i < 10; i++ {
		q := svc.Random("en", "", 1, TruthAll)
		if q == nil {
			t.Fatal("expected a random episode 1 quote")
		}
		if q.Episode != 1 {
			t.Errorf("Episode: got %d, want 1", q.Episode)
		}
	}
}

func TestService_Random_WithCharacterAndEpisode(t *testing.T) {
	svc := testService

	for i := 0; i < 10; i++ {
		q := svc.Random("en", "10", 1, TruthAll)
		if q == nil {
			t.Fatal("expected a random Battler ep1 quote")
		}
		if q.CharacterID != "10" {
			t.Errorf("CharacterID: got %q, want %q", q.CharacterID, "10")
		}
		if q.Episode != 1 {
			t.Errorf("Episode: got %d, want 1", q.Episode)
		}
	}
}

func TestService_Random_RedTruth(t *testing.T) {
	svc := testService

	for i := 0; i < 10; i++ {
		q := svc.Random("en", "", 0, TruthRed)
		if q == nil {
			t.Fatal("expected a random red truth quote")
		}
		if !strings.Contains(q.TextHtml, "red-truth") {
			t.Errorf("expected red-truth in HTML: %q", q.TextHtml)
		}
	}
}

func TestService_Random_DefaultLang(t *testing.T) {
	svc := testService

	q := svc.Random("", "", 0, TruthAll)

	if q == nil {
		t.Fatal("expected a random quote with default lang")
	}
}

func TestService_Random_UnknownLang(t *testing.T) {
	svc := testService

	q := svc.Random("fr", "", 0, TruthAll)

	if q != nil {
		t.Errorf("expected nil for unknown lang, got %+v", q)
	}
}

func TestService_GetContext(t *testing.T) {
	svc := testService

	// Use an audio ID that is not at the very start of the quotes slice
	resp := svc.Search("Beatrice", "en", 10, 0, "", 0, TruthAll)
	if resp.Total == 0 {
		t.Fatal("need search results to find a mid-slice audio ID")
	}
	var midAudioId string
	for _, r := range resp.Results {
		if r.Quote.AudioID != "" {
			midAudioId = r.Quote.AudioID
			break
		}
	}
	if midAudioId == "" {
		t.Skip("no quote with audioId found in search results")
	}
	// Handle composite audio IDs
	parts := strings.SplitN(midAudioId, ", ", 2)
	midAudioId = parts[0]

	result := svc.GetContext("en", midAudioId, 5)

	if result == nil {
		t.Fatal("expected context result")
	}
	if result.Quote.AudioID == "" {
		t.Error("expected quote to have an audio ID")
	}
	if len(result.Before) > 5 {
		t.Errorf("Before length exceeds lines: got %d", len(result.Before))
	}
	if len(result.After) > 5 {
		t.Errorf("After length exceeds lines: got %d", len(result.After))
	}
	totalContext := len(result.Before) + len(result.After)
	if totalContext == 0 {
		t.Error("expected at least some context lines")
	}
}

func TestService_GetContext_NotFound(t *testing.T) {
	svc := testService

	result := svc.GetContext("en", "99999999", 5)

	if result != nil {
		t.Errorf("expected nil for unknown audio ID, got %+v", result)
	}
}

func TestService_GetContext_EdgeOfSlice(t *testing.T) {
	svc := testService

	// 11900001 is near the start of the slice
	result := svc.GetContext("en", "11900001", 5)

	if result == nil {
		t.Fatal("expected context result for edge-of-slice quote")
	}
	// Before may be empty if the quote is at the start
	if len(result.Before) > 5 {
		t.Errorf("Before length exceeds lines: got %d", len(result.Before))
	}
	if len(result.After) > 5 {
		t.Errorf("After length exceeds lines: got %d", len(result.After))
	}
}

func TestService_GetContext_DefaultLines(t *testing.T) {
	svc := testService

	result := svc.GetContext("en", "11900001", 0)

	if result == nil {
		t.Fatal("expected context result with default lines")
	}
	if len(result.Before) > 5 {
		t.Errorf("Before with default lines: got %d, want <= 5", len(result.Before))
	}
	if len(result.After) > 5 {
		t.Errorf("After with default lines: got %d, want <= 5", len(result.After))
	}
}

func TestService_GetContext_CapsAtMax(t *testing.T) {
	svc := testService

	result := svc.GetContext("en", "11900001", 100)

	if result == nil {
		t.Fatal("expected context result")
	}
	if len(result.Before) > 20 {
		t.Errorf("Before exceeds max: got %d", len(result.Before))
	}
	if len(result.After) > 20 {
		t.Errorf("After exceeds max: got %d", len(result.After))
	}
}

func TestService_GetContext_UnknownLang(t *testing.T) {
	svc := testService

	result := svc.GetContext("fr", "11900001", 5)

	if result != nil {
		t.Errorf("expected nil for unknown lang, got %+v", result)
	}
}

func TestService_GetContext_DefaultLang(t *testing.T) {
	svc := testService

	result := svc.GetContext("", "11900001", 5)

	if result == nil {
		t.Fatal("expected context result with default lang")
	}
}

func TestService_GetCharacters(t *testing.T) {
	svc := testService

	chars := svc.GetCharacters()

	if len(chars) == 0 {
		t.Fatal("expected characters map to be non-empty")
	}
	if chars["10"] != CharacterNames["10"] {
		t.Errorf("chars[10]: got %q, want %q", chars["10"], CharacterNames["10"])
	}
	if chars["27"] != CharacterNames["27"] {
		t.Errorf("chars[27]: got %q, want %q", chars["27"], CharacterNames["27"])
	}
}

func TestService_GetStats(t *testing.T) {
	svc := testService

	stats := svc.GetStats()

	if stats == nil {
		t.Fatal("expected stats to be non-nil")
	}

	result := stats.Compute(AllEpisodes)
	if result == nil {
		t.Fatal("expected Compute(AllEpisodes) to return non-nil")
	}
}

func TestService_Browse_RedTruthFilter(t *testing.T) {
	svc := testService

	resp := svc.Browse("en", "", 10, 0, 0, TruthRed)

	for i := 0; i < len(resp.Quotes); i++ {
		if !strings.Contains(resp.Quotes[i].TextHtml, "red-truth") {
			t.Errorf("quote %d should contain red-truth in HTML", i)
		}
	}
}

func TestService_GetByCharacter_BlueTruthFilter(t *testing.T) {
	svc := testService

	resp := svc.GetByCharacter("en", "10", 100, 0, 0, TruthBlue)

	for i := 0; i < len(resp.Quotes); i++ {
		if !strings.Contains(resp.Quotes[i].TextHtml, "blue-truth") {
			t.Errorf("quote %d should contain blue-truth in HTML", i)
		}
	}
}

func TestService_GetByAudioID_CompositeAudioID(t *testing.T) {
	svc := testService

	q1 := svc.GetByAudioID("en", "11900001")
	if q1 == nil {
		t.Fatal("expected to find quote by first audio ID in composite")
	}

	q2 := svc.GetByAudioID("en", "11900002")
	if q2 == nil {
		t.Fatal("expected to find quote by second audio ID in composite")
	}

	if q1.CharacterID != q2.CharacterID {
		t.Errorf("both audio IDs should resolve to same character: %q vs %q", q1.CharacterID, q2.CharacterID)
	}
}
