package quote

import (
	"umineko_quote/internal/dto"
	"umineko_quote/internal/lexar"
)

type Parser interface {
	ParseAll(lines []string) []dto.ParsedQuote
	SubtitleRefs() []lexar.SubtitleRef
}

func NewParser() Parser {
	return NewScriptParser()
}
