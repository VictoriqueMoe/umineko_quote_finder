package store

import (
	"testing"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/umineko/character"

	"umineko_quote/internal/dto"
)

func resolveName(id string) string {
	return character.CharacterNames.GetCharacterName(character.CharacterFromID(id))
}

func TestNewCharacterResponse_NilQuotes(t *testing.T) {
	resp := newCharacterResponse("10", []dto.UminekoQuote{}, 10, 0, resolveName)

	quotes := resp.Quotes.([]dto.UminekoQuote)
	if len(quotes) != 0 {
		t.Errorf("Quotes length: got %d, want 0", len(quotes))
	}
	if resp.Total != 0 {
		t.Errorf("Total: got %d, want 0", resp.Total)
	}
	if resp.CharacterID != "10" {
		t.Errorf("CharacterID: got %q, want %q", resp.CharacterID, "10")
	}
	if resp.Character != character.CharacterNames[character.Battler] {
		t.Errorf("Character: got %q, want %q", resp.Character, character.CharacterNames[character.Battler])
	}
}

func TestNewCharacterResponse_EmptyCharacterID(t *testing.T) {
	resp := newCharacterResponse("", []dto.UminekoQuote{}, 10, 0, resolveName)

	if resp.CharacterID != "" {
		t.Errorf("CharacterID: got %q, want empty", resp.CharacterID)
	}
	if resp.Character != "" {
		t.Errorf("Character: got %q, want empty", resp.Character)
	}
}

func TestNewCharacterResponse_Pagination(t *testing.T) {
	quotes := make([]dto.UminekoQuote, 25)
	for i := 0; i < 25; i++ {
		quotes[i] = dto.UminekoQuote{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "quote", CharacterID: "27"}}}
	}

	resp := newCharacterResponse("27", quotes, 10, 0, resolveName)

	respQuotes := resp.Quotes.([]dto.UminekoQuote)
	if len(respQuotes) != 10 {
		t.Errorf("Quotes length: got %d, want 10", len(respQuotes))
	}
	if resp.Total != 25 {
		t.Errorf("Total: got %d, want 25", resp.Total)
	}
	if resp.Character != character.CharacterNames[character.Beatrice] {
		t.Errorf("Character: got %q, want %q", resp.Character, character.CharacterNames[character.Beatrice])
	}
}

func TestNewCharacterResponse_OffsetBeyondTotal(t *testing.T) {
	quotes := make([]dto.UminekoQuote, 5)
	for i := 0; i < 5; i++ {
		quotes[i] = dto.UminekoQuote{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "quote"}}}
	}

	resp := newCharacterResponse("10", quotes, 10, 100, resolveName)

	respQuotes := resp.Quotes.([]dto.UminekoQuote)
	if len(respQuotes) != 0 {
		t.Errorf("Quotes length: got %d, want 0", len(respQuotes))
	}
	if resp.Total != 5 {
		t.Errorf("Total: got %d, want 5", resp.Total)
	}
}

func TestNewCharacterResponse_PartialLastPage(t *testing.T) {
	quotes := make([]dto.UminekoQuote, 25)
	for i := 0; i < 25; i++ {
		quotes[i] = dto.UminekoQuote{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "quote"}}}
	}

	resp := newCharacterResponse("10", quotes, 10, 20, resolveName)

	respQuotes := resp.Quotes.([]dto.UminekoQuote)
	if len(respQuotes) != 5 {
		t.Errorf("Quotes length: got %d, want 5", len(respQuotes))
	}
	if resp.Total != 25 {
		t.Errorf("Total: got %d, want 25", resp.Total)
	}
}
