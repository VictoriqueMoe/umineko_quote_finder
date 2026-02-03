package lexar

import (
	"regexp"
	"strconv"
	"strings"

	"umineko_quote/internal/lexar/ast"
	"umineko_quote/internal/lexar/transformer"
)

type (
	QuoteExtractor struct {
		presets *transformer.PresetContext
	}

	ExtractedQuote struct {
		Content     []ast.DialogueElement
		CharacterID string
		AudioID     string
		Episode     int
		ContentType string
		Truth       TruthFlags
	}
)

var omakeRegex = regexp.MustCompile(`^o(\d+)_`)

func NewQuoteExtractor() *QuoteExtractor {
	return &QuoteExtractor{
		presets: transformer.NewPresetContext(),
	}
}

func (e *QuoteExtractor) ExtractQuotes(input string) []ExtractedQuote {
	script := Parse(input)
	return e.ExtractFromScript(script)
}

func (e *QuoteExtractor) ExtractFromScript(script *ast.Script) []ExtractedQuote {
	e.presets.CollectFromScript(script)

	var quotes []ExtractedQuote
	currentEpisode := 0
	currentContentType := ""

	for _, line := range script.Lines {
		switch l := line.(type) {
		case *ast.EpisodeMarkerLine:
			currentEpisode = l.Episode
			if l.Type == "episode" {
				currentContentType = ""
			} else {
				currentContentType = l.Type
			}

		case *ast.LabelLine:
			if matches := omakeRegex.FindStringSubmatch(l.Name); len(matches) >= 2 {
				if ep, err := strconv.Atoi(matches[1]); err == nil {
					currentEpisode = ep
					currentContentType = "omake"
				}
			}

		case *ast.DialogueLine:
			quote := e.extractFromDialogue(l)
			if quote != nil {
				if currentEpisode > 0 {
					quote.Episode = currentEpisode
				}
				quote.ContentType = currentContentType
				quotes = append(quotes, *quote)
			}
		}
	}

	return quotes
}

func (e *QuoteExtractor) extractFromDialogue(d *ast.DialogueLine) *ExtractedQuote {
	voices := d.GetVoiceCommands()
	truth := DetectTruth(d.Content, e.presets)

	if len(voices) == 0 {
		return &ExtractedQuote{
			Content:     d.Content,
			CharacterID: "narrator",
			Truth:       truth,
		}
	}

	characterID := voices[0].CharacterID

	seen := make(map[string]bool)
	var audioIDs []string
	for _, v := range voices {
		if !seen[v.AudioID] {
			seen[v.AudioID] = true
			audioIDs = append(audioIDs, v.AudioID)
		}
	}

	episode := 0
	if len(audioIDs) > 0 && len(audioIDs[0]) > 0 {
		ep := int(audioIDs[0][0] - '0')
		if ep >= 1 && ep <= 8 {
			episode = ep
		}
	}

	return &ExtractedQuote{
		Content:     d.Content,
		CharacterID: characterID,
		AudioID:     strings.Join(audioIDs, ", "),
		Episode:     episode,
		Truth:       truth,
	}
}

// Presets returns the preset context for use with transformers.
func (e *QuoteExtractor) Presets() *transformer.PresetContext {
	return e.presets
}
