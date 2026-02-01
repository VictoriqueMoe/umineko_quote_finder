package quote

import "testing"

func buildTestIndexer() (Indexer, map[string][]ParsedQuote) {
	quotes := map[string][]ParsedQuote{
		"en": {
			{Text: "Hello World", CharacterID: "10", Episode: 1},
			{Text: "Beatrice speaks", CharacterID: "27", Episode: 1},
			{Text: "Narrator text here", CharacterID: "narrator", Episode: 2},
			{Text: "Battler again", CharacterID: "10", Episode: 2},
			{Text: "Episode three line", CharacterID: "27", Episode: 3},
		},
	}
	return NewIndexer(quotes, ""), quotes
}

func TestIndexer_LowerTexts(t *testing.T) {
	idx, _ := buildTestIndexer()

	texts := idx.LowerTexts("en")
	if len(texts) != 5 {
		t.Fatalf("LowerTexts length: got %d, want 5", len(texts))
	}
	if texts[0] != "hello world" {
		t.Errorf("LowerTexts[0]: got %q, want %q", texts[0], "hello world")
	}
	if texts[1] != "beatrice speaks" {
		t.Errorf("LowerTexts[1]: got %q, want %q", texts[1], "beatrice speaks")
	}
}

func TestIndexer_LowerTexts_UnknownLang(t *testing.T) {
	idx, _ := buildTestIndexer()

	texts := idx.LowerTexts("fr")
	if texts != nil {
		t.Errorf("LowerTexts for unknown lang: got %v, want nil", texts)
	}
}

func TestIndexer_CharacterIndices(t *testing.T) {
	idx, _ := buildTestIndexer()

	battlerIdx := idx.CharacterIndices("en", "10")
	if len(battlerIdx) != 2 {
		t.Fatalf("CharacterIndices for Battler: got %d entries, want 2", len(battlerIdx))
	}
	if battlerIdx[0] != 0 || battlerIdx[1] != 3 {
		t.Errorf("CharacterIndices for Battler: got %v, want [0 3]", battlerIdx)
	}

	beatriceIdx := idx.CharacterIndices("en", "27")
	if len(beatriceIdx) != 2 {
		t.Fatalf("CharacterIndices for Beatrice: got %d entries, want 2", len(beatriceIdx))
	}

	narratorIdx := idx.CharacterIndices("en", "narrator")
	if len(narratorIdx) != 1 {
		t.Fatalf("CharacterIndices for narrator: got %d entries, want 1", len(narratorIdx))
	}
	if narratorIdx[0] != 2 {
		t.Errorf("CharacterIndices for narrator: got %v, want [2]", narratorIdx)
	}
}

func TestIndexer_CharacterIndices_UnknownLang(t *testing.T) {
	idx, _ := buildTestIndexer()

	result := idx.CharacterIndices("fr", "10")
	if result != nil {
		t.Errorf("CharacterIndices for unknown lang: got %v, want nil", result)
	}
}

func TestIndexer_CharacterIndices_UnknownCharacter(t *testing.T) {
	idx, _ := buildTestIndexer()

	result := idx.CharacterIndices("en", "99")
	if result != nil {
		t.Errorf("CharacterIndices for unknown character: got %v, want nil", result)
	}
}

func TestIndexer_NonNarratorIndices(t *testing.T) {
	idx, _ := buildTestIndexer()

	indices := idx.NonNarratorIndices("en")
	if len(indices) != 4 {
		t.Fatalf("NonNarratorIndices: got %d entries, want 4", len(indices))
	}
	for i := 0; i < len(indices); i++ {
		if indices[i] == 2 {
			t.Errorf("NonNarratorIndices should not contain narrator index 2")
		}
	}
}

func TestIndexer_NonNarratorIndices_UnknownLang(t *testing.T) {
	idx, _ := buildTestIndexer()

	result := idx.NonNarratorIndices("fr")
	if result != nil {
		t.Errorf("NonNarratorIndices for unknown lang: got %v, want nil", result)
	}
}

func TestIndexer_FilteredIndices_CharacterOnly(t *testing.T) {
	idx, _ := buildTestIndexer()

	indices := idx.FilteredIndices("en", "10", 0)
	if len(indices) != 2 {
		t.Fatalf("FilteredIndices (char only): got %d, want 2", len(indices))
	}
}

func TestIndexer_FilteredIndices_EpisodeOnly(t *testing.T) {
	idx, _ := buildTestIndexer()

	indices := idx.FilteredIndices("en", "", 1)
	if len(indices) != 2 {
		t.Fatalf("FilteredIndices (ep only): got %d, want 2", len(indices))
	}
}

func TestIndexer_FilteredIndices_CharacterAndEpisode(t *testing.T) {
	idx, _ := buildTestIndexer()

	indices := idx.FilteredIndices("en", "10", 1)
	if len(indices) != 1 {
		t.Fatalf("FilteredIndices (char+ep): got %d, want 1", len(indices))
	}
	if indices[0] != 0 {
		t.Errorf("FilteredIndices (char+ep): got index %d, want 0", indices[0])
	}
}

func TestIndexer_FilteredIndices_CharacterAndEpisode_NoMatch(t *testing.T) {
	idx, _ := buildTestIndexer()

	indices := idx.FilteredIndices("en", "10", 3)
	if len(indices) != 0 {
		t.Errorf("FilteredIndices (no match): got %d, want 0", len(indices))
	}
}

func TestIndexer_FilteredIndices_Neither(t *testing.T) {
	idx, _ := buildTestIndexer()

	result := idx.FilteredIndices("en", "", 0)
	if result != nil {
		t.Errorf("FilteredIndices (neither): got %v, want nil", result)
	}
}

func TestIndexer_FilteredIndices_UnknownLang(t *testing.T) {
	idx, _ := buildTestIndexer()

	result := idx.FilteredIndices("fr", "10", 0)
	if len(result) != 0 {
		t.Errorf("FilteredIndices unknown lang (char): got %v, want empty", result)
	}

	result = idx.FilteredIndices("fr", "", 1)
	if len(result) != 0 {
		t.Errorf("FilteredIndices unknown lang (ep): got %v, want empty", result)
	}

	result = idx.FilteredIndices("fr", "10", 1)
	if len(result) != 0 {
		t.Errorf("FilteredIndices unknown lang (both): got %v, want empty", result)
	}
}

func TestIndexer_AudioFilePath_EmptyDir(t *testing.T) {
	idx, _ := buildTestIndexer()

	path := idx.AudioFilePath("10", "10100001")
	if path != "" {
		t.Errorf("AudioFilePath with empty dir: got %q, want empty", path)
	}
}

func TestIndexer_AudioFilePath_NonexistentFile(t *testing.T) {
	quotes := map[string][]ParsedQuote{
		"en": {{Text: "test", CharacterID: "10", Episode: 1}},
	}
	idx := NewIndexer(quotes, "/nonexistent/audio/dir")

	path := idx.AudioFilePath("10", "10100001")
	if path != "" {
		t.Errorf("AudioFilePath with nonexistent file: got %q, want empty", path)
	}
}

func TestIndexer_MultipleLangs(t *testing.T) {
	quotes := map[string][]ParsedQuote{
		"en": {
			{Text: "English text", CharacterID: "10", Episode: 1},
		},
		"ja": {
			{Text: "日本語テキスト", CharacterID: "10", Episode: 1},
			{Text: "別の行", CharacterID: "27", Episode: 2},
		},
	}
	idx := NewIndexer(quotes, "")

	enTexts := idx.LowerTexts("en")
	if len(enTexts) != 1 {
		t.Errorf("EN LowerTexts: got %d, want 1", len(enTexts))
	}

	jaTexts := idx.LowerTexts("ja")
	if len(jaTexts) != 2 {
		t.Errorf("JA LowerTexts: got %d, want 2", len(jaTexts))
	}

	enBattler := idx.CharacterIndices("en", "10")
	if len(enBattler) != 1 {
		t.Errorf("EN CharacterIndices for Battler: got %d, want 1", len(enBattler))
	}

	jaBeatrice := idx.CharacterIndices("ja", "27")
	if len(jaBeatrice) != 1 {
		t.Errorf("JA CharacterIndices for Beatrice: got %d, want 1", len(jaBeatrice))
	}
}
