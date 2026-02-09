package quote

import "umineko_quote/internal/dto"

func NewCharacterResponse(characterID string, quotes []dto.ParsedQuote, limit int, offset int) dto.CharacterResponse {
	if quotes == nil {
		quotes = []dto.ParsedQuote{}
	}

	var characterName string
	if characterID != "" {
		characterName = CharacterNames.GetCharacterName(CharacterFromID(characterID))
	}

	total := len(quotes)

	if offset >= total {
		return dto.CharacterResponse{
			CharacterID: characterID,
			Character:   characterName,
			Quotes:      []dto.ParsedQuote{},
			Total:       total,
			Limit:       limit,
			Offset:      offset,
		}
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return dto.CharacterResponse{
		CharacterID: characterID,
		Character:   characterName,
		Quotes:      quotes[offset:end],
		Total:       total,
		Limit:       limit,
		Offset:      offset,
	}
}
