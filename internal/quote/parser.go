package quote

import (
	"umineko_quote/internal/dto"
	"umineko_quote/internal/lexar"
)

type Parser interface {
	ParseAll(lines []string) []dto.ParsedQuote
	SubtitleRefs() []lexar.SubtitleRef
	ValidationErrors() []lexar.ValidationError
}

func NewParser() Parser {
	return NewScriptParser()
}
