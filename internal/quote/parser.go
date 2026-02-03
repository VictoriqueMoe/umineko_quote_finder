package quote

import (
	"html"
	"regexp"
	"strconv"
	"strings"
)

type (
	Parser interface {
		ParseAll(lines []string) []ParsedQuote
	}

	ParsedQuote struct {
		Text        string `json:"text"`
		TextHtml    string `json:"textHtml"`
		CharacterID string `json:"characterId"`
		Character   string `json:"character"`
		AudioID     string `json:"audioId"`
		Episode     int    `json:"episode"`
		ContentType string `json:"contentType"`
	}

	textRule struct {
		pattern   *regexp.Regexp
		htmlRepl  string
		plainRepl string
	}

	parser struct {
		dialogueLineRegex *regexp.Regexp
		narratorLineRegex *regexp.Regexp
		voiceMetaRegex    *regexp.Regexp
		bracketRegex      *regexp.Regexp
		contentTypeRegex  *regexp.Regexp // matches new_episode, new_tea, new_ura
		omakeRegex        *regexp.Regexp
		presetRegex       *regexp.Regexp // matches {p:N:content}
		unclosedTagRegex  *regexp.Regexp
		cleanupPatterns   []string
		textRules         []textRule
		presetColours     map[string]string // parsed from script's preset_define lines
	}
)

var (
	specialCharTags = []struct {
		tag         string
		replacement string
	}{
		{"{0}", ""},
		{"{-}", ""},
		{"{qt}", `"`},
		{"{ob}", "{"},
		{"{eb}", "}"},
		{"{os}", "["},
		{"{es}", "]"},
		{"{t}", ""},
		{"{parallel}", ""},
	}

	// presetClasses maps preset IDs to semantic CSS class names (these override colours)
	presetClasses = map[string]string{
		"1": "red-truth",
		"2": "blue-truth",
	}

	// nestedContent matches content that may contain nested braces (for truth tags)
	nestedContent = `((?:[^{}]|\{[^{}]*\})+)`

	// simpleContent matches content without nested braces
	simpleContent = `([^{}]+)`

	presetDefineRegex = regexp.MustCompile(`^preset_define (\d+),\d+,-?\d+,(#[A-Fa-f0-9]{6})`)
)

func NewParser() Parser {
	return &parser{
		dialogueLineRegex: regexp.MustCompile(`^d2? (?:\[[^\]]*\])*\[lv`),
		voiceMetaRegex:    regexp.MustCompile(`\[lv \d+\*"(\d+)"\*"(.+?)"\]`),
		narratorLineRegex: regexp.MustCompile("^d2? `"),
		bracketRegex:      regexp.MustCompile(`\[[^\]]*\]`),
		contentTypeRegex:  regexp.MustCompile(`^new_(episode|tea|ura) (\d+)\r?$`),
		omakeRegex:        regexp.MustCompile(`^\*o(\d+)_`),
		presetRegex:       regexp.MustCompile(`\{p:(\d+):(` + nestedContent + `)\}?`),
		unclosedTagRegex:  regexp.MustCompile(`\{[a-zA-Z]+:(?:[^{}:]*:)*`),
		cleanupPatterns: []string{
			"`[@]", "`[\\]", "`[|]", "`\"", "\"`",
			"[@]", "[\\]", "[|]",
		},
		textRules: []textRule{
			// Special tags
			{regexp.MustCompile(`\{n\}`), "<br>", " "},
			{regexp.MustCompile(`\{i:` + simpleContent + `\}`), `<em>$1</em>`, "$1"},
			{regexp.MustCompile(`\{c:([A-Fa-f0-9]+):` + simpleContent + `\}`), `<span style="color:#$1">$2</span>`, "$2"},
			{regexp.MustCompile(`\{ruby:([^:]+):` + simpleContent + `\}`), `<ruby>$2<rp>(</rp><rt>$1</rt><rp>)</rp></ruby>`, "$2 ($1)"},
			// Tags that strip content entirely
			{regexp.MustCompile(`\{y:\d+:([^{}]*)\}`), "", ""},
			// Passthrough tags (content preserved, tag removed)
			{regexp.MustCompile(`\{(?:f|n|a|nobr):[^{}:]*:([^{}]*)\}`), "$1", "$1"},
			{regexp.MustCompile(`\{[a-zA-Z]+:(?:[^{}:]*:)?([^{}]*)\}`), "$1", "$1"},
		},
	}
}

func (p *parser) parseLine(line string) *ParsedQuote {
	if !p.dialogueLineRegex.MatchString(line) {
		return nil
	}

	allMatches := p.voiceMetaRegex.FindAllStringSubmatch(line, -1)
	if len(allMatches) == 0 || len(allMatches[0]) < 3 {
		return nil
	}

	characterID := allMatches[0][1]
	firstAudioID := allMatches[0][2]
	episode := p.parseEpisodeFromAudioID(firstAudioID)

	seen := map[string]bool{}
	var audioIDs []string
	for i := 0; i < len(allMatches); i++ {
		id := allMatches[i][2]
		if !seen[id] {
			seen[id] = true
			audioIDs = append(audioIDs, id)
		}
	}

	text, textHtml := p.extractText(line)
	if text == "" {
		return nil
	}

	return &ParsedQuote{
		Text:        text,
		TextHtml:    textHtml,
		CharacterID: characterID,
		Character:   CharacterNames.GetCharacterName(characterID),
		AudioID:     strings.Join(audioIDs, ", "),
		Episode:     episode,
	}
}

func (p *parser) parseEpisodeFromAudioID(audioID string) int {
	if len(audioID) < 1 {
		return 0
	}
	ep := int(audioID[0] - '0')
	if ep >= 1 && ep <= 8 {
		return ep
	}
	return 0
}

func (p *parser) replacePresets(text string, forHtml bool) string {
	return p.presetRegex.ReplaceAllStringFunc(text, func(match string) string {
		parts := p.presetRegex.FindStringSubmatch(match)
		if len(parts) < 3 {
			return match
		}
		presetID, content := parts[1], parts[2]

		if !forHtml {
			return content
		}

		if class, ok := presetClasses[presetID]; ok {
			return `<span class="` + class + `">` + content + `</span>`
		}

		if colour, ok := p.presetColours[presetID]; ok {
			return `<span style="color:` + colour + `">` + content + `</span>`
		}

		return content
	})
}

func (p *parser) parsePresetColours(lines []string) {
	p.presetColours = make(map[string]string)

	for _, line := range lines {
		if matches := presetDefineRegex.FindStringSubmatch(line); len(matches) >= 3 {
			presetID := matches[1]
			colour := strings.ToUpper(matches[2])
			if _, hasClass := presetClasses[presetID]; hasClass {
				continue
			}
			if colour == "#FFFFFF" {
				continue
			}
			p.presetColours[presetID] = colour
		}
	}
}

func (p *parser) extractText(line string) (string, string) {
	text := line

	for _, pattern := range p.cleanupPatterns {
		text = strings.ReplaceAll(text, pattern, "")
	}

	text = p.bracketRegex.ReplaceAllString(text, "")
	text = strings.TrimPrefix(text, "d2 ")
	text = strings.TrimPrefix(text, "d ")
	text = strings.TrimSpace(text)
	text = strings.Trim(text, "`\"")
	text = strings.ReplaceAll(text, "`", "")
	text = strings.TrimSpace(text)

	plainText := text
	textHtml := html.EscapeString(text)

	for _, sc := range specialCharTags {
		plainText = strings.ReplaceAll(plainText, sc.tag, sc.replacement)
		textHtml = strings.ReplaceAll(textHtml, sc.tag, sc.replacement)
	}

	for {
		prevHtml := textHtml
		prevPlain := plainText

		// Handle preset tags with custom logic
		textHtml = p.replacePresets(textHtml, true)
		plainText = p.replacePresets(plainText, false)

		// Handle other text rules
		for _, rule := range p.textRules {
			textHtml = rule.pattern.ReplaceAllString(textHtml, rule.htmlRepl)
			plainText = rule.pattern.ReplaceAllString(plainText, rule.plainRepl)
		}
		if textHtml == prevHtml && plainText == prevPlain {
			break
		}
	}

	textHtml = p.unclosedTagRegex.ReplaceAllString(textHtml, "")
	plainText = p.unclosedTagRegex.ReplaceAllString(plainText, "")

	textHtml = strings.ReplaceAll(textHtml, "{", "")
	textHtml = strings.ReplaceAll(textHtml, "}", "")
	plainText = strings.ReplaceAll(plainText, "{", "")
	plainText = strings.ReplaceAll(plainText, "}", "")

	return plainText, textHtml
}

func (p *parser) parseNarratorLine(line string) *ParsedQuote {
	if !p.narratorLineRegex.MatchString(line) {
		return nil
	}

	text, textHtml := p.extractText(line)
	if text == "" {
		return nil
	}

	characterID := "narrator"
	character := CharacterNames.GetCharacterName("narrator")
	audioID := ""
	episode := 0

	allMatches := p.voiceMetaRegex.FindAllStringSubmatch(line, -1)
	if len(allMatches) > 0 && len(allMatches[0]) >= 3 && strings.Contains(line, "{a:") {
		characterID = allMatches[0][1]
		character = CharacterNames.GetCharacterName(characterID)
		firstAudioID := allMatches[0][2]
		episode = p.parseEpisodeFromAudioID(firstAudioID)

		seen := map[string]bool{}
		var audioIDs []string
		for i := 0; i < len(allMatches); i++ {
			id := allMatches[i][2]
			if !seen[id] {
				seen[id] = true
				audioIDs = append(audioIDs, id)
			}
		}
		audioID = strings.Join(audioIDs, ", ")
	}

	return &ParsedQuote{
		Text:        text,
		TextHtml:    textHtml,
		CharacterID: characterID,
		Character:   character,
		AudioID:     audioID,
		Episode:     episode,
	}
}

func (p *parser) ParseAll(lines []string) []ParsedQuote {
	p.parsePresetColours(lines)

	quotes := make([]ParsedQuote, 0, len(lines)/6)
	currentEpisode := 0
	currentContentType := ""

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if matches := p.contentTypeRegex.FindStringSubmatch(line); len(matches) >= 3 {
			ep, err := strconv.Atoi(matches[2])
			if err == nil {
				currentEpisode = ep
				if matches[1] == "episode" {
					currentContentType = ""
				} else {
					currentContentType = matches[1]
				}
			}
			continue
		}
		if matches := p.omakeRegex.FindStringSubmatch(line); len(matches) >= 2 {
			ep, err := strconv.Atoi(matches[1])
			if err == nil {
				currentEpisode = ep
				currentContentType = "omake"
			}
			continue
		}

		parsed := p.parseLine(line)
		if parsed == nil {
			parsed = p.parseNarratorLine(line)
		}
		if parsed == nil || len(parsed.Text) <= 10 {
			continue
		}
		if currentEpisode > 0 {
			parsed.Episode = currentEpisode
		}
		parsed.ContentType = currentContentType
		quotes = append(quotes, *parsed)
	}

	return quotes
}
