package params

import (
	"strings"
	"umineko_quote/internal/quote/language"
)

type SearchParams struct {
	Query        string
	Lang         language.Language
	Limit        int
	Offset       int
	CharacterID  string
	Episode      int
	InteractionA string
	InteractionB string
}

func NewSearchParams(
	query string,
	lang language.Language,
	limit int,
	offset int,
	characterID string,
	episode int,
	interactionAParam string,
	interactionBParam string,
) SearchParams {
	return SearchParams{
		Query:        strings.TrimSpace(query),
		Lang:         lang,
		Limit:        limit,
		Offset:       offset,
		CharacterID:  strings.TrimSpace(characterID),
		Episode:      episode,
		InteractionA: strings.TrimSpace(interactionAParam),
		InteractionB: strings.TrimSpace(interactionBParam),
	}
}
