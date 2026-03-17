package params

import (
	"strings"
	"umineko_quote/internal/quote/character"
	"umineko_quote/internal/quote/language"
)

type BrowseParams struct {
	Lang         language.Language
	Character    character.Character
	Limit        int
	Offset       int
	Episode      int
	Truth        string
	InteractionA string
	InteractionB string
}

func NewBrowseParams(
	lang language.Language,
	limit int,
	offset int,
	characterParam string,
	episode int,
	truthParam string,
	interactionAParam string,
	interactionBParam string,
) BrowseParams {
	return BrowseParams{
		Lang:         lang,
		Character:    character.Character(strings.TrimSpace(characterParam)),
		Limit:        limit,
		Offset:       offset,
		Episode:      episode,
		Truth:        strings.TrimSpace(truthParam),
		InteractionA: character.Character(strings.TrimSpace(interactionAParam)).ID(),
		InteractionB: character.Character(strings.TrimSpace(interactionBParam)).ID(),
	}
}
