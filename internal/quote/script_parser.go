package quote

import (
	"runtime"
	"strings"
	"sync"
	"umineko_quote/internal/quote/character"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/lexar"
	"umineko_quote/internal/lexar/transformer"
)

type scriptParser struct {
	extractor *lexar.QuoteExtractor
	factory   *transformer.Factory
}

func NewScriptParser() Parser {
	extractor := lexar.NewQuoteExtractor()

	return &scriptParser{
		extractor: extractor,
		factory:   transformer.NewFactory(extractor.Presets()),
	}
}

func (p *scriptParser) ParseAll(lines []string) []dto.ParsedQuote {
	filtered := make([]string, 0, len(lines)/8)
	for _, line := range lines {
		if len(line) < 2 {
			continue
		}
		switch line[0] {
		case 'd': // d, d2 dialogue lines
			if line[1] == ' ' || (line[1] == '2' && len(line) > 2 && line[2] == ' ') {
				filtered = append(filtered, line)
			}
		case 'p': // preset_define
			if len(line) > 13 && line[:13] == "preset_define" {
				filtered = append(filtered, line)
			}
		case 'n': // new_episode, new_tea, new_ura
			if len(line) > 4 && line[:4] == "new_" {
				filtered = append(filtered, line)
			}
		case '*': // labels
			filtered = append(filtered, line)
		case 's': // stralias, ssa_load (subtitle references)
			if len(line) > 8 && line[:8] == "stralias" {
				filtered = append(filtered, line)
			} else if len(line) > 8 && line[:8] == "ssa_load" {
				filtered = append(filtered, line)
			}
		case 'l': // lv (top-level voice/video commands)
			if len(line) > 2 && line[:2] == "lv" && line[2] == ' ' {
				filtered = append(filtered, line)
			}
		}
	}

	input := strings.Join(filtered, "\n")

	extracted := p.extractor.ExtractQuotes(input)

	quotes := make([]dto.ParsedQuote, len(extracted))

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

				var audioTextMap map[string]string
				if len(eq.AudioTextMap) > 0 {
					audioTextMap = make(map[string]string, len(eq.AudioTextMap))
					for audioID, fragment := range eq.AudioTextMap {
						audioTextMap[audioID] = plainText.Transform(fragment)
					}
				}

				quotes[i] = dto.ParsedQuote{
					Text:         plainText.Transform(eq.Content),
					TextHtml:     htmlText.Transform(eq.Content),
					CharacterID:  eq.CharacterID,
					Character:    character.CharacterNames.GetCharacterName(character.CharacterFromID(eq.CharacterID)),
					AudioID:      eq.AudioID,
					AudioCharMap: eq.AudioCharMap,
					AudioTextMap: audioTextMap,
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

func (p *scriptParser) SubtitleRefs() []lexar.SubtitleRef {
	return p.extractor.SubtitleRefs()
}

func (p *scriptParser) ValidationErrors() []lexar.ValidationError {
	return p.extractor.ValidationErrors()
}
