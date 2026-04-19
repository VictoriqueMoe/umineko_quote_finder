package ciconia

import (
	"testing"
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/store"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
)

func makeQuote(charID string, chapter string) dto.CiconiaQuote {
	return dto.CiconiaQuote{
		CicroniaQuote: scriptdto.CicroniaQuote{
			BaseQuote: scriptdto.BaseQuote{
				CharacterID: charID,
				Character:   charID,
				Episode:     1,
				ContentType: "chapter",
			},
			Chapter: chapter,
		},
	}
}

func TestCiconiaStats_TopSpeakers(t *testing.T) {
	quotes := []dto.CiconiaQuote{
		makeQuote("miyao", "01"),
		makeQuote("miyao", "01"),
		makeQuote("miyao", "01"),
		makeQuote("jayden", "01"),
		makeQuote("jayden", "01"),
		makeQuote("gunhild", "01"),
	}

	stats := NewCiconiaStats(quotes)
	result := stats.Compute(store.AllEpisodes).(*dto.CiconiaStatsResult)

	if len(result.TopSpeakers) != 3 {
		t.Fatalf("expected 3 speakers, got %d", len(result.TopSpeakers))
	}
	if result.TopSpeakers[0].CharacterID != "miyao" {
		t.Errorf("top speaker should be miyao, got %q", result.TopSpeakers[0].CharacterID)
	}
	if result.TopSpeakers[0].Count != 3 {
		t.Errorf("top speaker count should be 3, got %d", result.TopSpeakers[0].Count)
	}
	if result.TopSpeakers[1].CharacterID != "jayden" {
		t.Errorf("second speaker should be jayden, got %q", result.TopSpeakers[1].CharacterID)
	}
}

func TestCiconiaStats_NarratorExcluded(t *testing.T) {
	quotes := []dto.CiconiaQuote{
		makeQuote("narrator", "01"),
		makeQuote("narrator", "01"),
		makeQuote("miyao", "01"),
	}

	stats := NewCiconiaStats(quotes)
	result := stats.Compute(store.AllEpisodes).(*dto.CiconiaStatsResult)

	if len(result.TopSpeakers) != 1 {
		t.Fatalf("expected 1 speaker (narrator excluded), got %d", len(result.TopSpeakers))
	}
	if result.TopSpeakers[0].CharacterID != "miyao" {
		t.Errorf("expected miyao, got %q", result.TopSpeakers[0].CharacterID)
	}
}

func TestCiconiaStats_Interactions(t *testing.T) {
	quotes := []dto.CiconiaQuote{
		makeQuote("miyao", "01"),
		makeQuote("jayden", "01"),
		makeQuote("miyao", "01"),
		makeQuote("jayden", "01"),
	}

	stats := NewCiconiaStats(quotes)
	result := stats.Compute(store.AllEpisodes).(*dto.CiconiaStatsResult)

	if len(result.Interactions) == 0 {
		t.Fatal("expected interactions between miyao and jayden")
	}
	if result.Interactions[0].Count < 2 {
		t.Errorf("expected at least 2 interactions, got %d", result.Interactions[0].Count)
	}
}

func TestCiconiaStats_InteractionsBreakAcrossChapters(t *testing.T) {
	quotes := []dto.CiconiaQuote{
		makeQuote("miyao", "01"),
		makeQuote("jayden", "02"),
	}

	stats := NewCiconiaStats(quotes)
	result := stats.Compute(store.AllEpisodes).(*dto.CiconiaStatsResult)

	if len(result.Interactions) != 0 {
		t.Errorf("expected no interactions across different chapters, got %d", len(result.Interactions))
	}
}

func TestCiconiaStats_LinesPerChapter(t *testing.T) {
	quotes := []dto.CiconiaQuote{
		makeQuote("miyao", "01"),
		makeQuote("miyao", "01"),
		makeQuote("jayden", "01"),
		makeQuote("miyao", "02"),
		makeQuote("gunhild", "02"),
	}

	stats := NewCiconiaStats(quotes)
	result := stats.Compute(store.AllEpisodes).(*dto.CiconiaStatsResult)

	if result.LinesPerChapter == nil {
		t.Fatal("expected LinesPerChapter")
	}
	ch01 := result.LinesPerChapter["01"]
	if ch01 == nil {
		t.Fatal("expected chapter 01 in LinesPerChapter")
	}
	if ch01["miyao"] != 2 {
		t.Errorf("chapter 01 miyao: got %d, want 2", ch01["miyao"])
	}
	if ch01["jayden"] != 1 {
		t.Errorf("chapter 01 jayden: got %d, want 1", ch01["jayden"])
	}

	ch02 := result.LinesPerChapter["02"]
	if ch02 == nil {
		t.Fatal("expected chapter 02 in LinesPerChapter")
	}
	if ch02["miyao"] != 1 {
		t.Errorf("chapter 02 miyao: got %d, want 1", ch02["miyao"])
	}
}

func TestCiconiaStats_CharacterNames(t *testing.T) {
	quotes := []dto.CiconiaQuote{
		makeQuote("miyao", "01"),
	}

	stats := NewCiconiaStats(quotes)
	result := stats.Compute(store.AllEpisodes).(*dto.CiconiaStatsResult)

	if result.CharacterNames == nil {
		t.Fatal("expected CharacterNames")
	}
	if _, ok := result.CharacterNames["miyao"]; !ok {
		t.Error("expected miyao in CharacterNames")
	}
}

func TestCiconiaStats_EmptyQuotes(t *testing.T) {
	stats := NewCiconiaStats(nil)
	result := stats.Compute(store.AllEpisodes).(*dto.CiconiaStatsResult)

	if result == nil {
		t.Fatal("expected non-nil result for empty quotes")
	}
	if len(result.TopSpeakers) != 0 {
		t.Errorf("expected 0 speakers, got %d", len(result.TopSpeakers))
	}
}

func TestCiconiaStats_InteractionCountsMap(t *testing.T) {
	quotes := []dto.CiconiaQuote{
		makeQuote("miyao", "01"),
		makeQuote("jayden", "01"),
		makeQuote("miyao", "01"),
	}

	stats := NewCiconiaStats(quotes)
	result := stats.Compute(store.AllEpisodes).(*dto.CiconiaStatsResult)

	if len(result.InteractionCounts) == 0 {
		t.Fatal("expected InteractionCounts to be populated")
	}
}
