package quote

import (
	"strings"
	"testing"
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/character"
	"umineko_quote/internal/quote/language"
	quoteparams "umineko_quote/internal/quote/params"
)

var testService = NewService()

func searchWithParams(
	svc Service,
	query string,
	lang language.Language,
	limit int,
	offset int,
	characterFilter character.Character,
	episode int,
	truth Truth,
	interactionA string,
	interactionB string,
) dto.SearchResponse {
	return svc.Search(quoteparams.SearchParams{
		Query:        query,
		Lang:         lang,
		Limit:        limit,
		Offset:       offset,
		Character:    characterFilter,
		Episode:      episode,
		Truth:        string(truth),
		InteractionA: interactionA,
		InteractionB: interactionB,
	})
}

func browseWithParams(
	svc Service,
	lang language.Language,
	characterFilter character.Character,
	limit int,
	offset int,
	episode int,
	truth Truth,
	interactionA string,
	interactionB string,
) dto.CharacterResponse {
	return svc.Browse(quoteparams.BrowseParams{
		Lang:         lang,
		Character:    characterFilter,
		Limit:        limit,
		Offset:       offset,
		Episode:      episode,
		Truth:        string(truth),
		InteractionA: interactionA,
		InteractionB: interactionB,
	})
}

func TestService_Search_ExactMatch(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "Beatrice", language.English, 10, 0, "", 0, TruthAll, "", "")

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

	resp := searchWithParams(svc, "witch", language.English, 0, -1, "", 0, TruthAll, "", "")

	if resp.Limit != 30 {
		t.Errorf("default limit: got %d, want 30", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Errorf("default offset: got %d, want 0", resp.Offset)
	}
}

func TestService_Search_WithCharacterFilter(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "witch", language.English, 10, 0, character.Battler, 0, TruthAll, "", "")

	for i := 0; i < len(resp.Results); i++ {
		if resp.Results[i].Quote.CharacterID != "10" {
			t.Errorf("result %d CharacterID: got %q, want %q", i, resp.Results[i].Quote.CharacterID, "10")
		}
	}
}

func TestService_Search_WithEpisodeFilter(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "witch", language.English, 10, 0, "", 1, TruthAll, "", "")

	for i := 0; i < len(resp.Results); i++ {
		if resp.Results[i].Quote.Episode != 1 {
			t.Errorf("result %d Episode: got %d, want 1", i, resp.Results[i].Quote.Episode)
		}
	}
}

func TestService_Search_RedTruthFilter(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "truth", language.English, 10, 0, "", 0, TruthRed, "", "")

	for i := 0; i < len(resp.Results); i++ {
		if !strings.Contains(resp.Results[i].Quote.TextHtml, "red-truth") {
			t.Errorf("result %d should contain red-truth in HTML", i)
		}
	}
}

func TestService_Search_NoResults(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "xyzzyxyzzyxyzzy", language.English, 10, 0, "", 0, TruthAll, "", "")

	if resp.Total != 0 {
		t.Errorf("Total: got %d, want 0", resp.Total)
	}
	if len(resp.Results) != 0 {
		t.Errorf("Results: got %d, want 0", len(resp.Results))
	}
}

func TestService_Search_Japanese(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "ベアトリーチェ", language.Japanese, 10, 0, "", 0, TruthAll, "", "")

	if resp.Total == 0 {
		t.Fatal("expected Japanese search results")
	}
}

func TestService_Search_UnknownLang(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "test", language.Language("fr"), 10, 0, "", 0, TruthAll, "", "")

	if resp.Total != 0 {
		t.Errorf("Total for unknown lang: got %d, want 0", resp.Total)
	}
}

func TestService_Search_WithInteractionPair(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "the", language.English, 100, 0, "", 0, TruthAll, "46", "47")
	if resp.Total == 0 {
		t.Fatal("expected interaction-filtered search results for Erika/Dlanor")
	}

	for i := 0; i < len(resp.Results); i++ {
		charID := resp.Results[i].Quote.CharacterID
		if charID != "46" && charID != "47" {
			t.Errorf("result %d CharacterID: got %q, want Erika(46) or Dlanor(47)", i, charID)
		}
	}
}

func TestService_Search_WithInteractionPair_SameCharacter(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "the", language.English, 100, 0, "", 0, TruthAll, "46", "46")
	if resp.Total != 0 {
		t.Errorf("same-character interaction filter should yield no results, got %d", resp.Total)
	}
}

func TestService_Browse(t *testing.T) {
	svc := testService

	resp := browseWithParams(svc, language.English, character.Battler, 10, 0, 0, TruthAll, "", "")

	if resp.Total == 0 {
		t.Fatal("expected browse results for Battler")
	}
	if resp.CharacterID != "10" {
		t.Errorf("CharacterID: got %q, want %q", resp.CharacterID, "10")
	}
	if resp.Character != character.CharacterNames[character.Battler] {
		t.Errorf("Character: got %q, want %q", resp.Character, character.CharacterNames[character.Battler])
	}
	if len(resp.Quotes) > 10 {
		t.Errorf("Quotes length exceeds limit: got %d", len(resp.Quotes))
	}
}

func TestService_Browse_WithEpisode(t *testing.T) {
	svc := testService

	resp := browseWithParams(svc, language.English, character.Battler, 10, 0, 1, TruthAll, "", "")

	for i := 0; i < len(resp.Quotes); i++ {
		if resp.Quotes[i].Episode != 1 {
			t.Errorf("quote %d Episode: got %d, want 1", i, resp.Quotes[i].Episode)
		}
	}
}

func TestService_Browse_DefaultValues(t *testing.T) {
	svc := testService

	resp := browseWithParams(svc, language.English, "", 0, -1, 0, TruthAll, "", "")

	if resp.Limit != 50 {
		t.Errorf("default limit: got %d, want 50", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Errorf("default offset: got %d, want 0", resp.Offset)
	}
}

func TestService_Browse_UnknownLang(t *testing.T) {
	svc := testService

	resp := browseWithParams(svc, language.Language("fr"), character.Battler, 10, 0, 0, TruthAll, "", "")

	if resp.Total != 0 {
		t.Errorf("Total for unknown lang: got %d, want 0", resp.Total)
	}
}

func TestService_Browse_WithInteractionPair(t *testing.T) {
	svc := testService

	resp := browseWithParams(svc, language.English, "", 100, 0, 0, TruthAll, "46", "47")
	if resp.Total == 0 {
		t.Fatal("expected browse results for Erika/Dlanor interaction pair")
	}
	for i := 0; i < len(resp.Quotes); i++ {
		charID := resp.Quotes[i].CharacterID
		if charID != "46" && charID != "47" {
			t.Errorf("quote %d CharacterID: got %q, want Erika(46) or Dlanor(47)", i, charID)
		}
	}
}

func TestService_Browse_WithInteractionPairAndCharacter(t *testing.T) {
	svc := testService

	resp := browseWithParams(svc, language.English, character.Erika, 100, 0, 0, TruthAll, "46", "47")
	if resp.Total == 0 {
		t.Fatal("expected browse results for Erika within Erika/Dlanor interactions")
	}
	for i := 0; i < len(resp.Quotes); i++ {
		if resp.Quotes[i].CharacterID != "46" {
			t.Errorf("quote %d CharacterID: got %q, want Erika(46)", i, resp.Quotes[i].CharacterID)
		}
	}
}

func TestService_GetByAudioID(t *testing.T) {
	svc := testService

	q := svc.GetByAudioID(language.English, "11900001")

	if q == nil {
		t.Fatal("expected to find quote by audio ID")
	}
	if q.CharacterID != "19" {
		t.Errorf("CharacterID: got %q, want %q", q.CharacterID, "19")
	}
}

func TestService_GetByAudioID_NotFound(t *testing.T) {
	svc := testService

	q := svc.GetByAudioID(language.English, "99999999")

	if q != nil {
		t.Errorf("expected nil for unknown audio ID, got %+v", q)
	}
}

func TestService_GetByAudioID_EmptyLang(t *testing.T) {
	svc := testService

	q := svc.GetByAudioID(language.Language(""), "11900001")

	if q != nil {
		t.Errorf("expected nil for empty lang (defaulting happens in controller), got %+v", q)
	}
}

func TestService_GetByAudioID_UnknownLang(t *testing.T) {
	svc := testService

	q := svc.GetByAudioID(language.Language("fr"), "11900001")

	if q != nil {
		t.Errorf("expected nil for unknown lang, got %+v", q)
	}
}

func TestService_Random(t *testing.T) {
	svc := testService

	q := svc.Random(language.English, "", 0, TruthAll)

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
		q := svc.Random(language.English, character.Beatrice, 0, TruthAll)
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
		q := svc.Random(language.English, "", 1, TruthAll)
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
		q := svc.Random(language.English, character.Battler, 1, TruthAll)
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
		q := svc.Random(language.English, "", 0, TruthRed)
		if q == nil {
			t.Fatal("expected a random red truth quote")
		}
		if !strings.Contains(q.TextHtml, "red-truth") {
			t.Errorf("expected red-truth in HTML: %q", q.TextHtml)
		}
	}
}

func TestService_Random_EmptyLang(t *testing.T) {
	svc := testService

	q := svc.Random(language.Language(""), "", 0, TruthAll)

	if q != nil {
		t.Errorf("expected nil for empty lang (defaulting happens in controller), got %+v", q)
	}
}

func TestService_Random_UnknownLang(t *testing.T) {
	svc := testService

	q := svc.Random(language.Language("fr"), "", 0, TruthAll)

	if q != nil {
		t.Errorf("expected nil for unknown lang, got %+v", q)
	}
}

func TestService_GetContext(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "Beatrice", language.English, 10, 0, "", 0, TruthAll, "", "")
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
	parts := strings.SplitN(midAudioId, ", ", 2)
	midAudioId = parts[0]

	result := svc.GetContext(language.English, midAudioId, 5)

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

	result := svc.GetContext(language.English, "99999999", 5)

	if result != nil {
		t.Errorf("expected nil for unknown audio ID, got %+v", result)
	}
}

func TestService_GetContext_EdgeOfSlice(t *testing.T) {
	svc := testService

	result := svc.GetContext(language.English, "11900001", 5)

	if result == nil {
		t.Fatal("expected context result for edge-of-slice quote")
	}
	if len(result.Before) > 5 {
		t.Errorf("Before length exceeds lines: got %d", len(result.Before))
	}
	if len(result.After) > 5 {
		t.Errorf("After length exceeds lines: got %d", len(result.After))
	}
}

func TestService_GetContext_DefaultLines(t *testing.T) {
	svc := testService

	result := svc.GetContext(language.English, "11900001", 0)

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

	result := svc.GetContext(language.English, "11900001", 100)

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

	result := svc.GetContext(language.Language("fr"), "11900001", 5)

	if result != nil {
		t.Errorf("expected nil for unknown lang, got %+v", result)
	}
}

func TestService_GetContext_EmptyLang(t *testing.T) {
	svc := testService

	result := svc.GetContext(language.Language(""), "11900001", 5)

	if result != nil {
		t.Errorf("expected nil for empty lang (defaulting happens in controller), got %+v", result)
	}
}

func TestService_GetCharacters(t *testing.T) {
	svc := testService

	chars := svc.GetCharacters()

	if len(chars) == 0 {
		t.Fatal("expected characters map to be non-empty")
	}
	if chars[character.Battler] != character.CharacterNames[character.Battler] {
		t.Errorf("chars[battler]: got %q, want %q", chars[character.Battler], character.CharacterNames[character.Battler])
	}
	if chars[character.Beatrice] != character.CharacterNames[character.Beatrice] {
		t.Errorf("chars[beatrice]: got %q, want %q", chars[character.Beatrice], character.CharacterNames[character.Beatrice])
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

	resp := browseWithParams(svc, language.English, "", 10, 0, 0, TruthRed, "", "")

	for i := 0; i < len(resp.Quotes); i++ {
		if !strings.Contains(resp.Quotes[i].TextHtml, "red-truth") {
			t.Errorf("quote %d should contain red-truth in HTML", i)
		}
	}
}

func TestService_Browse_BlueTruthFilter(t *testing.T) {
	svc := testService

	resp := browseWithParams(svc, language.English, character.Battler, 100, 0, 0, TruthBlue, "", "")

	for i := 0; i < len(resp.Quotes); i++ {
		if !strings.Contains(resp.Quotes[i].TextHtml, "blue-truth") {
			t.Errorf("quote %d should contain blue-truth in HTML", i)
		}
	}
}

func TestService_GetByAudioID_CompositeAudioID(t *testing.T) {
	svc := testService

	q1 := svc.GetByAudioID(language.English, "11900001")
	if q1 == nil {
		t.Fatal("expected to find quote by first audio ID in composite")
	}

	q2 := svc.GetByAudioID(language.English, "11900002")
	if q2 == nil {
		t.Fatal("expected to find quote by second audio ID in composite")
	}

	if q1.CharacterID != q2.CharacterID {
		t.Errorf("both audio IDs should resolve to same character: %q vs %q", q1.CharacterID, q2.CharacterID)
	}
}
