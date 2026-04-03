package umineko

import (
	"testing"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/store"

	"github.com/VictoriqueMoe/umineko_script_parser/umineko/character"
)

func BenchmarkSearch_BroadQuery(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		searchWithParams(svc, "the", language.English, 30, 0, "", 0, TruthAll, "", "")
	}
}

func BenchmarkSearch_NarrowQuery(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		searchWithParams(svc, "Beatrice", language.English, 30, 0, "", 0, TruthAll, "", "")
	}
}

func BenchmarkSearch_WithCharacterFilter(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		searchWithParams(svc, "witch", language.English, 30, 0, character.Battler, 0, TruthAll, "", "")
	}
}

func BenchmarkSearch_WithEpisodeFilter(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		searchWithParams(svc, "witch", language.English, 30, 0, "", 1, TruthAll, "", "")
	}
}

func BenchmarkSearch_WithCharacterAndEpisodeFilter(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		searchWithParams(svc, "witch", language.English, 30, 0, character.Battler, 1, TruthAll, "", "")
	}
}

func BenchmarkSearch_RedTruth(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		searchWithParams(svc, "truth", language.English, 30, 0, "", 0, TruthRed, "", "")
	}
}

func BenchmarkSearch_NoResults(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		searchWithParams(svc, "xyzzyxyzzyxyzzy", language.English, 30, 0, "", 0, TruthAll, "", "")
	}
}

func BenchmarkSearch_Japanese(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		searchWithParams(svc, "ベアトリーチェ", language.Japanese, 30, 0, "", 0, TruthAll, "", "")
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
		svc.Random(language.English, character.Beatrice.ID(), 0, TruthAll)
	}
}

func BenchmarkBrowse(b *testing.B) {
	svc := testService
	b.ResetTimer()
	for b.Loop() {
		browseWithParams(svc, language.English, character.Battler, 50, 0, 0, TruthAll, "", "")
	}
}

func BenchmarkIndexer_FilteredIndices_CharOnly(b *testing.B) {
	idx := buildTestIndexerLarge()
	b.ResetTimer()
	for b.Loop() {
		idx.FilteredIndices(language.English, "10", 0)
	}
}

func BenchmarkIndexer_FilteredIndices_EpOnly(b *testing.B) {
	idx := buildTestIndexerLarge()
	b.ResetTimer()
	for b.Loop() {
		idx.FilteredIndices(language.English, "", 1)
	}
}

func BenchmarkIndexer_FilteredIndices_Both(b *testing.B) {
	idx := buildTestIndexerLarge()
	b.ResetTimer()
	for b.Loop() {
		idx.FilteredIndices(language.English, "10", 1)
	}
}

func buildTestIndexerLarge() store.Indexer {
	chars := []string{"10", "27", "narrator", "19", "00"}
	n := 10000
	texts := make([]string, n)
	charIDs := make([]string, n)
	episodes := make([]int, n)
	audioIDs := make([]string, n)
	for i := 0; i < n; i++ {
		texts[i] = "Test quote text number"
		charIDs[i] = chars[i%len(chars)]
		episodes[i] = (i % 8) + 1
		audioIDs[i] = ""
	}
	fields := map[language.Language]store.QuoteFields{
		language.English: {
			Texts:        texts,
			CharacterIDs: charIDs,
			Episodes:     episodes,
			AudioIDs:     audioIDs,
		},
	}
	return store.NewIndexer(fields, "")
}
