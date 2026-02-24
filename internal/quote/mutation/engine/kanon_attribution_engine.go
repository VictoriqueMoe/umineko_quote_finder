package engine

import (
	"strings"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/character"
)

var (
	kanonID   = character.Kanon.ID()
	kanonName = character.CharacterNames.GetCharacterName(character.Kanon)
	erikaID   = character.Erika.ID()
)

type KanonAttributionEngine struct{}

func (k *KanonAttributionEngine) Apply(quotes []dto.ParsedQuote) []dto.ParsedQuote {
	for i := range quotes {
		if quotes[i].CharacterID != erikaID {
			continue
		}
		if !containsKanonAudioID(quotes[i].AudioID) {
			continue
		}

		quotes[i].CharacterID = kanonID
		quotes[i].Character = kanonName

		if quotes[i].AudioCharMap != nil {
			for audioID := range quotes[i].AudioCharMap {
				if isKanonEp2AudioID(audioID) {
					quotes[i].AudioCharMap[audioID] = kanonID
				}
			}
		}
	}
	return quotes
}

func containsKanonAudioID(audioIDField string) bool {
	for _, id := range strings.Split(audioIDField, ", ") {
		if isKanonEp2AudioID(id) {
			return true
		}
	}
	return false
}

func isKanonEp2AudioID(audioID string) bool {
	return len(audioID) == 8 && audioID >= "20600528" && audioID <= "20600543"
}
