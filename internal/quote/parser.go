package quote

import "umineko_quote/internal/dto"

type Parser interface {
	ParseAll(lines []string) []dto.ParsedQuote
}

func NewParser() Parser {
	return NewScriptParser()
}
