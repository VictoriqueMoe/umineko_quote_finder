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
	// Pre-filter to only relevant lines (dialogue, presets, episode markers, labels)
	filtered := make([]string, 0, len(lines)/8)
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		// Fast prefix check without allocations
		switch line[0] {
		case 'd':
			// d or d2 dialogue
			if line[1] == ' ' || (line[1] == '2' && len(line) > 2 && line[2] == ' ') {
				filtered = append(filtered, line)
			}
		case 'p':
			// preset_define
			if len(line) > 13 && line[:13] == "preset_define" {
				filtered = append(filtered, line)
			}
		case 'n':
			// new_episode, new_tea, new_ura
			if len(line) > 4 && line[:4] == "new_" {
				filtered = append(filtered, line)
			}
		case '*':
			// labels (for omake detection)
			filtered = append(filtered, line)
		}
	}

	input := strings.Join(filtered, "\n")

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
		wg.Go(func() {
			for i := start; i < end; i++ {
				eq := &extracted[i]
				quotes[i] = ParsedQuote{
					Text:         plainText.Transform(eq.Content),
					TextHtml:     htmlText.Transform(eq.Content),
					CharacterID:  eq.CharacterID,
					Character:    CharacterNames.GetCharacterName(eq.CharacterID),
					AudioID:      eq.AudioID,
					Episode:      eq.Episode,
					ContentType:  eq.ContentType,
					HasRedTruth:  eq.Truth.HasRed,
					HasBlueTruth: eq.Truth.HasBlue,
				}
			}
		})
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
