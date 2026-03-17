package params

import (
	"strings"
	"umineko_quote/internal/quote/character"
	"umineko_quote/internal/quote/language"
)

type SearchParams struct {
	Query        string
	Lang         language.Language
	Limit        int
	Offset       int
	Character    character.Character
	Episode      int
	Truth        string
	InteractionA string
	InteractionB string
}

func NewSearchParams(
	query string,
	lang language.Language,
	limit int,
	offset int,
	characterParam string,
	episode int,
	truthParam string,
	interactionAParam string,
	interactionBParam string,
) SearchParams {
	return SearchParams{
		Query:        strings.TrimSpace(query),
		Lang:         lang,
		Limit:        limit,
		Offset:       offset,
		Character:    character.Character(strings.TrimSpace(characterParam)),
		Episode:      episode,
		Truth:        strings.TrimSpace(truthParam),
		InteractionA: character.Character(strings.TrimSpace(interactionAParam)).ID(),
		InteractionB: character.Character(strings.TrimSpace(interactionBParam)).ID(),
	}
}
