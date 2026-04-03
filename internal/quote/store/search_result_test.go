package store

import (
	"testing"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/umineko/character"

	"umineko_quote/internal/dto"
)

func TestNewSearchResult(t *testing.T) {
	q := dto.UminekoQuote{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{
		Text:        "test quote",
		CharacterID: "10",
		Character:   character.CharacterNames[character.Battler],
		Episode:     1,
	}}}

	sr := NewSearchResult(q, 100)

	uq := sr.Quote.(dto.UminekoQuote)
	if uq.Text != "test quote" {
		t.Errorf("Quote.Text: got %q, want %q", uq.Text, "test quote")
	}
	if uq.CharacterID != "10" {
		t.Errorf("Quote.CharacterID: got %q, want %q", uq.CharacterID, "10")
	}
	if sr.Score != 100 {
		t.Errorf("Score: got %d, want 100", sr.Score)
	}
}

func TestNewSearchResult_ZeroScore(t *testing.T) {
	q := dto.UminekoQuote{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "something"}}}
	sr := NewSearchResult(q, 0)

	if sr.Score != 0 {
		t.Errorf("Score: got %d, want 0", sr.Score)
	}
}
