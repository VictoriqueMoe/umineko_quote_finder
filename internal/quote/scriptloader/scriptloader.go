package scriptloader

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"time"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/lexar"
	"umineko_quote/internal/quote/character"
	"umineko_quote/internal/quote/mutation"
	"umineko_quote/internal/subtitle"
)

type (
	ParseFunc func(lines []string) ([]dto.ParsedQuote, []lexar.SubtitleRef)

	Loader struct {
		fs        fs.ReadFileFS
		parse     ParseFunc
		mutations mutation.Pipeline
	}
)

var subtitleStyleCharacter = map[string]string{
	"Battler": "10",
}

func New(efs fs.ReadFileFS, parse ParseFunc) *Loader {
	return &Loader{
		fs:        efs,
		parse:     parse,
		mutations: *mutation.NewPipeline(),
	}
}

func (l *Loader) Load(lang string, path string) []dto.ParsedQuote {
	raw, err := l.fs.ReadFile(path)
	if err != nil {
		log.Printf("[%s] failed to read %s: %v", lang, path, err)
		return nil
	}

	decodeStart := time.Now()
	decoded, err := decode(raw)
	if err != nil {
		log.Printf("[%s] failed to decode %s: %v", lang, path, err)
		return nil
	}
	log.Printf("[%s] decoded %s (%d → %d bytes) in %v", lang, path, len(raw), len(decoded), time.Since(decodeStart).Round(time.Millisecond))

	lines := strings.Split(string(decoded), "\n")

	parseStart := time.Now()
	parsed, subtitleRefs := l.parse(lines)
	log.Printf("[%s] parsed %d lines → %d quotes in %v", lang, len(lines), len(parsed), time.Since(parseStart).Round(time.Millisecond))

	subQuotes := l.resolveSubtitleRefs(subtitleRefs)
	if len(subQuotes) > 0 {
		parsed = append(parsed, subQuotes...)
		log.Printf("[%s] added %d subtitle quotes", lang, len(subQuotes))
	}

	parsed = l.mutations.Apply(parsed)

	return parsed
}

func (l *Loader) resolveSubtitleRefs(refs []lexar.SubtitleRef) []dto.ParsedQuote {
	var quotes []dto.ParsedQuote

	for _, ref := range refs {
		filename := filepath.Base(strings.ReplaceAll(ref.SubPath, `\`, "/"))
		data, err := l.fs.ReadFile("data/sub/" + filename)
		if err != nil {
			log.Printf("[subtitle] could not read %s: %v", filename, err)
			continue
		}

		entries := subtitle.ParseASS(data)
		for i, entry := range entries {
			charID := ref.CharacterID
			if mapped, ok := subtitleStyleCharacter[entry.Style]; ok {
				charID = mapped
			}

			quotes = append(quotes, dto.ParsedQuote{
				Text:        entry.Text,
				TextHtml:    entry.Text,
				CharacterID: charID,
				Character:   character.CharacterNames.GetCharacterName(character.CharacterFromID(charID)),
				AudioID:     fmt.Sprintf("%s_s%d", ref.AudioID, i),
				Episode:     ref.Episode,
				ContentType: ref.ContentType,
			})
		}
	}

	return quotes
}
