package umineko

import (
	"testing"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/quotes"

	"github.com/VictoriqueMoe/umineko_script_parser/umineko/character"
)

func buildTestTallies() params.SpeakerTallies {
	t := params.NewSpeakerTallies()

	t.AddLines("10", "1", 2)
	t.AddLines("10", "2", 1)
	t.AddLines("10", "3", 2)
	t.AddLines("27", "1", 1)
	t.AddLines("27", "2", 1)
	t.AddLines("27", "3", 2)

	t.AddInteractions("1", "10", "27", 1)
	t.AddInteractions("2", "27", "10", 1)
	t.AddInteractions("3", "10", "27", 2)

	return t
}

func buildTestTruth() map[int]params.TruthCounts {
	return map[int]params.TruthCounts{
		2: {Red: 1, Blue: 1},
		3: {Red: 1, Blue: 1},
	}
}

func newTestStats(t *testing.T, tallies params.SpeakerTallies, truth map[int]params.TruthCounts) quotes.Stats {
	t.Helper()

	s, err := NewUminekoStats(tallies, truth)
	if err != nil {
		t.Fatalf("NewUminekoStats: %v", err)
	}

	return s
}

func TestNewStats_Compute_AllEpisodes(t *testing.T) {
	s := newTestStats(t, buildTestTallies(), buildTestTruth())

	result := s.Compute(quotes.AllEpisodes)
	sr, ok := result.(*dto.StatsResult)
	if !ok {
		t.Fatalf("Compute(AllEpisodes) returned unexpected type %T", result)
	}

	if len(sr.TopSpeakers) == 0 {
		t.Fatal("TopSpeakers should not be empty")
	}

	if sr.TopSpeakers[0].CharacterID != "10" {
		t.Errorf("TopSpeakers[0] should be Battler (10), got %q", sr.TopSpeakers[0].CharacterID)
	}
	if sr.TopSpeakers[0].Name != character.CharacterNames[character.Battler] {
		t.Errorf("TopSpeakers[0] name: got %q, want %q", sr.TopSpeakers[0].Name, character.CharacterNames[character.Battler])
	}
	if sr.TopSpeakers[0].Count != 5 {
		t.Errorf("TopSpeakers[0] count: got %d, want 5", sr.TopSpeakers[0].Count)
	}
}

func TestStats_TopSpeakers_Ranking(t *testing.T) {
	s := newTestStats(t, buildTestTallies(), buildTestTruth())
	result := s.Compute(quotes.AllEpisodes).(*dto.StatsResult)

	if len(result.TopSpeakers) < 2 {
		t.Fatalf("expected at least 2 top speakers, got %d", len(result.TopSpeakers))
	}

	if result.TopSpeakers[0].Count < result.TopSpeakers[1].Count {
		t.Errorf("TopSpeakers not sorted by count: %d < %d",
			result.TopSpeakers[0].Count, result.TopSpeakers[1].Count)
	}
}

func TestStats_TruthPerEpisode(t *testing.T) {
	s := newTestStats(t, buildTestTallies(), buildTestTruth())
	result := s.Compute(quotes.AllEpisodes).(*dto.StatsResult)

	if len(result.TruthPerEpisode) != 8 {
		t.Fatalf("TruthPerEpisode length: got %d, want 8", len(result.TruthPerEpisode))
	}

	if result.TruthPerEpisode[1].Red != 1 {
		t.Errorf("Episode 2 red truth: got %d, want 1", result.TruthPerEpisode[1].Red)
	}
	if result.TruthPerEpisode[1].Blue != 1 {
		t.Errorf("Episode 2 blue truth: got %d, want 1", result.TruthPerEpisode[1].Blue)
	}

	if result.TruthPerEpisode[2].Red != 1 {
		t.Errorf("Episode 3 red truth: got %d, want 1", result.TruthPerEpisode[2].Red)
	}
	if result.TruthPerEpisode[2].Blue != 1 {
		t.Errorf("Episode 3 blue truth: got %d, want 1", result.TruthPerEpisode[2].Blue)
	}

	if result.TruthPerEpisode[0].Red != 0 {
		t.Errorf("Episode 1 red truth: got %d, want 0", result.TruthPerEpisode[0].Red)
	}
}

func TestStats_LinesPerEpisode(t *testing.T) {
	s := newTestStats(t, buildTestTallies(), buildTestTruth())
	result := s.Compute(quotes.AllEpisodes).(*dto.StatsResult)

	if len(result.LinesPerEpisode) != 8 {
		t.Fatalf("LinesPerEpisode length: got %d, want 8", len(result.LinesPerEpisode))
	}

	ep1 := result.LinesPerEpisode[0]
	if ep1.Episode != 1 {
		t.Errorf("LinesPerEpisode[0].Episode: got %d, want 1", ep1.Episode)
	}
	if ep1.EpisodeName != "Legend" {
		t.Errorf("LinesPerEpisode[0].EpisodeName: got %q, want %q", ep1.EpisodeName, "Legend")
	}
	if ep1.Characters["10"] != 2 {
		t.Errorf("Episode 1 Battler lines: got %d, want 2", ep1.Characters["10"])
	}
	if ep1.Characters["27"] != 1 {
		t.Errorf("Episode 1 Beatrice lines: got %d, want 1", ep1.Characters["27"])
	}
}

func TestStats_Interactions(t *testing.T) {
	s := newTestStats(t, buildTestTallies(), buildTestTruth())
	result := s.Compute(quotes.AllEpisodes).(*dto.StatsResult)

	if len(result.Interactions) == 0 {
		t.Fatal("Interactions should not be empty")
	}

	found := false
	for _, pair := range result.Interactions {
		if (pair.CharA == "10" && pair.CharB == "27") || (pair.CharA == "27" && pair.CharB == "10") {
			found = true
			if pair.Count != 4 {
				t.Errorf("Battler-Beatrice interaction count: got %d, want 4", pair.Count)
			}
		}
	}
	if !found {
		t.Error("expected to find Battler-Beatrice interaction pair")
	}
}

func TestStats_InteractionCountsMap(t *testing.T) {
	s := newTestStats(t, buildTestTallies(), buildTestTruth())
	result := s.Compute(quotes.AllEpisodes).(*dto.StatsResult)

	if len(result.InteractionCounts) == 0 {
		t.Fatal("InteractionCounts should not be empty")
	}

	count, ok := result.InteractionCounts["10|27"]
	if !ok {
		t.Fatal(`expected InteractionCounts to include key "10|27"`)
	}
	if count == 0 {
		t.Error(`expected InteractionCounts["10|27"] to be > 0`)
	}
}

func TestStats_CharacterPresence(t *testing.T) {
	s := newTestStats(t, buildTestTallies(), buildTestTruth())
	result := s.Compute(quotes.AllEpisodes).(*dto.StatsResult)

	if len(result.CharacterPresence) == 0 {
		t.Fatal("CharacterPresence should not be empty")
	}

	if len(result.CharacterPresence[0].Episodes) != 8 {
		t.Errorf("CharacterPresence episodes length: got %d, want 8", len(result.CharacterPresence[0].Episodes))
	}
}

func TestStats_CharacterNames(t *testing.T) {
	s := newTestStats(t, buildTestTallies(), buildTestTruth())
	result := s.Compute(quotes.AllEpisodes).(*dto.StatsResult)

	if result.CharacterNames["10"] != character.CharacterNames[character.Battler] {
		t.Errorf("CharacterNames[10]: got %q, want %q", result.CharacterNames["10"], character.CharacterNames[character.Battler])
	}
	if result.CharacterNames["27"] != character.CharacterNames[character.Beatrice] {
		t.Errorf("CharacterNames[27]: got %q, want %q", result.CharacterNames["27"], character.CharacterNames[character.Beatrice])
	}
}

func TestStats_EpisodeNames(t *testing.T) {
	s := newTestStats(t, buildTestTallies(), buildTestTruth())
	result := s.Compute(quotes.AllEpisodes).(*dto.StatsResult)

	expected := map[int]string{
		1: "Legend", 2: "Turn", 3: "Banquet", 4: "Alliance",
		5: "End", 6: "Dawn", 7: "Requiem", 8: "Twilight",
	}

	for ep, name := range expected {
		if result.EpisodeNames[ep] != name {
			t.Errorf("EpisodeNames[%d]: got %q, want %q", ep, result.EpisodeNames[ep], name)
		}
	}
}

func TestStats_ComputeSpecificEpisode(t *testing.T) {
	s := newTestStats(t, buildTestTallies(), buildTestTruth())

	result := s.Compute(1).(*dto.StatsResult)

	for _, speaker := range result.TopSpeakers {
		if speaker.CharacterID == "10" && speaker.Count != 2 {
			t.Errorf("Battler count in ep1: got %d, want 2", speaker.Count)
		}
	}

	if result.InteractionCounts["10|27"] != 1 {
		t.Errorf("ep1 interactions: got %d, want 1", result.InteractionCounts["10|27"])
	}

	if result.LinesPerEpisode != nil {
		t.Error("LinesPerEpisode should be nil for specific episode")
	}
	if result.TruthPerEpisode != nil {
		t.Error("TruthPerEpisode should be nil for specific episode")
	}
	if result.CharacterPresence != nil {
		t.Error("CharacterPresence should be nil for specific episode")
	}
}

func TestStats_ComputeUnknownEpisodeIsEmpty(t *testing.T) {
	s := newTestStats(t, buildTestTallies(), buildTestTruth())

	result := s.Compute(99).(*dto.StatsResult)

	if len(result.TopSpeakers) != 0 {
		t.Errorf("TopSpeakers should be empty for an unknown episode, got %d", len(result.TopSpeakers))
	}
	if result.InteractionCounts == nil {
		t.Error("InteractionCounts should be an empty map, not nil")
	}
}

func TestStats_ComputeCached(t *testing.T) {
	s := newTestStats(t, buildTestTallies(), buildTestTruth())

	result1 := s.Compute(quotes.AllEpisodes)
	result2 := s.Compute(quotes.AllEpisodes)

	if result1 != result2 {
		t.Error("Compute(AllEpisodes) should return cached result")
	}
}

func TestStats_EmptyTallies(t *testing.T) {
	s := newTestStats(t, params.NewSpeakerTallies(), map[int]params.TruthCounts{})
	result := s.Compute(quotes.AllEpisodes).(*dto.StatsResult)

	if len(result.TopSpeakers) != 0 {
		t.Errorf("TopSpeakers should be empty, got %d", len(result.TopSpeakers))
	}
	if len(result.Interactions) != 0 {
		t.Errorf("Interactions should be empty, got %d", len(result.Interactions))
	}
}

func TestStats_NonNumericSegmentFails(t *testing.T) {
	tallies := params.NewSpeakerTallies()
	tallies.AddLines("10", "tea", 1)

	if _, err := NewUminekoStats(tallies, map[int]params.TruthCounts{}); err == nil {
		t.Error("expected an error for a segment that is not an episode number")
	}
}
