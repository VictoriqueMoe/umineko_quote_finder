package params

import (
	"strings"
	"umineko_quote/internal/quote/language"
)

type BrowseParams struct {
	Lang         language.Language
	CharacterID  string
	Limit        int
	Offset       int
	Episode      int
	InteractionA string
	InteractionB string
}

func NewBrowseParams(
	lang language.Language,
	limit int,
	offset int,
	characterID string,
	episode int,
	interactionAParam string,
	interactionBParam string,
) BrowseParams {
	return BrowseParams{
		Lang:         lang,
		CharacterID:  strings.TrimSpace(characterID),
		Limit:        limit,
		Offset:       offset,
		Episode:      episode,
		InteractionA: strings.TrimSpace(interactionAParam),
		InteractionB: strings.TrimSpace(interactionBParam),
	}
}
