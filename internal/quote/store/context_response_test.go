package store

import (
	"testing"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"

	"umineko_quote/internal/dto"
)

func TestContextResponse_Construction(t *testing.T) {
	before := []dto.UminekoQuote{
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Before line 1", CharacterID: "10"}}},
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Before line 2", CharacterID: "27"}}},
	}
	target := dto.UminekoQuote{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Target quote", CharacterID: "10", AudioID: "10100001"}}}
	after := []dto.UminekoQuote{
		{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "After line 1", CharacterID: "27"}}},
	}

	resp := dto.ContextResponse{
		Before: before,
		Quote:  target,
		After:  after,
	}

	respBefore := resp.Before.([]dto.UminekoQuote)
	if len(respBefore) != 2 {
		t.Errorf("Before length: got %d, want 2", len(respBefore))
	}
	respQuote := resp.Quote.(dto.UminekoQuote)
	if respQuote.Text != "Target quote" {
		t.Errorf("Quote.Text: got %q, want %q", respQuote.Text, "Target quote")
	}
	respAfter := resp.After.([]dto.UminekoQuote)
	if len(respAfter) != 1 {
		t.Errorf("After length: got %d, want 1", len(respAfter))
	}
}

func TestContextResponse_EmptyBeforeAndAfter(t *testing.T) {
	resp := dto.ContextResponse{
		Before: nil,
		Quote:  dto.UminekoQuote{UminekoQuote: scriptdto.UminekoQuote{BaseQuote: scriptdto.BaseQuote{Text: "Only quote", CharacterID: "10"}}},
		After:  nil,
	}

	if resp.Before != nil {
		t.Errorf("Before: got %v, want nil", resp.Before)
	}
	if resp.After != nil {
		t.Errorf("After: got %v, want nil", resp.After)
	}
}
