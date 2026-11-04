package memory

import (
	"slices"
	"testing"

	"umineko_quote/internal/quote/language"
)

func buildTestIndexer() *indexer {
	fields := map[language.Language]quoteFields{
		language.English: {
			Texts:        []string{"Hello World", "Beatrice speaks", "Narrator text here", "Battler again", "Episode three line"},
			CharacterIDs: []string{"10", "27", "narrator", "10", "27"},
			Episodes:     []int{1, 1, 2, 2, 3},
			AudioIDs:     []string{"", "", "", "", ""},
		},
	}
	return newIndexer(fields)
}

func TestIndexer_LowerTexts(t *testing.T) {
	idx := buildTestIndexer()

	texts := idx.LowerTexts(language.English)
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
	idx := buildTestIndexer()

	texts := idx.LowerTexts("fr")
	if texts != nil {
		t.Errorf("LowerTexts for unknown lang: got %v, want nil", texts)
	}
}

func TestIndexer_NonNarratorIndices(t *testing.T) {
	idx := buildTestIndexer()

	indices := idx.NonNarratorIndices(language.English)
	if len(indices) != 4 {
		t.Fatalf("NonNarratorIndices: got %d entries, want 4", len(indices))
	}
	if slices.Contains(indices, 2) {
		t.Errorf("NonNarratorIndices should not contain narrator index 2")
	}
}

func TestIndexer_NonNarratorIndices_UnknownLang(t *testing.T) {
	idx := buildTestIndexer()

	result := idx.NonNarratorIndices("fr")
	if result != nil {
		t.Errorf("NonNarratorIndices for unknown lang: got %v, want nil", result)
	}
}

func TestIndexer_FilteredIndices_CharacterOnly(t *testing.T) {
	idx := buildTestIndexer()

	indices := idx.FilteredIndices(language.English, "10", 0)
	if !slices.Equal(indices, []int{0, 3}) {
		t.Fatalf("FilteredIndices (char only): got %v, want [0 3]", indices)
	}
}

func TestIndexer_FilteredIndices_EpisodeOnly(t *testing.T) {
	idx := buildTestIndexer()

	indices := idx.FilteredIndices(language.English, "", 1)
	if len(indices) != 2 {
		t.Fatalf("FilteredIndices (ep only): got %d, want 2", len(indices))
	}
}

func TestIndexer_FilteredIndices_CharacterAndEpisode(t *testing.T) {
	idx := buildTestIndexer()

	indices := idx.FilteredIndices(language.English, "10", 1)
	if !slices.Equal(indices, []int{0}) {
		t.Fatalf("FilteredIndices (char+ep): got %v, want [0]", indices)
	}
}

func TestIndexer_FilteredIndices_CharacterAndEpisode_NoMatch(t *testing.T) {
	idx := buildTestIndexer()

	indices := idx.FilteredIndices(language.English, "10", 3)
	if len(indices) != 0 {
		t.Errorf("FilteredIndices (no match): got %d, want 0", len(indices))
	}
}

func TestIndexer_FilteredIndices_Neither(t *testing.T) {
	idx := buildTestIndexer()

	result := idx.FilteredIndices(language.English, "", 0)
	if result != nil {
		t.Errorf("FilteredIndices (neither): got %v, want nil", result)
	}
}

func TestIndexer_FilteredIndices_UnknownLang(t *testing.T) {
	idx := buildTestIndexer()

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

func TestIndexer_QuoteIndex_SubtitleIDs(t *testing.T) {
	fields := map[language.Language]quoteFields{
		language.English: {
			Texts:        []string{"Normal quote", "Welcome back, sir.", "It took me a while.", "Goodbye."},
			CharacterIDs: []string{"10", "00", "10", "00"},
			Episodes:     []int{0, 0, 0, 0},
			AudioIDs:     []string{"10100001", "end_all00_s0", "end_all00_s1", "end_all00_s2"},
		},
	}
	idx := newIndexer(fields)

	i, ok := idx.QuoteIndex(language.English, "end_all00_s0")
	if !ok {
		t.Fatal("QuoteIndex: expected to find end_all00_s0")
	}
	if i != 1 {
		t.Errorf("QuoteIndex end_all00_s0: got %d, want 1", i)
	}

	i, ok = idx.QuoteIndex(language.English, "end_all00_s2")
	if !ok {
		t.Fatal("QuoteIndex: expected to find end_all00_s2")
	}
	if i != 3 {
		t.Errorf("QuoteIndex end_all00_s2: got %d, want 3", i)
	}
}

func TestIndexer_SubtitleQuotes_EpisodeIndex(t *testing.T) {
	fields := map[language.Language]quoteFields{
		language.English: {
			Texts:        []string{"Episode 1 line", "Subtitle line 1", "Subtitle line 2", "Episode 3 line"},
			CharacterIDs: []string{"10", "00", "10", "27"},
			Episodes:     []int{1, 8, 8, 3},
			AudioIDs:     []string{"", "end_all00_s0", "end_all00_s1", ""},
		},
	}
	idx := newIndexer(fields)

	ep8 := idx.FilteredIndices(language.English, "", 8)
	if !slices.Equal(ep8, []int{1, 2}) {
		t.Errorf("FilteredIndices episode 8: got %v, want [1 2]", ep8)
	}
}

func TestIndexer_SubtitleQuotes_NonNarrator(t *testing.T) {
	fields := map[language.Language]quoteFields{
		language.English: {
			Texts:        []string{"Narrator line", "Sub line", "Battler sub line"},
			CharacterIDs: []string{"narrator", "00", "10"},
			Episodes:     []int{1, 8, 8},
			AudioIDs:     []string{"", "end_all00_s0", "end_all00_s1"},
		},
	}
	idx := newIndexer(fields)

	nonNarr := idx.NonNarratorIndices(language.English)
	if !slices.Equal(nonNarr, []int{1, 2}) {
		t.Errorf("NonNarratorIndices: got %v, want [1 2]", nonNarr)
	}
}

func TestIndexer_QuoteIndex_Found(t *testing.T) {
	fields := map[language.Language]quoteFields{
		language.English: {
			Texts:        []string{"First", "Second", "Third"},
			CharacterIDs: []string{"10", "27", "10"},
			Episodes:     []int{0, 0, 0},
			AudioIDs:     []string{"10100001", "12700001", "10100002"},
		},
	}
	idx := newIndexer(fields)

	i, ok := idx.QuoteIndex(language.English, "12700001")
	if !ok {
		t.Fatal("QuoteIndex: expected to find audio ID")
	}
	if i != 1 {
		t.Errorf("QuoteIndex: got index %d, want 1", i)
	}
}

func TestIndexer_QuoteIndex_NotFound(t *testing.T) {
	idx := buildTestIndexer()

	_, ok := idx.QuoteIndex(language.English, "99999999")
	if ok {
		t.Error("QuoteIndex: expected not found for unknown audio ID")
	}
}

func TestIndexer_QuoteIndex_CompositeIDs(t *testing.T) {
	fields := map[language.Language]quoteFields{
		language.English: {
			Texts:        []string{"Line one", "Line two"},
			CharacterIDs: []string{"10", "27"},
			Episodes:     []int{0, 0},
			AudioIDs:     []string{"10100001, 10100002", "12700001"},
		},
	}
	idx := newIndexer(fields)

	i1, ok1 := idx.QuoteIndex(language.English, "10100001")
	if !ok1 {
		t.Fatal("QuoteIndex: expected to find first sub-ID")
	}
	if i1 != 0 {
		t.Errorf("QuoteIndex first sub-ID: got %d, want 0", i1)
	}

	i2, ok2 := idx.QuoteIndex(language.English, "10100002")
	if !ok2 {
		t.Fatal("QuoteIndex: expected to find second sub-ID")
	}
	if i2 != 0 {
		t.Errorf("QuoteIndex second sub-ID: got %d, want 0", i2)
	}
}

func TestIndexer_QuoteIndex_UnknownLang(t *testing.T) {
	idx := buildTestIndexer()

	_, ok := idx.QuoteIndex("fr", "10100001")
	if ok {
		t.Error("QuoteIndex: expected not found for unknown lang")
	}
}

func TestIndexer_MultipleLangs(t *testing.T) {
	fields := map[language.Language]quoteFields{
		language.English: {
			Texts:        []string{"English text"},
			CharacterIDs: []string{"10"},
			Episodes:     []int{1},
			AudioIDs:     []string{""},
		},
		language.Japanese: {
			Texts:        []string{"日本語テキスト", "別の行"},
			CharacterIDs: []string{"10", "27"},
			Episodes:     []int{1, 2},
			AudioIDs:     []string{"", ""},
		},
	}
	idx := newIndexer(fields)

	if got := len(idx.LowerTexts(language.English)); got != 1 {
		t.Errorf("EN LowerTexts: got %d, want 1", got)
	}
	if got := len(idx.LowerTexts(language.Japanese)); got != 2 {
		t.Errorf("JA LowerTexts: got %d, want 2", got)
	}
	if got := len(idx.FilteredIndices(language.English, "10", 0)); got != 1 {
		t.Errorf("EN FilteredIndices for Battler: got %d, want 1", got)
	}
	if got := len(idx.FilteredIndices(language.Japanese, "27", 0)); got != 1 {
		t.Errorf("JA FilteredIndices for Beatrice: got %d, want 1", got)
	}
}
