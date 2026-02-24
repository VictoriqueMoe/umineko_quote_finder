package engine

import (
	"fmt"
	"testing"
	"umineko_quote/internal/quote/character"

	"umineko_quote/internal/dto"
)

func TestKanonAttributionEngine_FixesErikaToKanon(t *testing.T) {
	audioIDs := []string{
		"20600528", "20600529",
		"20600530", "20600531", "20600532", "20600533",
		"20600534", "20600535", "20600536", "20600537", "20600538", "20600539",
		"20600540", "20600541",
		"20600542", "20600543",
	}

	engine := &KanonAttributionEngine{}

	for _, audioID := range audioIDs {
		t.Run(audioID, func(t *testing.T) {
			quotes := []dto.ParsedQuote{
				{
					CharacterID: erikaID,
					Character:   "Erika",
					AudioID:     audioID,
					Episode:     2,
				},
			}

			result := engine.Apply(quotes)

			if result[0].CharacterID != kanonID {
				t.Errorf("expected characterID %s, got %s", kanonID, result[0].CharacterID)
			}
			if result[0].Character != kanonName {
				t.Errorf("expected character %s, got %s", kanonName, result[0].Character)
			}
		})
	}
}

func TestKanonAttributionEngine_FixesAudioCharMap(t *testing.T) {
	quotes := []dto.ParsedQuote{
		{
			CharacterID: erikaID,
			Character:   "Erika",
			AudioID:     "20600528",
			AudioCharMap: map[string]string{
				"20600528": erikaID,
			},
		},
	}

	engine := &KanonAttributionEngine{}
	result := engine.Apply(quotes)

	if result[0].AudioCharMap["20600528"] != kanonID {
		t.Errorf("expected AudioCharMap entry %s, got %s", kanonID, result[0].AudioCharMap["20600528"])
	}
}

func TestKanonAttributionEngine_IgnoresNonErikaQuotes(t *testing.T) {
	quotes := []dto.ParsedQuote{
		{
			CharacterID: "10",
			Character:   "Battler",
			AudioID:     "20600528",
		},
	}

	engine := &KanonAttributionEngine{}
	result := engine.Apply(quotes)

	if result[0].CharacterID != character.Battler.ID() {
		t.Errorf("should not modify non-Erika quotes, got %s", result[0].CharacterID)
	}
}

func TestKanonAttributionEngine_IgnoresErikaOutsideRange(t *testing.T) {
	quotes := []dto.ParsedQuote{
		{
			CharacterID: erikaID,
			Character:   "Erika",
			AudioID:     "50500100",
		},
	}

	engine := &KanonAttributionEngine{}
	result := engine.Apply(quotes)

	if result[0].CharacterID != erikaID {
		t.Errorf("should not modify Erika quotes outside audio range, got %s", result[0].CharacterID)
	}
}

func TestKanonAttributionEngine_HandlesCompositeAudioID(t *testing.T) {
	quotes := []dto.ParsedQuote{
		{
			CharacterID: erikaID,
			Character:   "Erika",
			AudioID:     fmt.Sprintf("20600528, 20600529"),
		},
	}

	engine := &KanonAttributionEngine{}
	result := engine.Apply(quotes)

	if result[0].CharacterID != kanonID {
		t.Errorf("expected characterID %s for composite audioID, got %s", kanonID, result[0].CharacterID)
	}
}
