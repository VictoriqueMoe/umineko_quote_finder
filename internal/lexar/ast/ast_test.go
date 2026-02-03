package ast

import "testing"

func TestTokenType_String(t *testing.T) {
	tests := []struct {
		typ  TokenType
		want string
	}{
		{TokenEOF, "EOF"},
		{TokenNewline, "Newline"},
		{TokenComment, "Comment"},
		{TokenLabel, "Label"},
		{TokenCommand, "Command"},
		{TokenNumber, "Number"},
		{TokenString, "String"},
		{TokenComma, "Comma"},
		{TokenBacktick, "Backtick"},
		{TokenText, "Text"},
		{TokenBracketOpen, "BracketOpen"},
		{TokenBracketClose, "BracketClose"},
		{TokenBraceOpen, "BraceOpen"},
		{TokenBraceClose, "BraceClose"},
		{TokenColon, "Colon"},
		{TokenAsterisk, "Asterisk"},
		{TokenPipe, "Pipe"},
		{TokenAt, "At"},
		{TokenBackslash, "Backslash"},
		{TokenHash, "Hash"},
		{TokenExclaim, "Exclaim"},
		{TokenInlineCommand, "InlineCommand"},
		{TokenFormatTag, "FormatTag"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.typ.String(); got != tt.want {
				t.Errorf("TokenType(%d).String() = %q, want %q", tt.typ, got, tt.want)
			}
		})
	}
}

func TestTokenType_String_Unknown(t *testing.T) {
	unknown := TokenType(999)
	if got := unknown.String(); got != "Unknown" {
		t.Errorf("TokenType(999).String() = %q, want %q", got, "Unknown")
	}
}

func TestToken_Position(t *testing.T) {
	tests := []struct {
		name  string
		token Token
		want  string
	}{
		{"origin", Token{Line: 1, Column: 1}, "1:1"},
		{"mid-file", Token{Line: 42, Column: 13}, "42:13"},
		{"large", Token{Line: 10000, Column: 500}, "10000:500"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.token.Position(); got != tt.want {
				t.Errorf("Position() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNodeTypes(t *testing.T) {
	tests := []struct {
		name string
		node Node
		want string
	}{
		{"Script", &Script{}, "Script"},
		{"CommentLine", &CommentLine{}, "CommentLine"},
		{"LabelLine", &LabelLine{}, "LabelLine"},
		{"CommandLine", &CommandLine{}, "CommandLine"},
		{"PresetDefineLine", &PresetDefineLine{}, "PresetDefineLine"},
		{"EpisodeMarkerLine", &EpisodeMarkerLine{}, "EpisodeMarkerLine"},
		{"DialogueLine", &DialogueLine{}, "DialogueLine"},
		{"PlainText", &PlainText{}, "PlainText"},
		{"FormatTag", &FormatTag{}, "FormatTag"},
		{"SpecialChar", &SpecialChar{}, "SpecialChar"},
		{"InlineCommand", &InlineCommand{}, "InlineCommand"},
		{"VoiceCommand", &VoiceCommand{}, "VoiceCommand"},
		{"ClickWait", &ClickWait{}, "ClickWait"},
		{"TimedWait", &TimedWait{}, "TimedWait"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.node.nodeType(); got != tt.want {
				t.Errorf("nodeType() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLineInterface(t *testing.T) {
	// Verify all Line types satisfy the interface.
	lines := []Line{
		&CommentLine{},
		&LabelLine{},
		&CommandLine{},
		&PresetDefineLine{},
		&EpisodeMarkerLine{},
		&DialogueLine{},
	}

	for _, l := range lines {
		l.lineNode()
		if l.nodeType() == "" {
			t.Errorf("Line %T has empty nodeType", l)
		}
	}
}

func TestDialogueElementInterface(t *testing.T) {
	// Verify all DialogueElement types satisfy the interface.
	elements := []DialogueElement{
		&PlainText{},
		&FormatTag{},
		&SpecialChar{},
		&InlineCommand{},
		&VoiceCommand{},
		&ClickWait{},
		&TimedWait{},
	}

	for _, e := range elements {
		e.dialogueElement()
		if e.nodeType() == "" {
			t.Errorf("DialogueElement %T has empty nodeType", e)
		}
	}
}

func TestGetVoiceCommands_Flat(t *testing.T) {
	v1 := &VoiceCommand{CharacterID: "10", AudioID: "10100001"}
	v2 := &VoiceCommand{CharacterID: "10", AudioID: "10100002"}

	d := &DialogueLine{
		Content: []DialogueElement{
			v1,
			&PlainText{Text: "Hello."},
			v2,
			&PlainText{Text: "World."},
		},
	}

	voices := d.GetVoiceCommands()
	if len(voices) != 2 {
		t.Fatalf("got %d voice commands, want 2", len(voices))
	}
	if voices[0].AudioID != "10100001" {
		t.Errorf("voices[0].AudioID = %q, want %q", voices[0].AudioID, "10100001")
	}
	if voices[1].AudioID != "10100002" {
		t.Errorf("voices[1].AudioID = %q, want %q", voices[1].AudioID, "10100002")
	}
}

func TestGetVoiceCommands_NestedInFormatTag(t *testing.T) {
	v := &VoiceCommand{CharacterID: "27", AudioID: "40700001"}

	d := &DialogueLine{
		Content: []DialogueElement{
			&FormatTag{
				Name:  "p",
				Param: "1",
				Content: []DialogueElement{
					v,
					&PlainText{Text: "Red truth!"},
				},
			},
		},
	}

	voices := d.GetVoiceCommands()
	if len(voices) != 1 {
		t.Fatalf("got %d voice commands, want 1", len(voices))
	}
	if voices[0].AudioID != "40700001" {
		t.Errorf("AudioID = %q, want %q", voices[0].AudioID, "40700001")
	}
}

func TestGetVoiceCommands_DeeplyNested(t *testing.T) {
	v := &VoiceCommand{CharacterID: "10", AudioID: "50100001"}

	d := &DialogueLine{
		Content: []DialogueElement{
			&FormatTag{
				Name: "p",
				Content: []DialogueElement{
					&FormatTag{
						Name: "c",
						Content: []DialogueElement{
							v,
							&PlainText{Text: "Deep."},
						},
					},
				},
			},
		},
	}

	voices := d.GetVoiceCommands()
	if len(voices) != 1 {
		t.Fatalf("got %d voice commands, want 1", len(voices))
	}
	if voices[0].CharacterID != "10" {
		t.Errorf("CharacterID = %q, want %q", voices[0].CharacterID, "10")
	}
}

func TestGetVoiceCommands_Empty(t *testing.T) {
	d := &DialogueLine{
		Content: []DialogueElement{
			&PlainText{Text: "No voice here."},
		},
	}

	voices := d.GetVoiceCommands()
	if len(voices) != 0 {
		t.Errorf("got %d voice commands, want 0", len(voices))
	}
}

func TestGetVoiceCommands_NilContent(t *testing.T) {
	d := &DialogueLine{}

	voices := d.GetVoiceCommands()
	if len(voices) != 0 {
		t.Errorf("got %d voice commands, want 0", len(voices))
	}
}

func TestGetVoiceCommands_MixedElements(t *testing.T) {
	v1 := &VoiceCommand{CharacterID: "27", AudioID: "10700001"}
	v2 := &VoiceCommand{CharacterID: "27", AudioID: "10700002"}

	d := &DialogueLine{
		Content: []DialogueElement{
			v1,
			&ClickWait{Type: "@"},
			&PlainText{Text: "text"},
			&SpecialChar{Name: "n"},
			&FormatTag{
				Name: "i",
				Content: []DialogueElement{
					v2,
					&PlainText{Text: "italic"},
				},
			},
			&TimedWait{Duration: 100},
			&InlineCommand{Command: "foo"},
		},
	}

	voices := d.GetVoiceCommands()
	if len(voices) != 2 {
		t.Fatalf("got %d voice commands, want 2", len(voices))
	}
	if voices[0].AudioID != "10700001" {
		t.Errorf("voices[0].AudioID = %q, want %q", voices[0].AudioID, "10700001")
	}
	if voices[1].AudioID != "10700002" {
		t.Errorf("voices[1].AudioID = %q, want %q", voices[1].AudioID, "10700002")
	}
}
