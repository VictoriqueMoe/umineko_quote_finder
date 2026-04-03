package subtitle

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"umineko_quote/internal/dto"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/umineko/character"
	umilexer "github.com/VictoriqueMoe/umineko_script_parser/umineko/lexer"
)

var styleCharacter = map[string]string{
	"Battler": "10",
}

func ResolveRefs(efs fs.ReadFileFS, refs []umilexer.SubtitleRef) []dto.UminekoQuote {
	var quotes []dto.UminekoQuote

	for _, ref := range refs {
		filename := filepath.Base(strings.ReplaceAll(ref.SubPath, `\`, "/"))
		data, err := efs.ReadFile("data/sub/" + filename)
		if err != nil {
			log.Printf("[subtitle] could not read %s: %v", filename, err)
			continue
		}

		entries := ParseASS(data)
		for i, entry := range entries {
			charID := ref.CharacterID
			if mapped, ok := styleCharacter[entry.Style]; ok {
				charID = mapped
			}

			quotes = append(quotes, dto.UminekoQuote{
				UminekoQuote: scriptdto.UminekoQuote{
					BaseQuote: scriptdto.BaseQuote{
						Text:        entry.Text,
						TextHtml:    entry.Text,
						CharacterID: charID,
						Character:   character.CharacterNames.GetCharacterName(character.CharacterFromID(charID)),
						AudioID:     fmt.Sprintf("%s_s%d", ref.AudioID, i),
						Episode:     ref.Episode,
						ContentType: ref.ContentType,
					},
				},
			})
		}
	}

	return quotes
}
