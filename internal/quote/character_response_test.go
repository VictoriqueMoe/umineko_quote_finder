package quote

import "testing"

func TestNewCharacterResponse_NilQuotes(t *testing.T) {
	resp := NewCharacterResponse("10", nil, 10, 0)

	if len(resp.Quotes) != 0 {
		t.Errorf("Quotes length: got %d, want 0", len(resp.Quotes))
	}
	if resp.Total != 0 {
		t.Errorf("Total: got %d, want 0", resp.Total)
	}
	if resp.CharacterID != "10" {
		t.Errorf("CharacterID: got %q, want %q", resp.CharacterID, "10")
	}
	if resp.Character != "Battler" {
		t.Errorf("Character: got %q, want %q", resp.Character, "Battler")
	}
}

func TestNewCharacterResponse_EmptyCharacterID(t *testing.T) {
	resp := NewCharacterResponse("", nil, 10, 0)

	if resp.CharacterID != "" {
		t.Errorf("CharacterID: got %q, want empty", resp.CharacterID)
	}
	if resp.Character != "" {
		t.Errorf("Character: got %q, want empty", resp.Character)
	}
}

func TestNewCharacterResponse_Pagination(t *testing.T) {
	quotes := make([]ParsedQuote, 25)
	for i := 0; i < 25; i++ {
		quotes[i] = ParsedQuote{Text: "quote", CharacterID: "27"}
	}

	resp := NewCharacterResponse("27", quotes, 10, 0)

	if len(resp.Quotes) != 10 {
		t.Errorf("Quotes length: got %d, want 10", len(resp.Quotes))
	}
	if resp.Total != 25 {
		t.Errorf("Total: got %d, want 25", resp.Total)
	}
	if resp.Character != "Beatrice" {
		t.Errorf("Character: got %q, want %q", resp.Character, "Beatrice")
	}
}

func TestNewCharacterResponse_OffsetBeyondTotal(t *testing.T) {
	quotes := make([]ParsedQuote, 5)
	for i := 0; i < 5; i++ {
		quotes[i] = ParsedQuote{Text: "quote"}
	}

	resp := NewCharacterResponse("10", quotes, 10, 100)

	if len(resp.Quotes) != 0 {
		t.Errorf("Quotes length: got %d, want 0", len(resp.Quotes))
	}
	if resp.Total != 5 {
		t.Errorf("Total: got %d, want 5", resp.Total)
	}
}

func TestNewCharacterResponse_PartialLastPage(t *testing.T) {
	quotes := make([]ParsedQuote, 25)
	for i := 0; i < 25; i++ {
		quotes[i] = ParsedQuote{Text: "quote"}
	}

	resp := NewCharacterResponse("10", quotes, 10, 20)

	if len(resp.Quotes) != 5 {
		t.Errorf("Quotes length: got %d, want 5", len(resp.Quotes))
	}
	if resp.Total != 25 {
		t.Errorf("Total: got %d, want 25", resp.Total)
	}
}
