package quote

import (
	"runtime"
	"strings"
	"sync"

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

	quotes := make([]ParsedQuote, len(extracted))

	plainText := p.factory.MustGet(transformer.FormatPlainText)
	htmlText := p.factory.MustGet(transformer.FormatHTML)

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (len(extracted) + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > len(extracted) {
			end = len(extracted)
		}
		if start >= end {
			break
		}

		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				eq := &extracted[i]
				quotes[i] = ParsedQuote{
					Text:        plainText.Transform(eq.Content),
					TextHtml:    htmlText.Transform(eq.Content),
					CharacterID: eq.CharacterID,
					Character:   CharacterNames.GetCharacterName(eq.CharacterID),
					AudioID:     eq.AudioID,
					Episode:     eq.Episode,
					ContentType: eq.ContentType,
					TruthType:   eq.TruthType.String(),
				}
			}
		}(start, end)
	}
	wg.Wait()

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
