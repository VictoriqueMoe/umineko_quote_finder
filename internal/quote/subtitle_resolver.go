package quote

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/subtitle"

	"github.com/VictoriqueMoe/umineko_script_parser/lexer"
	"github.com/VictoriqueMoe/umineko_script_parser/quote/character"
)

var subtitleStyleCharacter = map[string]string{
	"Battler": "10",
}

func resolveSubtitleRefs(efs fs.ReadFileFS, refs []lexer.SubtitleRef) []dto.ParsedQuote {
	var quotes []dto.ParsedQuote

	for _, ref := range refs {
		filename := filepath.Base(strings.ReplaceAll(ref.SubPath, `\`, "/"))
		data, err := efs.ReadFile("data/sub/" + filename)
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
