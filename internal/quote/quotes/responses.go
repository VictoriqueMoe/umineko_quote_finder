package quotes

import (
	"maps"

	"umineko_quote/internal/dto"
)

func NewCharactersResult(curated map[string]string, speakers map[string]string) dto.CharactersResult {
	characters := make(map[string]string, len(curated))
	maps.Copy(characters, curated)

	additional := make(map[string]string, len(speakers))
	maps.Copy(additional, speakers)

	for id := range characters {
		delete(additional, id)
	}
	delete(additional, "")

	return dto.CharactersResult{
		Characters: characters,
		Additional: additional,
	}
}
