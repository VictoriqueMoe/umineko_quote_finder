package quote

import "testing"

func TestNewSearchResult(t *testing.T) {
	q := ParsedQuote{
		Text:        "test quote",
		CharacterID: "10",
		Character:   "Battler",
		Episode:     1,
	}

	sr := NewSearchResult(q, 100)

	if sr.Quote.Text != "test quote" {
		t.Errorf("Quote.Text: got %q, want %q", sr.Quote.Text, "test quote")
	}
	if sr.Quote.CharacterID != "10" {
		t.Errorf("Quote.CharacterID: got %q, want %q", sr.Quote.CharacterID, "10")
	}
	if sr.Score != 100 {
		t.Errorf("Score: got %d, want 100", sr.Score)
	}
}

func TestNewSearchResult_ZeroScore(t *testing.T) {
	q := ParsedQuote{Text: "something"}
	sr := NewSearchResult(q, 0)

	if sr.Score != 0 {
		t.Errorf("Score: got %d, want 0", sr.Score)
	}
}
