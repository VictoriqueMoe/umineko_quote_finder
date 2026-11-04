package umineko

import (
	"testing"

	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"

	"github.com/VictoriqueMoe/umineko_script_parser/umineko/character"
)

func BenchmarkSearch_BroadQuery(b *testing.B) {
	svc := testService
	for b.Loop() {
		searchWithParams(b, svc, "the", language.English, 30, 0, "", 0, TruthAll, "", "")
	}
}

func BenchmarkSearch_NarrowQuery(b *testing.B) {
	svc := testService
	for b.Loop() {
		searchWithParams(b, svc, "Beatrice", language.English, 30, 0, "", 0, TruthAll, "", "")
	}
}

func BenchmarkSearch_WithCharacterFilter(b *testing.B) {
	svc := testService
	for b.Loop() {
		searchWithParams(b, svc, "witch", language.English, 30, 0, character.Battler, 0, TruthAll, "", "")
	}
}

func BenchmarkSearch_WithEpisodeFilter(b *testing.B) {
	svc := testService
	for b.Loop() {
		searchWithParams(b, svc, "witch", language.English, 30, 0, "", 1, TruthAll, "", "")
	}
}

func BenchmarkSearch_WithCharacterAndEpisodeFilter(b *testing.B) {
	svc := testService
	for b.Loop() {
		searchWithParams(b, svc, "witch", language.English, 30, 0, character.Battler, 1, TruthAll, "", "")
	}
}

func BenchmarkSearch_RedTruth(b *testing.B) {
	svc := testService
	for b.Loop() {
		searchWithParams(b, svc, "truth", language.English, 30, 0, "", 0, TruthRed, "", "")
	}
}

func BenchmarkSearch_NoResults(b *testing.B) {
	svc := testService
	for b.Loop() {
		searchWithParams(b, svc, "xyzzyxyzzyxyzzy", language.English, 30, 0, "", 0, TruthAll, "", "")
	}
}

func BenchmarkSearch_Japanese(b *testing.B) {
	svc := testService
	for b.Loop() {
		searchWithParams(b, svc, "ベアトリーチェ", language.Japanese, 30, 0, "", 0, TruthAll, "", "")
	}
}

func BenchmarkGetByAudioID(b *testing.B) {
	for b.Loop() {
		byAudioID(b, language.English, "11900001")
	}
}

func BenchmarkGetContext(b *testing.B) {
	for b.Loop() {
		context(b, language.English, "11900001", 5)
	}
}

func BenchmarkRandom_Unfiltered(b *testing.B) {
	for b.Loop() {
		randomQuote(b, language.English, "", 0, params.QuoteScope{})
	}
}

func BenchmarkRandom_WithCharacter(b *testing.B) {
	for b.Loop() {
		randomQuote(b, language.English, character.Beatrice.ID(), 0, params.QuoteScope{})
	}
}

func BenchmarkBrowse(b *testing.B) {
	svc := testService
	for b.Loop() {
		browseWithParams(b, svc, language.English, character.Battler, 50, 0, 0, TruthAll, "", "")
	}
}
