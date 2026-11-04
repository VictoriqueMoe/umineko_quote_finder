package ciconia

import (
	"testing"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/quotes"
)

func TestCiconiaStats_TopSpeakers(t *testing.T) {
	tallies := params.NewSpeakerTallies()
	tallies.AddLines("miyao", "01", 3)
	tallies.AddLines("jayden", "01", 2)
	tallies.AddLines("gunhild", "01", 1)

	result := NewCiconiaStats(tallies).Compute(quotes.AllEpisodes).(*dto.CiconiaStatsResult)

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

func TestCiconiaStats_TopSpeakersSumAcrossChapters(t *testing.T) {
	tallies := params.NewSpeakerTallies()
	tallies.AddLines("miyao", "01", 1)
	tallies.AddLines("miyao", "02", 2)
	tallies.AddLines("jayden", "01", 2)

	result := NewCiconiaStats(tallies).Compute(quotes.AllEpisodes).(*dto.CiconiaStatsResult)

	if result.TopSpeakers[0].CharacterID != "miyao" || result.TopSpeakers[0].Count != 3 {
		t.Errorf("top speaker: got %+v, want miyao with 3", result.TopSpeakers[0])
	}
}

func TestCiconiaStats_Interactions(t *testing.T) {
	tallies := params.NewSpeakerTallies()
	tallies.AddLines("miyao", "01", 2)
	tallies.AddLines("jayden", "01", 2)
	tallies.AddInteractions("01", "jayden", "miyao", 2)
	tallies.AddInteractions("02", "miyao", "jayden", 1)

	result := NewCiconiaStats(tallies).Compute(quotes.AllEpisodes).(*dto.CiconiaStatsResult)

	if len(result.Interactions) != 1 {
		t.Fatalf("expected 1 interaction pair, got %d", len(result.Interactions))
	}
	if result.Interactions[0].Count != 3 {
		t.Errorf("expected interactions summed across chapters to be 3, got %d", result.Interactions[0].Count)
	}
}

func TestCiconiaStats_LinesPerChapter(t *testing.T) {
	tallies := params.NewSpeakerTallies()
	tallies.AddLines("miyao", "01", 2)
	tallies.AddLines("jayden", "01", 1)
	tallies.AddLines("miyao", "02", 1)
	tallies.AddLines("gunhild", "02", 1)

	result := NewCiconiaStats(tallies).Compute(quotes.AllEpisodes).(*dto.CiconiaStatsResult)

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
	tallies := params.NewSpeakerTallies()
	tallies.AddLines("miyao", "01", 1)

	result := NewCiconiaStats(tallies).Compute(quotes.AllEpisodes).(*dto.CiconiaStatsResult)

	if result.CharacterNames == nil {
		t.Fatal("expected CharacterNames")
	}
	if _, ok := result.CharacterNames["miyao"]; !ok {
		t.Error("expected miyao in CharacterNames")
	}
}

func TestCiconiaStats_EmptyTallies(t *testing.T) {
	result := NewCiconiaStats(params.NewSpeakerTallies()).Compute(quotes.AllEpisodes).(*dto.CiconiaStatsResult)

	if result == nil {
		t.Fatal("expected non-nil result for empty tallies")
	}
	if len(result.TopSpeakers) != 0 {
		t.Errorf("expected 0 speakers, got %d", len(result.TopSpeakers))
	}
}

func TestCiconiaStats_InteractionCountsMap(t *testing.T) {
	tallies := params.NewSpeakerTallies()
	tallies.AddLines("miyao", "01", 2)
	tallies.AddLines("jayden", "01", 1)
	tallies.AddInteractions("01", "miyao", "jayden", 2)

	result := NewCiconiaStats(tallies).Compute(quotes.AllEpisodes).(*dto.CiconiaStatsResult)

	if result.InteractionCounts["jayden|miyao"] != 2 {
		t.Errorf("expected InteractionCounts[jayden|miyao] = 2, got %d", result.InteractionCounts["jayden|miyao"])
	}
}
