package ast

import "fmt"

type (
	TokenType int

	Token struct {
		Type   TokenType
		Value  string
		Line   int
		Column int
	}
)

const (
	TokenEOF TokenType = iota
	TokenNewline
	TokenComment
	TokenLabel
	TokenCommand
	TokenNumber
	TokenString
	TokenComma
	TokenBacktick
	TokenText
	TokenBracketOpen
	TokenBracketClose
	TokenBraceOpen
	TokenBraceClose
	TokenColon
	TokenAsterisk
	TokenPipe
	TokenAt
	TokenBackslash
	TokenHash
	TokenExclaim
	TokenInlineCommand
	TokenFormatTag
)

func (t TokenType) String() string {
	names := map[TokenType]string{
		TokenEOF:           "EOF",
		TokenNewline:       "Newline",
		TokenComment:       "Comment",
		TokenLabel:         "Label",
		TokenCommand:       "Command",
		TokenNumber:        "Number",
		TokenString:        "String",
		TokenComma:         "Comma",
		TokenBacktick:      "Backtick",
		TokenText:          "Text",
		TokenBracketOpen:   "BracketOpen",
		TokenBracketClose:  "BracketClose",
		TokenBraceOpen:     "BraceOpen",
		TokenBraceClose:    "BraceClose",
		TokenColon:         "Colon",
		TokenAsterisk:      "Asterisk",
		TokenPipe:          "Pipe",
		TokenAt:            "At",
		TokenBackslash:     "Backslash",
		TokenHash:          "Hash",
		TokenExclaim:       "Exclaim",
		TokenInlineCommand: "InlineCommand",
		TokenFormatTag:     "FormatTag",
	}
	if name, ok := names[t]; ok {
		return name
	}
	return "Unknown"
}

func (t Token) Position() string {
	return fmt.Sprintf("%d:%d", t.Line, t.Column)
}

// AST node interfaces

type (
	Node interface {
		nodeType() string
	}

	Line interface {
		Node
		lineNode()
	}

	DialogueElement interface {
		Node
		dialogueElement()
	}
)

// AST node types

type (
	Script struct {
		Lines []Line
	}

	CommentLine struct {
		Text string
		Pos  Token
	}

	LabelLine struct {
		Name string
		Pos  Token
	}

	CommandLine struct {
		Command string
		Args    []Token
		Pos     Token
	}

	PresetDefineLine struct {
		ID     int
		FontID int
		Size   int
		Colour string
		Flags  []int
		Pos    Token
	}

	EpisodeMarkerLine struct {
		Type    string
		Episode int
		Pos     Token
	}

	DialogueLine struct {
		Command string
		Content []DialogueElement
		Pos     Token
	}

	PlainText struct {
		Text string
		Pos  Token
	}

	FormatTag struct {
		Name    string
		Param   string
		Content []DialogueElement
		Pos     Token
	}

	SpecialChar struct {
		Name string
		Pos  Token
	}

	InlineCommand struct {
		Command string
		Args    string
		Pos     Token
	}

	VoiceCommand struct {
		Channel     int
		CharacterID string
		AudioID     string
		Pos         Token
	}

	ClickWait struct {
		Type string
		Pos  Token
	}

	TimedWait struct {
		Skippable bool
		Duration  int
		Pos       Token
	}
)

// Marker methods to restrict interface implementations.

func (s *Script) nodeType() string            { return "Script" }
func (c *CommentLine) nodeType() string       { return "CommentLine" }
func (c *CommentLine) lineNode()              {}
func (l *LabelLine) nodeType() string         { return "LabelLine" }
func (l *LabelLine) lineNode()                {}
func (c *CommandLine) nodeType() string       { return "CommandLine" }
func (c *CommandLine) lineNode()              {}
func (p *PresetDefineLine) nodeType() string  { return "PresetDefineLine" }
func (p *PresetDefineLine) lineNode()         {}
func (e *EpisodeMarkerLine) nodeType() string { return "EpisodeMarkerLine" }
func (e *EpisodeMarkerLine) lineNode()        {}
func (d *DialogueLine) nodeType() string      { return "DialogueLine" }
func (d *DialogueLine) lineNode()             {}
func (p *PlainText) nodeType() string         { return "PlainText" }
func (p *PlainText) dialogueElement()         {}
func (f *FormatTag) nodeType() string         { return "FormatTag" }
func (f *FormatTag) dialogueElement()         {}
func (s *SpecialChar) nodeType() string       { return "SpecialChar" }
func (s *SpecialChar) dialogueElement()       {}
func (i *InlineCommand) nodeType() string     { return "InlineCommand" }
func (i *InlineCommand) dialogueElement()     {}
func (v *VoiceCommand) nodeType() string      { return "VoiceCommand" }
func (v *VoiceCommand) dialogueElement()      {}
func (c *ClickWait) nodeType() string         { return "ClickWait" }
func (c *ClickWait) dialogueElement()         {}
func (t *TimedWait) nodeType() string         { return "TimedWait" }
func (t *TimedWait) dialogueElement()         {}

// Helper methods

func (d *DialogueLine) GetVoiceCommands() []*VoiceCommand {
	var voices []*VoiceCommand
	d.collectVoiceCommands(d.Content, &voices)
	return voices
}

func (d *DialogueLine) collectVoiceCommands(elements []DialogueElement, voices *[]*VoiceCommand) {
	for _, elem := range elements {
		switch el := elem.(type) {
		case *VoiceCommand:
			*voices = append(*voices, el)
		case *FormatTag:
			d.collectVoiceCommands(el.Content, voices)
		}
	}
}
