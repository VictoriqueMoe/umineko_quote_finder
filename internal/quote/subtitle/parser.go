package subtitle

import (
	"regexp"
	"strings"
)

type Entry struct {
	Style string
	Text  string
}

var assOverrideBlock = regexp.MustCompile(`\{[^}]*\}`)

func ParseASS(data []byte) []Entry {
	lines := strings.Split(string(data), "\n")

	inEvents := false
	textIdx := -1
	styleIdx := -1

	var entries []Entry

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")

		if strings.HasPrefix(line, "[Events]") {
			inEvents = true
			continue
		}

		if strings.HasPrefix(line, "[") && inEvents {
			break
		}

		if !inEvents {
			continue
		}

		if strings.HasPrefix(line, "Format:") {
			fields := strings.Split(line[len("Format:"):], ",")
			for i, f := range fields {
				f = strings.TrimSpace(f)
				if f == "Text" {
					textIdx = i
				}
				if f == "Style" {
					styleIdx = i
				}
			}
			continue
		}

		if !strings.HasPrefix(line, "Dialogue:") {
			continue
		}

		if textIdx < 0 {
			continue
		}

		content := line[len("Dialogue:"):]
		parts := strings.SplitN(content, ",", textIdx+1)
		if len(parts) <= textIdx {
			continue
		}

		rawText := parts[textIdx]
		style := ""
		if styleIdx >= 0 && styleIdx < len(parts) {
			style = strings.TrimSpace(parts[styleIdx])
		}

		cleaned := stripASSTags(rawText)
		cleaned = strings.TrimSpace(cleaned)

		if cleaned == "" {
			continue
		}

		entries = append(entries, Entry{
			Style: style,
			Text:  cleaned,
		})
	}

	return entries
}

func stripASSTags(s string) string {
	s = assOverrideBlock.ReplaceAllString(s, "")
	s = strings.ReplaceAll(s, `\h`, " ")
	s = strings.ReplaceAll(s, `\N`, " ")
	s = strings.ReplaceAll(s, `\n`, " ")
	return s
}
