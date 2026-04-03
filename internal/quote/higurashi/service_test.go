package higurashi

import (
	"strings"
	"testing"
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"
	quoteparams "umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/store"

	hicharacter "github.com/VictoriqueMoe/umineko_script_parser/higurashi/character"
)

var testService = NewService()

func searchWithParams(
	svc Service,
	query string,
	lang language.Language,
	limit int,
	offset int,
	characterID string,
	episode int,
	arc string,
	interactionA string,
	interactionB string,
) dto.SearchResponse {
	return svc.Search(quoteparams.SearchParams{
		Query:        query,
		Lang:         lang,
		Limit:        limit,
		Offset:       offset,
		CharacterID:  characterID,
		Episode:      episode,
		InteractionA: interactionA,
		InteractionB: interactionB,
	}, arc)
}

func browseWithParams(
	svc Service,
	lang language.Language,
	characterID string,
	limit int,
	offset int,
	episode int,
	arc string,
	interactionA string,
	interactionB string,
) dto.CharacterResponse {
	return svc.Browse(quoteparams.BrowseParams{
		Lang:         lang,
		CharacterID:  characterID,
		Limit:        limit,
		Offset:       offset,
		Episode:      episode,
		InteractionA: interactionA,
		InteractionB: interactionB,
	}, arc)
}

func TestService_Search_ExactMatch(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "Rena", language.English, 10, 0, "", 0, "", "", "")

	if resp.Total == 0 {
		t.Fatal("expected search results for 'Rena'")
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

	resp := searchWithParams(svc, "school", language.English, 0, -1, "", 0, "", "", "")

	if resp.Limit != 30 {
		t.Errorf("default limit: got %d, want 30", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Errorf("default offset: got %d, want 0", resp.Offset)
	}
}

func TestService_Search_WithCharacterFilter(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "school", language.English, 10, 0, hicharacter.Keiichi.ID(), 0, "", "", "")

	for i := 0; i < len(resp.Results); i++ {
		q := resp.Results[i].Quote.(dto.HigurashiQuote)
		if q.CharacterID != hicharacter.Keiichi.ID() {
			t.Errorf("result %d CharacterID: got %q, want %q", i, q.CharacterID, hicharacter.Keiichi.ID())
		}
	}
}

func TestService_Search_WithArcFilter(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "the", language.English, 10, 0, "", 0, "onikakushi", "", "")

	if resp.Total == 0 {
		t.Fatal("expected results for onikakushi arc")
	}
	for i := 0; i < len(resp.Results); i++ {
		q := resp.Results[i].Quote.(dto.HigurashiQuote)
		if q.Arc != "onikakushi" {
			t.Errorf("result %d Arc: got %q, want onikakushi", i, q.Arc)
		}
	}
}

func TestService_Search_WithEpisodeFilter(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "the", language.English, 10, 0, "", 1, "", "", "")

	for i := 0; i < len(resp.Results); i++ {
		q := resp.Results[i].Quote.(dto.HigurashiQuote)
		if q.Episode != 1 {
			t.Errorf("result %d Episode: got %d, want 1", i, q.Episode)
		}
	}
}

func TestService_Search_NoResults(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "xyzzyxyzzyxyzzy", language.English, 10, 0, "", 0, "", "", "")

	if resp.Total != 0 {
		t.Errorf("Total: got %d, want 0", resp.Total)
	}
}

func TestService_Browse(t *testing.T) {
	svc := testService

	resp := browseWithParams(svc, language.English, hicharacter.Keiichi.ID(), 10, 0, 0, "", "", "")

	if resp.Total == 0 {
		t.Fatal("expected browse results for Keiichi")
	}
	if resp.CharacterID != hicharacter.Keiichi.ID() {
		t.Errorf("CharacterID: got %q, want %q", resp.CharacterID, hicharacter.Keiichi.ID())
	}
	quotes := resp.Quotes.([]dto.HigurashiQuote)
	if len(quotes) > 10 {
		t.Errorf("Quotes length exceeds limit: got %d", len(quotes))
	}
}

func TestService_Browse_WithArcFilter(t *testing.T) {
	svc := testService

	resp := browseWithParams(svc, language.English, "", 10, 0, 0, "watanagashi", "", "")

	if resp.Total == 0 {
		t.Fatal("expected browse results for watanagashi")
	}
	quotes := resp.Quotes.([]dto.HigurashiQuote)
	for i := 0; i < len(quotes); i++ {
		if quotes[i].Arc != "watanagashi" {
			t.Errorf("quote %d Arc: got %q, want watanagashi", i, quotes[i].Arc)
		}
	}
}

func TestService_Random(t *testing.T) {
	svc := testService

	q := svc.Random(language.English, "", 0, "")

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
		q := svc.Random(language.English, hicharacter.Rena.ID(), 0, "")
		if q == nil {
			t.Fatal("expected a random Rena quote")
		}
		if q.CharacterID != hicharacter.Rena.ID() {
			t.Errorf("CharacterID: got %q, want %q", q.CharacterID, hicharacter.Rena.ID())
		}
	}
}

func TestService_Random_WithArc(t *testing.T) {
	svc := testService

	for i := 0; i < 10; i++ {
		q := svc.Random(language.English, "", 0, "onikakushi")
		if q == nil {
			t.Fatal("expected a random onikakushi quote")
		}
		if q.Arc != "onikakushi" {
			t.Errorf("Arc: got %q, want onikakushi", q.Arc)
		}
	}
}

func TestService_GetByIndex(t *testing.T) {
	svc := testService

	q := svc.GetByIndex(language.English, 0)

	if q == nil {
		t.Fatal("expected quote at index 0")
	}
	if q.Index != 0 {
		t.Errorf("Index: got %d, want 0", q.Index)
	}
}

func TestService_GetByIndex_NotFound(t *testing.T) {
	svc := testService

	q := svc.GetByIndex(language.English, 999999)

	if q != nil {
		t.Errorf("expected nil for out-of-bounds index")
	}
}

func TestService_GetCharacters(t *testing.T) {
	svc := testService

	chars := svc.GetCharacters()

	if len(chars) == 0 {
		t.Fatal("expected characters map to be non-empty")
	}
	if chars[hicharacter.Keiichi] != hicharacter.CharacterNames[hicharacter.Keiichi] {
		t.Errorf("chars[Keiichi]: got %q, want %q", chars[hicharacter.Keiichi], hicharacter.CharacterNames[hicharacter.Keiichi])
	}
}

func TestService_GetStats(t *testing.T) {
	svc := testService

	stats := svc.GetStats()

	if stats == nil {
		t.Fatal("expected stats to be non-nil")
	}

	result := stats.Compute(store.AllEpisodes)
	if result == nil {
		t.Fatal("expected Compute(AllEpisodes) to return non-nil")
	}

	hStats, ok := result.(*dto.HigurashiStatsResult)
	if !ok {
		t.Fatalf("expected *dto.HigurashiStatsResult, got %T", result)
	}
	if len(hStats.TopSpeakers) == 0 {
		t.Error("expected non-empty TopSpeakers")
	}
	if len(hStats.Interactions) == 0 {
		t.Error("expected non-empty Interactions")
	}
}

func TestService_Search_HasJapaneseText(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "Keiichi", language.English, 5, 0, "", 0, "", "", "")

	if resp.Total == 0 {
		t.Fatal("expected results")
	}

	q := resp.Results[0].Quote.(dto.HigurashiQuote)
	if q.TextJP == "" {
		t.Error("expected TextJP to be populated")
	}
}

func TestService_Search_HasArc(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "morning", language.English, 5, 0, "", 0, "", "", "")

	if resp.Total == 0 {
		t.Fatal("expected results")
	}

	q := resp.Results[0].Quote.(dto.HigurashiQuote)
	if q.Arc == "" {
		t.Error("expected Arc to be populated")
	}
}

func TestService_LoadedLanguages(t *testing.T) {
	svc := testService

	langs := svc.LoadedLanguages()

	if len(langs) == 0 {
		t.Fatal("expected at least one loaded language")
	}
	if langs[language.English] == 0 {
		t.Error("expected English to have quotes")
	}
}

func TestService_Browse_WithInteractionPair(t *testing.T) {
	svc := testService

	resp := browseWithParams(svc, language.English, "", 100, 0, 0, "", hicharacter.Keiichi.ID(), hicharacter.Rena.ID())

	if resp.Total == 0 {
		t.Fatal("expected browse results for Keiichi/Rena interaction pair")
	}
	quotes := resp.Quotes.([]dto.HigurashiQuote)
	for i := 0; i < len(quotes); i++ {
		charID := quotes[i].CharacterID
		if charID != hicharacter.Keiichi.ID() && charID != hicharacter.Rena.ID() {
			t.Errorf("quote %d CharacterID: got %q, want Keiichi or Rena", i, charID)
		}
	}
}

func TestService_QuoteCount(t *testing.T) {
	svc := testService

	langs := svc.LoadedLanguages()
	count := langs[language.English]

	t.Logf("Higurashi quotes loaded: %d", count)

	if count < 50000 {
		t.Errorf("expected at least 50000 quotes, got %d", count)
	}
}

func TestService_AllArcsPresent(t *testing.T) {
	svc := testService

	expectedArcs := []string{"onikakushi", "watanagashi", "tatarigoroshi", "himatsubushi", "meakashi", "tsumihoroboshi", "minagoroshi", "matsuribayashi"}

	for _, arc := range expectedArcs {
		resp := searchWithParams(svc, "the", language.English, 1, 0, "", 0, arc, "", "")
		if resp.Total == 0 {
			t.Errorf("no quotes found for arc %q", arc)
		}
	}
}

func TestStats_TopSpeakersHaveKeiichi(t *testing.T) {
	svc := testService

	result := svc.GetStats().Compute(store.AllEpisodes).(*dto.HigurashiStatsResult)

	found := false
	for i := 0; i < len(result.TopSpeakers); i++ {
		if strings.Contains(result.TopSpeakers[i].Name, "Keiichi") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected Keiichi in TopSpeakers")
	}
}

func TestStats_InteractionsExist(t *testing.T) {
	svc := testService

	result := svc.GetStats().Compute(store.AllEpisodes).(*dto.HigurashiStatsResult)

	if len(result.Interactions) == 0 {
		t.Fatal("expected interactions")
	}

	top := result.Interactions[0]
	if top.Count == 0 {
		t.Error("top interaction should have non-zero count")
	}
	if top.NameA == "" || top.NameB == "" {
		t.Error("top interaction should have resolved names")
	}
}

func TestStats_LinesPerArcExists(t *testing.T) {
	svc := testService

	result := svc.GetStats().Compute(store.AllEpisodes).(*dto.HigurashiStatsResult)

	if len(result.LinesPerArc) == 0 {
		t.Fatal("expected LinesPerArc to be populated")
	}

	if _, ok := result.LinesPerArc["onikakushi"]; !ok {
		t.Error("expected onikakushi in LinesPerArc")
	}
}
