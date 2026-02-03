package quote

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
		TruthType   string `json:"truthType,omitempty"`
	}
)

// NewParser creates a new parser using the script-based lexer/parser.
func NewParser() Parser {
	return NewScriptParser()
}
