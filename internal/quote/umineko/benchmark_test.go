package umineko

import (
	"testing"

	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"

	"github.com/VictoriqueMoe/umineko_script_parser/umineko/character"
)

func BenchmarkSearch_BroadQuery(b *testing.B) {
	for b.Loop() {
		if _, err := testService.Search(params.SearchParams{Query: "the", Lang: language.English, Limit: 30}, TruthAll.Scope()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSearch_NarrowQuery(b *testing.B) {
	for b.Loop() {
		if _, err := testService.Search(params.SearchParams{Query: "Beatrice", Lang: language.English, Limit: 30}, TruthAll.Scope()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSearch_WithCharacterFilter(b *testing.B) {
	for b.Loop() {
		if _, err := testService.Search(params.SearchParams{Query: "witch", Lang: language.English, Limit: 30, CharacterID: character.Battler.ID()}, TruthAll.Scope()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSearch_WithEpisodeFilter(b *testing.B) {
	for b.Loop() {
		if _, err := testService.Search(params.SearchParams{Query: "witch", Lang: language.English, Limit: 30, Episode: 1}, TruthAll.Scope()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSearch_WithCharacterAndEpisodeFilter(b *testing.B) {
	for b.Loop() {
		if _, err := testService.Search(params.SearchParams{Query: "witch", Lang: language.English, Limit: 30, CharacterID: character.Battler.ID(), Episode: 1}, TruthAll.Scope()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSearch_RedTruth(b *testing.B) {
	for b.Loop() {
		if _, err := testService.Search(params.SearchParams{Query: "truth", Lang: language.English, Limit: 30}, TruthRed.Scope()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSearch_NoResults(b *testing.B) {
	for b.Loop() {
		if _, err := testService.Search(params.SearchParams{Query: "xyzzyxyzzyxyzzy", Lang: language.English, Limit: 30}, TruthAll.Scope()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSearch_Japanese(b *testing.B) {
	for b.Loop() {
		if _, err := testService.Search(params.SearchParams{Query: "ベアトリーチェ", Lang: language.Japanese, Limit: 30}, TruthAll.Scope()); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetByAudioID(b *testing.B) {
	for b.Loop() {
		if _, err := testService.GetByAudioID(language.English, "11900001"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkGetContext(b *testing.B) {
	for b.Loop() {
		if _, err := testService.GetContext(language.English, "11900001", 5); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRandom_Unfiltered(b *testing.B) {
	for b.Loop() {
		if _, err := testService.Random(language.English, "", 0, params.QuoteScope{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRandom_WithCharacter(b *testing.B) {
	for b.Loop() {
		if _, err := testService.Random(language.English, character.Beatrice.ID(), 0, params.QuoteScope{}); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBrowse(b *testing.B) {
	for b.Loop() {
		if _, err := testService.Browse(params.BrowseParams{Lang: language.English, CharacterID: character.Battler.ID(), Limit: 50}, TruthAll.Scope()); err != nil {
			b.Fatal(err)
		}
	}
}
