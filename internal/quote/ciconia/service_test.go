package ciconia

import (
	"log"
	"os"
	"strings"
	"testing"

	"umineko_quote/internal/dao"
	"umineko_quote/internal/dao/memory"
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/game"
	"umineko_quote/internal/quote/language"
	quoteparams "umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/quotes"
	"umineko_quote/internal/repository"

	cicharacter "github.com/VictoriqueMoe/umineko_script_parser/ciconia/character"
)

var testService Service

func TestMain(m *testing.M) {
	loaded, err := game.Ciconia{}.Load()
	if err != nil {
		log.Fatalf("Load: %v", err)
	}

	repos := repository.NewFromDAOs(dao.DAOs{Ciconia: memory.NewQuotes(game.Ciconia{}, loaded)}, nil)

	svc, err := NewService(repos.Ciconia)
	if err != nil {
		log.Fatalf("NewService: %v", err)
	}
	testService = svc

	os.Exit(m.Run())
}

func searchWithParams(
	t *testing.T,
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
	t.Helper()

	resp, err := svc.Search(quoteparams.SearchParams{
		Query:        query,
		Lang:         lang,
		Limit:        limit,
		Offset:       offset,
		CharacterID:  characterID,
		InteractionA: interactionA,
		InteractionB: interactionB,
	}, quoteparams.QuoteScope{Segment: chapter})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	return resp
}

func browseWithParams(
	t *testing.T,
	svc Service,
	lang language.Language,
	characterID string,
	limit int,
	offset int,
	chapter string,
	interactionA string,
	interactionB string,
) dto.CharacterResponse {
	t.Helper()

	resp, err := svc.Browse(quoteparams.BrowseParams{
		Lang:         lang,
		CharacterID:  characterID,
		Limit:        limit,
		Offset:       offset,
		InteractionA: interactionA,
		InteractionB: interactionB,
	}, quoteparams.QuoteScope{Segment: chapter})
	if err != nil {
		t.Fatalf("Browse: %v", err)
	}
	return resp
}

func randomQuote(t *testing.T, characterID string, scope quoteparams.QuoteScope) *dto.CiconiaQuote {
	t.Helper()

	q, err := testService.Random(language.English, characterID, 0, scope)
	if err != nil {
		t.Fatalf("Random: %v", err)
	}
	return q
}

func loadedLanguages(t *testing.T) map[language.Language]int {
	t.Helper()

	langs, err := testService.LoadedLanguages()
	if err != nil {
		t.Fatalf("LoadedLanguages: %v", err)
	}
	return langs
}

func TestService_Search_ExactMatch(t *testing.T) {
	svc := testService

	resp := searchWithParams(t, svc, "Miyao", language.English, 10, 0, "", "", "", "")

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

	resp := searchWithParams(t, svc, "the", language.English, 0, -1, "", "", "", "")

	if resp.Limit != 30 {
		t.Errorf("default limit: got %d, want 30", resp.Limit)
	}
	if resp.Offset != 0 {
		t.Errorf("default offset: got %d, want 0", resp.Offset)
	}
}

func TestService_Search_WithCharacterFilter(t *testing.T) {
	svc := testService

	resp := searchWithParams(t, svc, "the", language.English, 10, 0, cicharacter.Miyao.ID(), "", "", "")

	for i, result := range resp.Results {
		q := result.Quote.(dto.CiconiaQuote)
		if q.CharacterID != cicharacter.Miyao.ID() {
			t.Errorf("result %d CharacterID: got %q, want %q", i, q.CharacterID, cicharacter.Miyao.ID())
		}
	}
}

func TestService_Search_WithChapterFilter(t *testing.T) {
	svc := testService

	resp := searchWithParams(t, svc, "the", language.English, 10, 0, "", "01", "", "")

	if resp.Total == 0 {
		t.Fatal("expected results for chapter 01")
	}
	for i, result := range resp.Results {
		q := result.Quote.(dto.CiconiaQuote)
		if q.Chapter != "01" {
			t.Errorf("result %d Chapter: got %q, want 01", i, q.Chapter)
		}
	}
}

func TestService_Search_NoResults(t *testing.T) {
	svc := testService

	resp := searchWithParams(t, svc, "xyzzyxyzzyxyzzy", language.English, 10, 0, "", "", "", "")

	if resp.Total != 0 {
		t.Errorf("Total: got %d, want 0", resp.Total)
	}
}

func TestService_Browse(t *testing.T) {
	svc := testService

	resp := browseWithParams(t, svc, language.English, cicharacter.Miyao.ID(), 10, 0, "", "", "")

	if resp.Total == 0 {
		t.Fatal("expected browse results for Miyao")
	}
	if resp.CharacterID != cicharacter.Miyao.ID() {
		t.Errorf("CharacterID: got %q, want %q", resp.CharacterID, cicharacter.Miyao.ID())
	}
	browsed := resp.Quotes.([]dto.CiconiaQuote)
	if len(browsed) > 10 {
		t.Errorf("Quotes length exceeds limit: got %d", len(browsed))
	}
}

func TestService_Browse_WithChapterFilter(t *testing.T) {
	svc := testService

	resp := browseWithParams(t, svc, language.English, "", 10, 0, "df01", "", "")

	if resp.Total == 0 {
		t.Fatal("expected browse results for df01")
	}
	for i, q := range resp.Quotes.([]dto.CiconiaQuote) {
		if q.Chapter != "df01" {
			t.Errorf("quote %d Chapter: got %q, want df01", i, q.Chapter)
		}
	}
}

func TestService_Random(t *testing.T) {
	q := randomQuote(t, "", quoteparams.QuoteScope{})

	if q == nil {
		t.Fatal("expected a random quote")
	}
	if q.CharacterID == game.NarratorID {
		t.Error("Random with no filters should exclude narrator")
	}
}

func TestService_Random_WithCharacter(t *testing.T) {
	for range 10 {
		q := randomQuote(t, cicharacter.Jayden.ID(), quoteparams.QuoteScope{})
		if q == nil {
			t.Fatal("expected a random Jayden quote")
		}
		if q.CharacterID != cicharacter.Jayden.ID() {
			t.Errorf("CharacterID: got %q, want %q", q.CharacterID, cicharacter.Jayden.ID())
		}
	}
}

func TestService_Random_WithChapter(t *testing.T) {
	for range 10 {
		q := randomQuote(t, "", quoteparams.QuoteScope{Segment: "05"})
		if q == nil {
			t.Fatal("expected a random chapter 05 quote")
		}
		if q.Chapter != "05" {
			t.Errorf("Chapter: got %q, want 05", q.Chapter)
		}
	}
}

func TestService_GetByIndex(t *testing.T) {
	q, err := testService.GetByIndex(language.English, 0)
	if err != nil {
		t.Fatalf("GetByIndex: %v", err)
	}

	if q == nil {
		t.Fatal("expected quote at index 0")
	}
	if q.Index != 0 {
		t.Errorf("Index: got %d, want 0", q.Index)
	}
}

func TestService_GetByIndex_NotFound(t *testing.T) {
	q, err := testService.GetByIndex(language.English, 999999)
	if err != nil {
		t.Fatalf("GetByIndex: %v", err)
	}

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

	result := stats.Compute(quotes.AllEpisodes)
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

	resp := searchWithParams(t, svc, "Miyao", language.English, 5, 0, "", "", "", "")

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

	resp := searchWithParams(t, svc, "the", language.English, 5, 0, "", "", "", "")

	if resp.Total == 0 {
		t.Fatal("expected results")
	}

	q := resp.Results[0].Quote.(dto.CiconiaQuote)
	if q.Chapter == "" {
		t.Error("expected Chapter to be populated")
	}
}

func TestService_LoadedLanguages(t *testing.T) {
	langs := loadedLanguages(t)

	if len(langs) == 0 {
		t.Fatal("expected at least one loaded language")
	}
	if langs[language.English] == 0 {
		t.Error("expected English to have quotes")
	}
}

func TestService_QuoteCount(t *testing.T) {
	count := loadedLanguages(t)[language.English]

	t.Logf("Ciconia quotes loaded: %d", count)

	if count < 5000 {
		t.Errorf("expected at least 5000 quotes, got %d", count)
	}
}

func TestService_AllChaptersPresent(t *testing.T) {
	svc := testService

	expectedChapters := []string{"00", "01", "05", "10", "15", "20", "25", "25b", "df01", "df08", "df16"}

	for _, ch := range expectedChapters {
		resp := searchWithParams(t, svc, "the", language.English, 1, 0, "", ch, "", "")
		if resp.Total == 0 {
			t.Errorf("no quotes found for chapter %q", ch)
		}
	}
}

func TestStats_TopSpeakersHaveMiyao(t *testing.T) {
	svc := testService

	result := svc.GetStats().Compute(quotes.AllEpisodes).(*dto.CiconiaStatsResult)

	found := false
	for _, speaker := range result.TopSpeakers {
		if strings.Contains(speaker.Name, "Miyao") {
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

	result := svc.GetStats().Compute(quotes.AllEpisodes).(*dto.CiconiaStatsResult)

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

	result := svc.GetStats().Compute(quotes.AllEpisodes).(*dto.CiconiaStatsResult)

	if len(result.LinesPerChapter) == 0 {
		t.Fatal("expected LinesPerChapter to be populated")
	}

	if _, ok := result.LinesPerChapter["01"]; !ok {
		t.Error("expected chapter 01 in LinesPerChapter")
	}
}

func TestService_KeropoyoPresent(t *testing.T) {
	q := randomQuote(t, cicharacter.Keropoyo.ID(), quoteparams.QuoteScope{})
	if q == nil {
		t.Fatal("expected a random Keropoyo quote")
	}
	if q.CharacterID != cicharacter.Keropoyo.ID() {
		t.Errorf("CharacterID: got %q, want %q", q.CharacterID, cicharacter.Keropoyo.ID())
	}
}
