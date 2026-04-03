package store

import (
	"os"
	"path/filepath"
	"testing"

	"umineko_quote/internal/quote/language"
)

func buildTestIndexer() Indexer {
	fields := map[language.Language]QuoteFields{
		language.English: {
			Texts:        []string{"Hello World", "Beatrice speaks", "Narrator text here", "Battler again", "Episode three line"},
			CharacterIDs: []string{"10", "27", "narrator", "10", "27"},
			Episodes:     []int{1, 1, 2, 2, 3},
			AudioIDs:     []string{"", "", "", "", ""},
		},
	}
	return NewIndexer(fields, "")
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

	texts := idx.LowerTexts(language.Language("fr"))
	if texts != nil {
		t.Errorf("LowerTexts for unknown lang: got %v, want nil", texts)
	}
}

func TestIndexer_CharacterIndices(t *testing.T) {
	idx := buildTestIndexer()

	battlerIdx := idx.CharacterIndices(language.English, "10")
	if len(battlerIdx) != 2 {
		t.Fatalf("CharacterIndices for Battler: got %d entries, want 2", len(battlerIdx))
	}
	if battlerIdx[0] != 0 || battlerIdx[1] != 3 {
		t.Errorf("CharacterIndices for Battler: got %v, want [0 3]", battlerIdx)
	}

	beatriceIdx := idx.CharacterIndices(language.English, "27")
	if len(beatriceIdx) != 2 {
		t.Fatalf("CharacterIndices for Beatrice: got %d entries, want 2", len(beatriceIdx))
	}

	narratorIdx := idx.CharacterIndices(language.English, "narrator")
	if len(narratorIdx) != 1 {
		t.Fatalf("CharacterIndices for narrator: got %d entries, want 1", len(narratorIdx))
	}
	if narratorIdx[0] != 2 {
		t.Errorf("CharacterIndices for narrator: got %v, want [2]", narratorIdx)
	}
}

func TestIndexer_CharacterIndices_UnknownLang(t *testing.T) {
	idx := buildTestIndexer()

	result := idx.CharacterIndices(language.Language("fr"), "10")
	if result != nil {
		t.Errorf("CharacterIndices for unknown lang: got %v, want nil", result)
	}
}

func TestIndexer_CharacterIndices_UnknownCharacter(t *testing.T) {
	idx := buildTestIndexer()

	result := idx.CharacterIndices(language.English, "99")
	if result != nil {
		t.Errorf("CharacterIndices for unknown character: got %v, want nil", result)
	}
}

func TestIndexer_NonNarratorIndices(t *testing.T) {
	idx := buildTestIndexer()

	indices := idx.NonNarratorIndices(language.English)
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
	idx := buildTestIndexer()

	result := idx.NonNarratorIndices(language.Language("fr"))
	if result != nil {
		t.Errorf("NonNarratorIndices for unknown lang: got %v, want nil", result)
	}
}

func TestIndexer_FilteredIndices_CharacterOnly(t *testing.T) {
	idx := buildTestIndexer()

	indices := idx.FilteredIndices(language.English, "10", 0)
	if len(indices) != 2 {
		t.Fatalf("FilteredIndices (char only): got %d, want 2", len(indices))
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
	if len(indices) != 1 {
		t.Fatalf("FilteredIndices (char+ep): got %d, want 1", len(indices))
	}
	if indices[0] != 0 {
		t.Errorf("FilteredIndices (char+ep): got index %d, want 0", indices[0])
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

	result := idx.FilteredIndices(language.Language("fr"), "10", 0)
	if len(result) != 0 {
		t.Errorf("FilteredIndices unknown lang (char): got %v, want empty", result)
	}

	result = idx.FilteredIndices(language.Language("fr"), "", 1)
	if len(result) != 0 {
		t.Errorf("FilteredIndices unknown lang (ep): got %v, want empty", result)
	}

	result = idx.FilteredIndices(language.Language("fr"), "10", 1)
	if len(result) != 0 {
		t.Errorf("FilteredIndices unknown lang (both): got %v, want empty", result)
	}
}

func TestIndexer_InteractionIndices(t *testing.T) {
	idx := buildTestIndexer()

	indices := idx.InteractionIndices(language.English, "10", "27")
	if len(indices) != 2 {
		t.Fatalf("InteractionIndices (10,27): got %d, want 2", len(indices))
	}
	if indices[0] != 0 || indices[1] != 1 {
		t.Errorf("InteractionIndices (10,27): got %v, want [0 1]", indices)
	}
}

func TestIndexer_InteractionIndices_SameCharacter(t *testing.T) {
	idx := buildTestIndexer()

	indices := idx.InteractionIndices(language.English, "10", "10")
	if len(indices) != 0 {
		t.Errorf("InteractionIndices same character: got %v, want empty", indices)
	}
}

func TestIndexer_InteractionIndices_UnknownLang(t *testing.T) {
	idx := buildTestIndexer()

	indices := idx.InteractionIndices(language.Language("fr"), "10", "27")
	if len(indices) != 0 {
		t.Errorf("InteractionIndices unknown lang: got %v, want empty", indices)
	}
}

func TestIndexer_AudioFilePath_EmptyDir(t *testing.T) {
	idx := buildTestIndexer()

	path := idx.AudioFilePath("10", "10100001")
	if path != "" {
		t.Errorf("AudioFilePath with empty dir: got %q, want empty", path)
	}
}

func TestIndexer_AudioFilePath_NonexistentFile(t *testing.T) {
	fields := map[language.Language]QuoteFields{
		language.English: {
			Texts:        []string{"test"},
			CharacterIDs: []string{"10"},
			Episodes:     []int{1},
			AudioIDs:     []string{""},
		},
	}
	idx := NewIndexer(fields, "/nonexistent/audio/dir")

	path := idx.AudioFilePath("10", "10100001")
	if path != "" {
		t.Errorf("AudioFilePath with nonexistent file: got %q, want empty", path)
	}
}

func TestIndexer_AudioFilePath_SubtitleSuffixStrip(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "00")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "end_all00.ogg"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	fields := map[language.Language]QuoteFields{
		language.English: {
			Texts:        []string{"test"},
			CharacterIDs: []string{"00"},
			Episodes:     []int{8},
			AudioIDs:     []string{""},
		},
	}
	idx := NewIndexer(fields, dir)

	path := idx.AudioFilePath("00", "end_all00_s0")
	if path == "" {
		t.Fatal("AudioFilePath should strip _s0 suffix and find end_all00.ogg")
	}
	expected := filepath.Join(dir, "00", "end_all00.ogg")
	if path != expected {
		t.Errorf("AudioFilePath: got %q, want %q", path, expected)
	}

	path = idx.AudioFilePath("00", "end_all00_s15")
	if path == "" {
		t.Fatal("AudioFilePath should strip _s15 suffix and find end_all00.ogg")
	}
	if path != expected {
		t.Errorf("AudioFilePath: got %q, want %q", path, expected)
	}
}

func TestIndexer_AudioFilePath_SubtitleSuffixNoBase(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "00")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	fields := map[language.Language]QuoteFields{
		language.English: {
			Texts:        []string{"test"},
			CharacterIDs: []string{"00"},
			Episodes:     []int{8},
			AudioIDs:     []string{""},
		},
	}
	idx := NewIndexer(fields, dir)

	path := idx.AudioFilePath("00", "end_all00_s5")
	if path != "" {
		t.Errorf("AudioFilePath should return empty when base file also missing: got %q", path)
	}
}

func TestIndexer_AudioFilePath_ExactMatchPreferred(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "00")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "end_all00_s0.ogg"), []byte("exact"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "end_all00.ogg"), []byte("base"), 0644); err != nil {
		t.Fatal(err)
	}

	fields := map[language.Language]QuoteFields{
		language.English: {
			Texts:        []string{"test"},
			CharacterIDs: []string{"00"},
			Episodes:     []int{8},
			AudioIDs:     []string{""},
		},
	}
	idx := NewIndexer(fields, dir)

	path := idx.AudioFilePath("00", "end_all00_s0")
	expected := filepath.Join(dir, "00", "end_all00_s0.ogg")
	if path != expected {
		t.Errorf("AudioFilePath should prefer exact match: got %q, want %q", path, expected)
	}
}

func TestIndexer_QuoteIndex_SubtitleIDs(t *testing.T) {
	fields := map[language.Language]QuoteFields{
		language.English: {
			Texts:        []string{"Normal quote", "Welcome back, sir.", "It took me a while.", "Goodbye."},
			CharacterIDs: []string{"10", "00", "10", "00"},
			Episodes:     []int{0, 0, 0, 0},
			AudioIDs:     []string{"10100001", "end_all00_s0", "end_all00_s1", "end_all00_s2"},
		},
	}
	idx := NewIndexer(fields, "")

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
	fields := map[language.Language]QuoteFields{
		language.English: {
			Texts:        []string{"Episode 1 line", "Subtitle line 1", "Subtitle line 2", "Episode 3 line"},
			CharacterIDs: []string{"10", "00", "10", "27"},
			Episodes:     []int{1, 8, 8, 3},
			AudioIDs:     []string{"", "end_all00_s0", "end_all00_s1", ""},
		},
	}
	idx := NewIndexer(fields, "")

	ep8 := idx.FilteredIndices(language.English, "", 8)
	if len(ep8) != 2 {
		t.Fatalf("FilteredIndices episode 8: got %d, want 2", len(ep8))
	}
	if ep8[0] != 1 || ep8[1] != 2 {
		t.Errorf("FilteredIndices episode 8: got %v, want [1 2]", ep8)
	}
}

func TestIndexer_SubtitleQuotes_NonNarrator(t *testing.T) {
	fields := map[language.Language]QuoteFields{
		language.English: {
			Texts:        []string{"Narrator line", "Sub line", "Battler sub line"},
			CharacterIDs: []string{"narrator", "00", "10"},
			Episodes:     []int{1, 8, 8},
			AudioIDs:     []string{"", "end_all00_s0", "end_all00_s1"},
		},
	}
	idx := NewIndexer(fields, "")

	nonNarr := idx.NonNarratorIndices(language.English)
	if len(nonNarr) != 2 {
		t.Fatalf("NonNarratorIndices: got %d, want 2", len(nonNarr))
	}
	if nonNarr[0] != 1 || nonNarr[1] != 2 {
		t.Errorf("NonNarratorIndices: got %v, want [1 2]", nonNarr)
	}
}

func TestIndexer_QuoteIndex_Found(t *testing.T) {
	fields := map[language.Language]QuoteFields{
		language.English: {
			Texts:        []string{"First", "Second", "Third"},
			CharacterIDs: []string{"10", "27", "10"},
			Episodes:     []int{0, 0, 0},
			AudioIDs:     []string{"10100001", "12700001", "10100002"},
		},
	}
	idx := NewIndexer(fields, "")

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
	fields := map[language.Language]QuoteFields{
		language.English: {
			Texts:        []string{"Line one", "Line two"},
			CharacterIDs: []string{"10", "27"},
			Episodes:     []int{0, 0},
			AudioIDs:     []string{"10100001, 10100002", "12700001"},
		},
	}
	idx := NewIndexer(fields, "")

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

	_, ok := idx.QuoteIndex(language.Language("fr"), "10100001")
	if ok {
		t.Error("QuoteIndex: expected not found for unknown lang")
	}
}

func TestIndexer_MultipleLangs(t *testing.T) {
	fields := map[language.Language]QuoteFields{
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
	idx := NewIndexer(fields, "")

	enTexts := idx.LowerTexts(language.English)
	if len(enTexts) != 1 {
		t.Errorf("EN LowerTexts: got %d, want 1", len(enTexts))
	}

	jaTexts := idx.LowerTexts(language.Japanese)
	if len(jaTexts) != 2 {
		t.Errorf("JA LowerTexts: got %d, want 2", len(jaTexts))
	}

	enBattler := idx.CharacterIndices(language.English, "10")
	if len(enBattler) != 1 {
		t.Errorf("EN CharacterIndices for Battler: got %d, want 1", len(enBattler))
	}

	jaBeatrice := idx.CharacterIndices(language.Japanese, "27")
	if len(jaBeatrice) != 1 {
		t.Errorf("JA CharacterIndices for Beatrice: got %d, want 1", len(jaBeatrice))
	}
}

func TestIndexer_HasAudio_EmptyDir(t *testing.T) {
	idx := buildTestIndexer()

	if idx.HasAudio() {
		t.Error("HasAudio with empty dir string: expected false")
	}
}

func TestIndexer_HasAudio_NonexistentDir(t *testing.T) {
	fields := map[language.Language]QuoteFields{
		language.English: {
			Texts:        []string{"test"},
			CharacterIDs: []string{"10"},
			Episodes:     []int{1},
			AudioIDs:     []string{""},
		},
	}
	idx := NewIndexer(fields, "/nonexistent/audio/dir")

	if idx.HasAudio() {
		t.Error("HasAudio with nonexistent dir: expected false")
	}
}

func TestIndexer_HasAudio_EmptyExistingDir(t *testing.T) {
	dir := t.TempDir()
	fields := map[language.Language]QuoteFields{
		language.English: {
			Texts:        []string{"test"},
			CharacterIDs: []string{"10"},
			Episodes:     []int{1},
			AudioIDs:     []string{""},
		},
	}
	idx := NewIndexer(fields, dir)

	if idx.HasAudio() {
		t.Error("HasAudio with empty existing dir: expected false")
	}
}

func TestIndexer_HasAudio_WithFiles(t *testing.T) {
	dir := t.TempDir()
	subDir := filepath.Join(dir, "10")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(subDir, "10100001.ogg"), []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	fields := map[language.Language]QuoteFields{
		language.English: {
			Texts:        []string{"test"},
			CharacterIDs: []string{"10"},
			Episodes:     []int{1},
			AudioIDs:     []string{""},
		},
	}
	idx := NewIndexer(fields, dir)

	if !idx.HasAudio() {
		t.Error("HasAudio with files: expected true")
	}
}
