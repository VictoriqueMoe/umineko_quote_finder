package quote

import (
	"strings"
	"testing"
)

func getTestService() Service {
	return NewService()
}

func TestService_Search_ExactMatch(t *testing.T) {
	svc := getTestService()

	resp := svc.Search("Beatrice", "en", 10, 0, "", 0, false, TruthAll)

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
	svc := getTestService()

	resp := svc.Search("witch", "", 0, -1, "", 0, false, TruthAll)

	if resp.Limit != 30 {
		t.Errorf("default limit: got %d, want 30", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Errorf("default offset: got %d, want 0", resp.Offset)
	}
}

func TestService_Search_WithCharacterFilter(t *testing.T) {
	svc := getTestService()

	resp := svc.Search("witch", "en", 10, 0, "10", 0, false, TruthAll)

	for i := 0; i < len(resp.Results); i++ {
		if resp.Results[i].Quote.CharacterID != "10" {
			t.Errorf("result %d CharacterID: got %q, want %q", i, resp.Results[i].Quote.CharacterID, "10")
		}
	}
}

func TestService_Search_WithEpisodeFilter(t *testing.T) {
	svc := getTestService()

	resp := svc.Search("witch", "en", 10, 0, "", 1, false, TruthAll)

	for i := 0; i < len(resp.Results); i++ {
		if resp.Results[i].Quote.Episode != 1 {
			t.Errorf("result %d Episode: got %d, want 1", i, resp.Results[i].Quote.Episode)
		}
	}
}

func TestService_Search_RedTruthFilter(t *testing.T) {
	svc := getTestService()

	resp := svc.Search("truth", "en", 10, 0, "", 0, false, TruthRed)

	for i := 0; i < len(resp.Results); i++ {
		if !strings.Contains(resp.Results[i].Quote.TextHtml, "red-truth") {
			t.Errorf("result %d should contain red-truth in HTML", i)
		}
	}
}

func TestService_Search_NoResults(t *testing.T) {
	svc := getTestService()

	resp := svc.Search("xyzzyxyzzyxyzzy", "en", 10, 0, "", 0, false, TruthAll)

	if resp.Total != 0 {
		t.Errorf("Total: got %d, want 0", resp.Total)
	}
	if len(resp.Results) != 0 {
		t.Errorf("Results: got %d, want 0", len(resp.Results))
	}
}

func TestService_Search_ForceFuzzy(t *testing.T) {
	svc := getTestService()

	resp := svc.Search("Beatrice", "en", 10, 0, "", 0, true, TruthAll)

	if resp.Total == 0 {
		t.Fatal("expected fuzzy search results for 'Beatrice'")
	}
}

func TestService_Search_Japanese(t *testing.T) {
	svc := getTestService()

	resp := svc.Search("ベアトリーチェ", "ja", 10, 0, "", 0, false, TruthAll)

	if resp.Total == 0 {
		t.Fatal("expected Japanese search results")
	}
}

func TestService_Search_UnknownLang(t *testing.T) {
	svc := getTestService()

	resp := svc.Search("test", "fr", 10, 0, "", 0, false, TruthAll)

	if resp.Total != 0 {
		t.Errorf("Total for unknown lang: got %d, want 0", resp.Total)
	}
}

func TestService_Browse(t *testing.T) {
	svc := getTestService()

	resp := svc.Browse("en", "10", 10, 0, 0, TruthAll)

	if resp.Total == 0 {
		t.Fatal("expected browse results for Battler")
	}
	if resp.CharacterID != "10" {
		t.Errorf("CharacterID: got %q, want %q", resp.CharacterID, "10")
	}
	if resp.Character != "Battler" {
		t.Errorf("Character: got %q, want %q", resp.Character, "Battler")
	}
	if len(resp.Quotes) > 10 {
		t.Errorf("Quotes length exceeds limit: got %d", len(resp.Quotes))
	}
}

func TestService_Browse_WithEpisode(t *testing.T) {
	svc := getTestService()

	resp := svc.Browse("en", "10", 10, 0, 1, TruthAll)

	for i := 0; i < len(resp.Quotes); i++ {
		if resp.Quotes[i].Episode != 1 {
			t.Errorf("quote %d Episode: got %d, want 1", i, resp.Quotes[i].Episode)
		}
	}
}

func TestService_Browse_DefaultValues(t *testing.T) {
	svc := getTestService()

	resp := svc.Browse("", "", 0, -1, 0, TruthAll)

	if resp.Limit != 50 {
		t.Errorf("default limit: got %d, want 50", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Errorf("default offset: got %d, want 0", resp.Offset)
	}
}

func TestService_Browse_UnknownLang(t *testing.T) {
	svc := getTestService()

	resp := svc.Browse("fr", "10", 10, 0, 0, TruthAll)

	if resp.Total != 0 {
		t.Errorf("Total for unknown lang: got %d, want 0", resp.Total)
	}
}

func TestService_GetByCharacter(t *testing.T) {
	svc := getTestService()

	resp := svc.GetByCharacter("en", "27", 10, 0, 0, TruthAll)

	if resp.Total == 0 {
		t.Fatal("expected results for Beatrice")
	}
	if resp.Character != "Beatrice" {
		t.Errorf("Character: got %q, want %q", resp.Character, "Beatrice")
	}
	for i := 0; i < len(resp.Quotes); i++ {
		if resp.Quotes[i].CharacterID != "27" {
			t.Errorf("quote %d CharacterID: got %q, want %q", i, resp.Quotes[i].CharacterID, "27")
		}
	}
}

func TestService_GetByCharacter_WithEpisode(t *testing.T) {
	svc := getTestService()

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
	svc := getTestService()

	resp := svc.GetByCharacter("en", "999", 10, 0, 0, TruthAll)

	if resp.Total != 0 {
		t.Errorf("Total for unknown character: got %d, want 0", resp.Total)
	}
}

func TestService_GetByCharacter_DefaultValues(t *testing.T) {
	svc := getTestService()

	resp := svc.GetByCharacter("", "10", 0, -1, 0, TruthAll)

	if resp.Limit != 50 {
		t.Errorf("default limit: got %d, want 50", resp.Limit)
	}
}

func TestService_GetByAudioID(t *testing.T) {
	svc := getTestService()

	q := svc.GetByAudioID("en", "11900001")

	if q == nil {
		t.Fatal("expected to find quote by audio ID")
	}
	if q.CharacterID != "19" {
		t.Errorf("CharacterID: got %q, want %q", q.CharacterID, "19")
	}
}

func TestService_GetByAudioID_NotFound(t *testing.T) {
	svc := getTestService()

	q := svc.GetByAudioID("en", "99999999")

	if q != nil {
		t.Errorf("expected nil for unknown audio ID, got %+v", q)
	}
}

func TestService_GetByAudioID_DefaultLang(t *testing.T) {
	svc := getTestService()

	q := svc.GetByAudioID("", "11900001")

	if q == nil {
		t.Fatal("expected to find quote with empty lang (should default to en)")
	}
}

func TestService_GetByAudioID_UnknownLang(t *testing.T) {
	svc := getTestService()

	q := svc.GetByAudioID("fr", "11900001")

	if q != nil {
		t.Errorf("expected nil for unknown lang, got %+v", q)
	}
}

func TestService_Random(t *testing.T) {
	svc := getTestService()

	q := svc.Random("en", "", 0, TruthAll)

	if q == nil {
		t.Fatal("expected a random quote")
	}
	if q.CharacterID == "narrator" {
		t.Error("Random with no filters should exclude narrator")
	}
}

func TestService_Random_WithCharacter(t *testing.T) {
	svc := getTestService()

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
	svc := getTestService()

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
	svc := getTestService()

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
	svc := getTestService()

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
	svc := getTestService()

	q := svc.Random("", "", 0, TruthAll)

	if q == nil {
		t.Fatal("expected a random quote with default lang")
	}
}

func TestService_Random_UnknownLang(t *testing.T) {
	svc := getTestService()

	q := svc.Random("fr", "", 0, TruthAll)

	if q != nil {
		t.Errorf("expected nil for unknown lang, got %+v", q)
	}
}

func TestService_GetCharacters(t *testing.T) {
	svc := getTestService()

	chars := svc.GetCharacters()

	if len(chars) == 0 {
		t.Fatal("expected characters map to be non-empty")
	}
	if chars["10"] != "Battler" {
		t.Errorf("chars[10]: got %q, want %q", chars["10"], "Battler")
	}
	if chars["27"] != "Beatrice" {
		t.Errorf("chars[27]: got %q, want %q", chars["27"], "Beatrice")
	}
}

func TestService_GetStats(t *testing.T) {
	svc := getTestService()

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
	svc := getTestService()

	resp := svc.Browse("en", "", 10, 0, 0, TruthRed)

	for i := 0; i < len(resp.Quotes); i++ {
		if !strings.Contains(resp.Quotes[i].TextHtml, "red-truth") {
			t.Errorf("quote %d should contain red-truth in HTML", i)
		}
	}
}

func TestService_GetByCharacter_BlueTruthFilter(t *testing.T) {
	svc := getTestService()

	resp := svc.GetByCharacter("en", "10", 100, 0, 0, TruthBlue)

	for i := 0; i < len(resp.Quotes); i++ {
		if !strings.Contains(resp.Quotes[i].TextHtml, "blue-truth") {
			t.Errorf("quote %d should contain blue-truth in HTML", i)
		}
	}
}

func TestService_GetByAudioID_CompositeAudioID(t *testing.T) {
	svc := getTestService()

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
