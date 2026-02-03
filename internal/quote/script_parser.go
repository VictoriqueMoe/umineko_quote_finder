package quote

import (
	"strings"

	"umineko_quote/internal/lexar"
	"umineko_quote/internal/lexar/transformer"
)

// scriptParser implements Parser using the script lexer/parser.
type scriptParser struct {
	extractor *lexar.QuoteExtractor
	factory   *transformer.Factory
}

// ParseAll parses all lines and returns quotes.
func (p *scriptParser) ParseAll(lines []string) []ParsedQuote {
	input := strings.Join(lines, "\n")

	extracted := p.extractor.ExtractQuotes(input)

	var quotes []ParsedQuote
	for _, eq := range extracted {
		quotes = append(quotes, ParsedQuote{
			Text:        p.factory.MustGet(transformer.FormatPlainText).Transform(eq.Content),
			TextHtml:    p.factory.MustGet(transformer.FormatHTML).Transform(eq.Content),
			CharacterID: eq.CharacterID,
			Character:   CharacterNames.GetCharacterName(eq.CharacterID),
			AudioID:     eq.AudioID,
			Episode:     eq.Episode,
			ContentType: eq.ContentType,
			TruthType:   eq.TruthType.String(),
		})
	}

	return quotes
}

// NewScriptParser creates a new parser using the script package.
func NewScriptParser() Parser {
	extractor := lexar.NewQuoteExtractor()

	return &scriptParser{
		extractor: extractor,
		factory:   transformer.NewFactory(extractor.Presets()),
	}
}
