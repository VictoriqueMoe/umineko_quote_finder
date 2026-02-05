package lexar

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"umineko_quote/internal/lexar/ast"
	"umineko_quote/internal/lexar/transformer"
)

type (
	QuoteExtractor struct {
		presets *transformer.PresetContext
	}

	ExtractedQuote struct {
		Content      []ast.DialogueElement
		CharacterID  string
		AudioID      string
		AudioCharMap map[string]string // audioID → characterID, only for multi-character quotes
		Episode      int
		ContentType  string
		Truth        TruthFlags
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

	if len(voices) == 0 || hasWordsBeforeVoice(d.Content) {
		return &ExtractedQuote{
			Content:     d.Content,
			CharacterID: "narrator",
			Truth:       truth,
		}
	}

	characterID := voices[0].CharacterID

	seen := make(map[string]bool)
	var audioIDs []string
	multiChar := false
	for _, v := range voices {
		if !seen[v.AudioID] {
			seen[v.AudioID] = true
			audioIDs = append(audioIDs, v.AudioID)
			if v.CharacterID != characterID {
				multiChar = true
			}
		}
	}

	var audioCharMap map[string]string
	if multiChar {
		audioCharMap = make(map[string]string, len(audioIDs))
		for _, v := range voices {
			if audioCharMap[v.AudioID] == "" {
				audioCharMap[v.AudioID] = v.CharacterID
			}
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
		Content:      d.Content,
		CharacterID:  characterID,
		AudioID:      strings.Join(audioIDs, ", "),
		AudioCharMap: audioCharMap,
		Episode:      episode,
		Truth:        truth,
	}
}

// hasWordsBeforeVoice walks dialogue elements in document order and returns
// true if actual word characters (letters) appear before the first voice
// command. This detects narration lines that have embedded voice clips at the
// end (e.g. omake commentary ending with "Mii, nipah~☆").
// Lines where only dots/punctuation appear before the voice command (character
// pauses like "............") are not treated as narration.
func hasWordsBeforeVoice(elements []ast.DialogueElement) bool {
	for _, elem := range elements {
		switch el := elem.(type) {
		case *ast.VoiceCommand:
			return false
		case *ast.PlainText:
			if containsLetters(el.Text) {
				return true
			}
		case *ast.FormatTag:
			if result := hasWordsBeforeVoice(el.Content); result {
				return true
			}
		}
	}
	return false
}

func containsLetters(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

// Presets returns the preset context for use with transformers.
func (e *QuoteExtractor) Presets() *transformer.PresetContext {
	return e.presets
}
