package quote

import (
	"testing"
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/character"
	"umineko_quote/internal/quote/language"
)

func BenchmarkSearch_BroadQuery(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		svc.Search("the", language.English, 30, 0, "", 0, TruthAll)
	}
}

func BenchmarkSearch_NarrowQuery(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		svc.Search("Beatrice", language.English, 30, 0, "", 0, TruthAll)
	}
}

func BenchmarkSearch_WithCharacterFilter(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		svc.Search("witch", language.English, 30, 0, character.Battler, 0, TruthAll)
	}
}

func BenchmarkSearch_WithEpisodeFilter(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		svc.Search("witch", language.English, 30, 0, "", 1, TruthAll)
	}
}

func BenchmarkSearch_WithCharacterAndEpisodeFilter(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		svc.Search("witch", language.English, 30, 0, character.Battler, 1, TruthAll)
	}
}

func BenchmarkSearch_RedTruth(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		svc.Search("truth", language.English, 30, 0, "", 0, TruthRed)
	}
}

func BenchmarkSearch_NoResults(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		svc.Search("xyzzyxyzzyxyzzy", language.English, 30, 0, "", 0, TruthAll)
	}
}

func BenchmarkSearch_Japanese(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		svc.Search("ベアトリーチェ", language.Japanese, 30, 0, "", 0, TruthAll)
	}
}

func BenchmarkGetByAudioID(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		svc.GetByAudioID(language.English, "11900001")
	}
}

func BenchmarkGetContext(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		svc.GetContext(language.English, "11900001", 5)
	}
}

func BenchmarkRandom_Unfiltered(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		svc.Random(language.English, "", 0, TruthAll)
	}
}

func BenchmarkRandom_WithCharacter(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		svc.Random(language.English, character.Beatrice, 0, TruthAll)
	}
}

func BenchmarkBrowse(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		svc.Browse(language.English, character.Battler, 50, 0, 0, TruthAll)
	}
}

func BenchmarkIndexer_FilteredIndices_CharOnly(b *testing.B) {
	idx, _ := buildTestIndexerLarge()
	b.ResetTimer()
	for b.Loop() {
		idx.FilteredIndices(language.English, "10", 0)
	}
}

func BenchmarkIndexer_FilteredIndices_EpOnly(b *testing.B) {
	idx, _ := buildTestIndexerLarge()
	b.ResetTimer()
	for b.Loop() {
		idx.FilteredIndices(language.English, "", 1)
	}
}

func BenchmarkIndexer_FilteredIndices_Both(b *testing.B) {
	idx, _ := buildTestIndexerLarge()
	b.ResetTimer()
	for b.Loop() {
		idx.FilteredIndices(language.English, "10", 1)
	}
}

func BenchmarkConcurrentExactSearch(b *testing.B) {
	n := 50000
	quotes := make([]dto.ParsedQuote, n)
	lowerTexts := make([]string, n)
	indices := make([]int, n)
	for i := 0; i < n; i++ {
		quotes[i] = dto.ParsedQuote{Text: "some text about witches and magic", CharacterID: "10"}
		lowerTexts[i] = "some text about witches and magic"
		indices[i] = i
	}
	quotes[n/2].Text = "Beatrice is the golden witch"
	lowerTexts[n/2] = "beatrice is the golden witch"

	b.ResetTimer()
	for b.Loop() {
		concurrentExactSearch(indices, lowerTexts, quotes, "beatrice", func(q dto.ParsedQuote) bool { return true })
	}
}

func buildTestIndexerLarge() (Indexer, map[language.Language][]dto.ParsedQuote) {
	chars := []string{"10", "27", "narrator", "19", "00"}
	n := 10000
	quotes := make([]dto.ParsedQuote, n)
	for i := 0; i < n; i++ {
		quotes[i] = dto.ParsedQuote{
			Text:        "Test quote text number",
			CharacterID: chars[i%len(chars)],
			Episode:     (i % 8) + 1,
		}
	}
	m := map[language.Language][]dto.ParsedQuote{
		language.English: quotes,
	}
	return NewIndexer(m, ""), m
}
