package ciconia

import (
	"strings"
	"testing"
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"
	quoteparams "umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/store"

	cicharacter "github.com/VictoriqueMoe/umineko_script_parser/ciconia/character"
)

var testService = NewService()

func searchWithParams(
	svc Service,
	query string,
	lang language.Language,
	limit int,
	offset int,
	characterID string,
	chapter string,
	interactionA string,
	interactionB string,
) dto.SearchResponse {
	return svc.Search(quoteparams.SearchParams{
		Query:        query,
		Lang:         lang,
		Limit:        limit,
		Offset:       offset,
		CharacterID:  characterID,
		InteractionA: interactionA,
		InteractionB: interactionB,
	}, chapter)
}

func browseWithParams(
	svc Service,
	lang language.Language,
	characterID string,
	limit int,
	offset int,
	chapter string,
	interactionA string,
	interactionB string,
) dto.CharacterResponse {
	return svc.Browse(quoteparams.BrowseParams{
		Lang:         lang,
		CharacterID:  characterID,
		Limit:        limit,
		Offset:       offset,
		InteractionA: interactionA,
		InteractionB: interactionB,
	}, chapter)
}

func TestService_Search_ExactMatch(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "Miyao", language.English, 10, 0, "", "", "", "")

	if resp.Total == 0 {
		t.Fatal("expected search results for 'Miyao'")
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

	resp := searchWithParams(svc, "the", language.English, 0, -1, "", "", "", "")

	if resp.Limit != 30 {
		t.Errorf("default limit: got %d, want 30", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Errorf("default offset: got %d, want 0", resp.Offset)
	}
}

func TestService_Search_WithCharacterFilter(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "the", language.English, 10, 0, cicharacter.Miyao.ID(), "", "", "")

	for i := 0; i < len(resp.Results); i++ {
		q := resp.Results[i].Quote.(dto.CiconiaQuote)
		if q.CharacterID != cicharacter.Miyao.ID() {
			t.Errorf("result %d CharacterID: got %q, want %q", i, q.CharacterID, cicharacter.Miyao.ID())
		}
	}
}

func TestService_Search_WithChapterFilter(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "the", language.English, 10, 0, "", "01", "", "")

	if resp.Total == 0 {
		t.Fatal("expected results for chapter 01")
	}
	for i := 0; i < len(resp.Results); i++ {
		q := resp.Results[i].Quote.(dto.CiconiaQuote)
		if q.Chapter != "01" {
			t.Errorf("result %d Chapter: got %q, want 01", i, q.Chapter)
		}
	}
}

func TestService_Search_NoResults(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "xyzzyxyzzyxyzzy", language.English, 10, 0, "", "", "", "")

	if resp.Total != 0 {
		t.Errorf("Total: got %d, want 0", resp.Total)
	}
}

func TestService_Browse(t *testing.T) {
	svc := testService

	resp := browseWithParams(svc, language.English, cicharacter.Miyao.ID(), 10, 0, "", "", "")

	if resp.Total == 0 {
		t.Fatal("expected browse results for Miyao")
	}
	if resp.CharacterID != cicharacter.Miyao.ID() {
		t.Errorf("CharacterID: got %q, want %q", resp.CharacterID, cicharacter.Miyao.ID())
	}
	quotes := resp.Quotes.([]dto.CiconiaQuote)
	if len(quotes) > 10 {
		t.Errorf("Quotes length exceeds limit: got %d", len(quotes))
	}
}

func TestService_Browse_WithChapterFilter(t *testing.T) {
	svc := testService

	resp := browseWithParams(svc, language.English, "", 10, 0, "df01", "", "")

	if resp.Total == 0 {
		t.Fatal("expected browse results for df01")
	}
	quotes := resp.Quotes.([]dto.CiconiaQuote)
	for i := 0; i < len(quotes); i++ {
		if quotes[i].Chapter != "df01" {
			t.Errorf("quote %d Chapter: got %q, want df01", i, quotes[i].Chapter)
		}
	}
}

func TestService_Random(t *testing.T) {
	svc := testService

	q := svc.Random(language.English, "", "")

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
		q := svc.Random(language.English, cicharacter.Jayden.ID(), "")
		if q == nil {
			t.Fatal("expected a random Jayden quote")
		}
		if q.CharacterID != cicharacter.Jayden.ID() {
			t.Errorf("CharacterID: got %q, want %q", q.CharacterID, cicharacter.Jayden.ID())
		}
	}
}

func TestService_Random_WithChapter(t *testing.T) {
	svc := testService

	for i := 0; i < 10; i++ {
		q := svc.Random(language.English, "", "05")
		if q == nil {
			t.Fatal("expected a random chapter 05 quote")
		}
		if q.Chapter != "05" {
			t.Errorf("Chapter: got %q, want 05", q.Chapter)
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

	result := svc.GetCharacters()

	if len(result.Characters) == 0 {
		t.Fatal("expected characters map to be non-empty")
	}
	miyaoID := cicharacter.Miyao.ID()
	expectedName := cicharacter.MainCharacterNames[cicharacter.Miyao]
	if result.Characters[miyaoID] != expectedName {
		t.Errorf("chars[Miyao]: got %q, want %q", result.Characters[miyaoID], expectedName)
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
		t.Fatal("expected Compute to return non-nil")
	}

	cStats, ok := result.(*dto.CiconiaStatsResult)
	if !ok {
		t.Fatalf("expected *dto.CiconiaStatsResult, got %T", result)
	}
	if len(cStats.TopSpeakers) == 0 {
		t.Error("expected non-empty TopSpeakers")
	}
	if len(cStats.Interactions) == 0 {
		t.Error("expected non-empty Interactions")
	}
}

func TestService_Search_HasJapaneseText(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "Miyao", language.English, 5, 0, "", "", "", "")

	if resp.Total == 0 {
		t.Fatal("expected results")
	}

	q := resp.Results[0].Quote.(dto.CiconiaQuote)
	if q.TextJP == "" {
		t.Error("expected TextJP to be populated")
	}
}

func TestService_Search_HasChapter(t *testing.T) {
	svc := testService

	resp := searchWithParams(svc, "the", language.English, 5, 0, "", "", "", "")

	if resp.Total == 0 {
		t.Fatal("expected results")
	}

	q := resp.Results[0].Quote.(dto.CiconiaQuote)
	if q.Chapter == "" {
		t.Error("expected Chapter to be populated")
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

func TestService_QuoteCount(t *testing.T) {
	svc := testService

	langs := svc.LoadedLanguages()
	count := langs[language.English]

	t.Logf("Ciconia quotes loaded: %d", count)

	if count < 5000 {
		t.Errorf("expected at least 5000 quotes, got %d", count)
	}
}

func TestService_AllChaptersPresent(t *testing.T) {
	svc := testService

	expectedChapters := []string{"00", "01", "05", "10", "15", "20", "25", "25b", "df01", "df08", "df16"}

	for _, ch := range expectedChapters {
		resp := searchWithParams(svc, "the", language.English, 1, 0, "", ch, "", "")
		if resp.Total == 0 {
			t.Errorf("no quotes found for chapter %q", ch)
		}
	}
}

func TestStats_TopSpeakersHaveMiyao(t *testing.T) {
	svc := testService

	result := svc.GetStats().Compute(store.AllEpisodes).(*dto.CiconiaStatsResult)

	found := false
	for i := 0; i < len(result.TopSpeakers); i++ {
		if strings.Contains(result.TopSpeakers[i].Name, "Miyao") {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected Miyao in TopSpeakers")
	}
}

func TestStats_InteractionsExist(t *testing.T) {
	svc := testService

	result := svc.GetStats().Compute(store.AllEpisodes).(*dto.CiconiaStatsResult)

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

func TestStats_LinesPerChapterExists(t *testing.T) {
	svc := testService

	result := svc.GetStats().Compute(store.AllEpisodes).(*dto.CiconiaStatsResult)

	if len(result.LinesPerChapter) == 0 {
		t.Fatal("expected LinesPerChapter to be populated")
	}

	if _, ok := result.LinesPerChapter["01"]; !ok {
		t.Error("expected chapter 01 in LinesPerChapter")
	}
}

func TestService_KeropoyoPresent(t *testing.T) {
	svc := testService

	q := svc.Random(language.English, cicharacter.Keropoyo.ID(), "")
	if q == nil {
		t.Fatal("expected a random Keropoyo quote")
	}
	if q.CharacterID != cicharacter.Keropoyo.ID() {
		t.Errorf("CharacterID: got %q, want %q", q.CharacterID, cicharacter.Keropoyo.ID())
	}
}
