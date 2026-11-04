package memory

import (
	"testing"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/game"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
)

func BenchmarkIndexer_FilteredIndices_CharOnly(b *testing.B) {
	idx := buildTestIndexerLarge()
	for b.Loop() {
		idx.FilteredIndices(language.English, "10", 0)
	}
}

func BenchmarkIndexer_FilteredIndices_EpOnly(b *testing.B) {
	idx := buildTestIndexerLarge()
	for b.Loop() {
		idx.FilteredIndices(language.English, "", 1)
	}
}

func BenchmarkIndexer_FilteredIndices_Both(b *testing.B) {
	idx := buildTestIndexerLarge()
	for b.Loop() {
		idx.FilteredIndices(language.English, "10", 1)
	}
}

func buildTestIndexerLarge() *indexer {
	chars := []string{"10", "27", "narrator", "19", "00"}
	n := 10000
	texts := make([]string, n)
	charIDs := make([]string, n)
	episodes := make([]int, n)
	for i := range n {
		texts[i] = "Test quote text number"
		charIDs[i] = chars[i%len(chars)]
		episodes[i] = (i % 8) + 1
	}
	fields := map[language.Language]quoteFields{
		language.English: {
			Texts:        texts,
			CharacterIDs: charIDs,
			Episodes:     episodes,
			AudioIDs:     make([]string, n),
		},
	}
	return newIndexer(fields)
}

func BenchmarkConcurrentMatch(b *testing.B) {
	n := 50000
	lowerTexts := make([]string, n)
	indices := make([]int, n)
	for i := range n {
		lowerTexts[i] = "some text about witches and magic"
		indices[i] = i
	}
	lowerTexts[n/2] = "beatrice is the golden witch"

	for b.Loop() {
		concurrentMatch(indices, lowerTexts, "beatrice", false, keepAll)
	}
}

func BenchmarkFindQuotes_Text(b *testing.B) {
	n := 50000
	en := make([]dto.UminekoQuote, n)
	for i := range n {
		en[i] = dto.UminekoQuote{
			UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "some text about witches and magic", CharacterID: "10"}},
			Index:        i,
		}
	}
	en[n/2].Text = "Beatrice is the golden witch"

	d := NewQuotes(game.Umineko{}, map[language.Language][]dto.UminekoQuote{language.English: en})
	f := params.QuoteFilter{Lang: language.English, Text: "beatrice"}

	for b.Loop() {
		if _, _, err := d.FindQuotes(f, 30, 0); err != nil {
			b.Fatal(err)
		}
	}
}
