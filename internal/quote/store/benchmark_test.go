package store

import (
	"testing"

	"umineko_quote/internal/dto"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
)

func BenchmarkConcurrentExactSearch(b *testing.B) {
	n := 50000
	quotes := make([]dto.UminekoQuote, n)
	lowerTexts := make([]string, n)
	indices := make([]int, n)
	for i := 0; i < n; i++ {
		quotes[i] = dto.UminekoQuote{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "some text about witches and magic", CharacterID: "10"}}}
		lowerTexts[i] = "some text about witches and magic"
		indices[i] = i
	}
	quotes[n/2].Text = "Beatrice is the golden witch"
	lowerTexts[n/2] = "beatrice is the golden witch"

	b.ResetTimer()
	for b.Loop() {
		concurrentExactSearchGeneric(indices, lowerTexts, quotes, "beatrice", false, func(idx int) bool { return true })
	}
}
