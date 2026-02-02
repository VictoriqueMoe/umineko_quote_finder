package quote

import "testing"

func TestContextResponse_Construction(t *testing.T) {
	before := []ParsedQuote{
		{Text: "Before line 1", CharacterID: "10"},
		{Text: "Before line 2", CharacterID: "27"},
	}
	target := ParsedQuote{Text: "Target quote", CharacterID: "10", AudioID: "10100001"}
	after := []ParsedQuote{
		{Text: "After line 1", CharacterID: "27"},
	}

	resp := ContextResponse{
		Before: before,
		Quote:  target,
		After:  after,
	}

	if len(resp.Before) != 2 {
		t.Errorf("Before length: got %d, want 2", len(resp.Before))
	}
	if resp.Quote.Text != "Target quote" {
		t.Errorf("Quote.Text: got %q, want %q", resp.Quote.Text, "Target quote")
	}
	if len(resp.After) != 1 {
		t.Errorf("After length: got %d, want 1", len(resp.After))
	}
}

func TestContextResponse_EmptyBeforeAndAfter(t *testing.T) {
	resp := ContextResponse{
		Before: nil,
		Quote:  ParsedQuote{Text: "Only quote", CharacterID: "10"},
		After:  nil,
	}

	if resp.Before != nil {
		t.Errorf("Before: got %v, want nil", resp.Before)
	}
	if resp.After != nil {
		t.Errorf("After: got %v, want nil", resp.After)
	}
}
