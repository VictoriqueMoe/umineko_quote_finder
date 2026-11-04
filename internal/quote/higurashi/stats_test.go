package higurashi

import (
	"testing"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/quotes"
)

func TestHigurashiStats_TopSpeakers(t *testing.T) {
	tallies := params.NewSpeakerTallies()
	tallies.AddLines("1", "onikakushi", 3)
	tallies.AddLines("2", "onikakushi", 2)
	tallies.AddLines("3", "onikakushi", 1)

	result := NewHigurashiStats(tallies).Compute(quotes.AllEpisodes).(*dto.HigurashiStatsResult)

	if len(result.TopSpeakers) != 3 {
		t.Fatalf("expected 3 speakers, got %d", len(result.TopSpeakers))
	}
	if result.TopSpeakers[0].CharacterID != "1" {
		t.Errorf("top speaker should be 1, got %q", result.TopSpeakers[0].CharacterID)
	}
	if result.TopSpeakers[0].Count != 3 {
		t.Errorf("top speaker count should be 3, got %d", result.TopSpeakers[0].Count)
	}
	if result.TopSpeakers[1].CharacterID != "2" {
		t.Errorf("second speaker should be 2, got %q", result.TopSpeakers[1].CharacterID)
	}
}

func TestHigurashiStats_Interactions(t *testing.T) {
	tallies := params.NewSpeakerTallies()
	tallies.AddLines("1", "onikakushi", 2)
	tallies.AddLines("2", "onikakushi", 2)
	tallies.AddInteractions("onikakushi", "1", "2", 3)

	result := NewHigurashiStats(tallies).Compute(quotes.AllEpisodes).(*dto.HigurashiStatsResult)

	if len(result.Interactions) == 0 {
		t.Fatal("expected interactions between chars 1 and 2")
	}
	if result.Interactions[0].Count != 3 {
		t.Errorf("expected 3 interactions, got %d", result.Interactions[0].Count)
	}
}

func TestHigurashiStats_LinesPerArc(t *testing.T) {
	tallies := params.NewSpeakerTallies()
	tallies.AddLines("1", "onikakushi", 2)
	tallies.AddLines("2", "onikakushi", 1)
	tallies.AddLines("1", "watanagashi", 1)
	tallies.AddLines("3", "watanagashi", 1)

	result := NewHigurashiStats(tallies).Compute(quotes.AllEpisodes).(*dto.HigurashiStatsResult)

	if result.LinesPerArc == nil {
		t.Fatal("expected LinesPerArc")
	}
	oni := result.LinesPerArc["onikakushi"]
	if oni == nil {
		t.Fatal("expected onikakushi in LinesPerArc")
	}
	if oni["1"] != 2 {
		t.Errorf("onikakushi char 1: got %d, want 2", oni["1"])
	}
	if oni["2"] != 1 {
		t.Errorf("onikakushi char 2: got %d, want 1", oni["2"])
	}

	wata := result.LinesPerArc["watanagashi"]
	if wata == nil {
		t.Fatal("expected watanagashi in LinesPerArc")
	}
	if wata["1"] != 1 {
		t.Errorf("watanagashi char 1: got %d, want 1", wata["1"])
	}
}

func TestHigurashiStats_CharacterNames(t *testing.T) {
	tallies := params.NewSpeakerTallies()
	tallies.AddLines("1", "onikakushi", 1)

	result := NewHigurashiStats(tallies).Compute(quotes.AllEpisodes).(*dto.HigurashiStatsResult)

	if result.CharacterNames == nil {
		t.Fatal("expected CharacterNames")
	}
	if _, ok := result.CharacterNames["1"]; !ok {
		t.Error("expected char 1 in CharacterNames")
	}
}

func TestHigurashiStats_EmptyTallies(t *testing.T) {
	result := NewHigurashiStats(params.NewSpeakerTallies()).Compute(quotes.AllEpisodes).(*dto.HigurashiStatsResult)

	if result == nil {
		t.Fatal("expected non-nil result for empty tallies")
	}
	if len(result.TopSpeakers) != 0 {
		t.Errorf("expected 0 speakers, got %d", len(result.TopSpeakers))
	}
}

func TestHigurashiStats_InteractionCountsMap(t *testing.T) {
	tallies := params.NewSpeakerTallies()
	tallies.AddLines("1", "onikakushi", 2)
	tallies.AddLines("2", "onikakushi", 1)
	tallies.AddInteractions("onikakushi", "2", "1", 2)

	result := NewHigurashiStats(tallies).Compute(quotes.AllEpisodes).(*dto.HigurashiStatsResult)

	if result.InteractionCounts["1|2"] != 2 {
		t.Errorf("expected InteractionCounts[1|2] = 2, got %d", result.InteractionCounts["1|2"])
	}
}
