package quote

import "testing"

func buildTestQuotes() []ParsedQuote {
	return []ParsedQuote{
		{Text: "Line 1", TextHtml: "Line 1", CharacterID: "10", Episode: 1},
		{Text: "Line 2", TextHtml: "Line 2", CharacterID: "10", Episode: 1},
		{Text: "Line 3", TextHtml: "Line 3", CharacterID: "27", Episode: 1},
		{Text: "Narrator line", TextHtml: "Narrator line", CharacterID: "narrator", Episode: 1},
		{Text: "Red truth", TextHtml: `<span class="red-truth">Red</span>`, CharacterID: "27", Episode: 2, HasRedTruth: true},
		{Text: "Blue truth", TextHtml: `<span class="blue-truth">Blue</span>`, CharacterID: "10", Episode: 2, HasBlueTruth: true},
		{Text: "Red truth ep3", TextHtml: `<span class="red-truth">Red</span>`, CharacterID: "10", Episode: 3, HasRedTruth: true},
		{Text: "Blue truth ep3", TextHtml: `<span class="blue-truth">Blue</span>`, CharacterID: "27", Episode: 3, HasBlueTruth: true},
		{Text: "Line ep3", TextHtml: "Line ep3", CharacterID: "27", Episode: 3},
		{Text: "Line ep3 b", TextHtml: "Line ep3 b", CharacterID: "10", Episode: 3},
	}
}

func TestNewStats_Compute_AllEpisodes(t *testing.T) {
	quotes := buildTestQuotes()
	s := NewStats(quotes)

	result := s.Compute(AllEpisodes)
	sr, ok := result.(*statsResult)
	if !ok {
		t.Fatalf("Compute(AllEpisodes) returned unexpected type %T", result)
	}

	if len(sr.TopSpeakers) == 0 {
		t.Fatal("TopSpeakers should not be empty")
	}

	if sr.TopSpeakers[0].CharacterID != "10" {
		t.Errorf("TopSpeakers[0] should be Battler (10), got %q", sr.TopSpeakers[0].CharacterID)
	}
	if sr.TopSpeakers[0].Name != CharacterNames["10"] {
		t.Errorf("TopSpeakers[0] name: got %q, want %q", sr.TopSpeakers[0].Name, CharacterNames["10"])
	}
	if sr.TopSpeakers[0].Count != 5 {
		t.Errorf("TopSpeakers[0] count: got %d, want 5", sr.TopSpeakers[0].Count)
	}
}

func TestStats_TopSpeakers_Ranking(t *testing.T) {
	quotes := buildTestQuotes()
	s := NewStats(quotes)
	result := s.Compute(AllEpisodes).(*statsResult)

	if len(result.TopSpeakers) < 2 {
		t.Fatalf("expected at least 2 top speakers, got %d", len(result.TopSpeakers))
	}

	if result.TopSpeakers[0].Count < result.TopSpeakers[1].Count {
		t.Errorf("TopSpeakers not sorted by count: %d < %d",
			result.TopSpeakers[0].Count, result.TopSpeakers[1].Count)
	}
}

func TestStats_TopSpeakers_ExcludesNarrator(t *testing.T) {
	quotes := buildTestQuotes()
	s := NewStats(quotes)
	result := s.Compute(AllEpisodes).(*statsResult)

	for i := 0; i < len(result.TopSpeakers); i++ {
		if result.TopSpeakers[i].CharacterID == "narrator" {
			t.Error("TopSpeakers should not include narrator")
		}
	}
}

func TestStats_TruthPerEpisode(t *testing.T) {
	quotes := buildTestQuotes()
	s := NewStats(quotes)
	result := s.Compute(AllEpisodes).(*statsResult)

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
	quotes := buildTestQuotes()
	s := NewStats(quotes)
	result := s.Compute(AllEpisodes).(*statsResult)

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
	quotes := buildTestQuotes()
	s := NewStats(quotes)
	result := s.Compute(AllEpisodes).(*statsResult)

	if len(result.Interactions) == 0 {
		t.Fatal("Interactions should not be empty")
	}

	found := false
	for i := 0; i < len(result.Interactions); i++ {
		pair := result.Interactions[i]
		if (pair.CharA == "10" && pair.CharB == "27") || (pair.CharA == "27" && pair.CharB == "10") {
			found = true
			if pair.Count == 0 {
				t.Error("Battler-Beatrice interaction count should be > 0")
			}
		}
	}
	if !found {
		t.Error("expected to find Battler-Beatrice interaction pair")
	}
}

func TestStats_CharacterPresence(t *testing.T) {
	quotes := buildTestQuotes()
	s := NewStats(quotes)
	result := s.Compute(AllEpisodes).(*statsResult)

	if len(result.CharacterPresence) == 0 {
		t.Fatal("CharacterPresence should not be empty")
	}

	if len(result.CharacterPresence[0].Episodes) != 8 {
		t.Errorf("CharacterPresence episodes length: got %d, want 8", len(result.CharacterPresence[0].Episodes))
	}
}

func TestStats_CharacterNames(t *testing.T) {
	quotes := buildTestQuotes()
	s := NewStats(quotes)
	result := s.Compute(AllEpisodes).(*statsResult)

	if result.CharacterNames["10"] != CharacterNames["10"] {
		t.Errorf("CharacterNames[10]: got %q, want %q", result.CharacterNames["10"], CharacterNames["10"])
	}
	if result.CharacterNames["27"] != CharacterNames["27"] {
		t.Errorf("CharacterNames[27]: got %q, want %q", result.CharacterNames["27"], CharacterNames["27"])
	}
}

func TestStats_EpisodeNames(t *testing.T) {
	quotes := buildTestQuotes()
	s := NewStats(quotes)
	result := s.Compute(AllEpisodes).(*statsResult)

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
	quotes := buildTestQuotes()
	s := NewStats(quotes)

	result := s.Compute(1).(*statsResult)

	for i := 0; i < len(result.TopSpeakers); i++ {
		speaker := result.TopSpeakers[i]
		if speaker.CharacterID == "10" && speaker.Count != 2 {
			t.Errorf("Battler count in ep1: got %d, want 2", speaker.Count)
		}
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

func TestStats_ComputeCached(t *testing.T) {
	quotes := buildTestQuotes()
	s := NewStats(quotes)

	result1 := s.Compute(AllEpisodes)
	result2 := s.Compute(AllEpisodes)

	if result1 != result2 {
		t.Error("Compute(AllEpisodes) should return cached result")
	}
}

func TestStats_EmptyQuotes(t *testing.T) {
	s := NewStats([]ParsedQuote{})
	result := s.Compute(AllEpisodes).(*statsResult)

	if len(result.TopSpeakers) != 0 {
		t.Errorf("TopSpeakers should be empty, got %d", len(result.TopSpeakers))
	}
	if len(result.Interactions) != 0 {
		t.Errorf("Interactions should be empty, got %d", len(result.Interactions))
	}
}

func TestStats_InteractionPairOrdering(t *testing.T) {
	quotes := []ParsedQuote{
		{Text: "Line 1", TextHtml: "Line 1", CharacterID: "27", Episode: 1},
		{Text: "Line 2", TextHtml: "Line 2", CharacterID: "10", Episode: 1},
	}
	s := NewStats(quotes)
	result := s.Compute(AllEpisodes).(*statsResult)

	if len(result.Interactions) != 1 {
		t.Fatalf("expected 1 interaction, got %d", len(result.Interactions))
	}
	if result.Interactions[0].CharA != "10" || result.Interactions[0].CharB != "27" {
		t.Errorf("interaction pair should be ordered: got %q-%q, want 10-27",
			result.Interactions[0].CharA, result.Interactions[0].CharB)
	}
}

func TestStats_NarratorBreaksInteraction(t *testing.T) {
	quotes := []ParsedQuote{
		{Text: "Line 1", TextHtml: "Line 1", CharacterID: "10", Episode: 1},
		{Text: "Narration", TextHtml: "Narration", CharacterID: "narrator", Episode: 1},
		{Text: "Line 2", TextHtml: "Line 2", CharacterID: "27", Episode: 1},
	}
	s := NewStats(quotes)
	result := s.Compute(AllEpisodes).(*statsResult)

	if len(result.Interactions) != 0 {
		t.Errorf("narrator should break interaction chain, got %d interactions", len(result.Interactions))
	}
}
