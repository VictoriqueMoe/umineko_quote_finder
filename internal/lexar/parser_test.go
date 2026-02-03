package lexar

import (
	"testing"

	"umineko_quote/internal/lexar/ast"
	"umineko_quote/internal/lexar/transformer"
)

var plainTextTransformer = transformer.NewPlainTextTransformer()

func getPlainText(d *ast.DialogueLine) string {
	return plainTextTransformer.Transform(d.Content)
}

func TestParse_DialogueLine(t *testing.T) {
	input := `d ` + "`Hello, world!`"
	script := Parse(input)

	if len(script.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(script.Lines))
	}

	d, ok := script.Lines[0].(*ast.DialogueLine)
	if !ok {
		t.Fatalf("expected DialogueLine, got %T", script.Lines[0])
	}

	if d.Command != "d" {
		t.Errorf("command: got %q, want 'd'", d.Command)
	}

	text := getPlainText(d)
	if text != "Hello, world!" {
		t.Errorf("plain text: got %q, want 'Hello, world!'", text)
	}
}

func TestParse_DialogueWithVoice(t *testing.T) {
	input := `d2 [lv 0*"10"*"10100001"]` + "`Test line`"
	script := Parse(input)

	if len(script.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(script.Lines))
	}

	d, ok := script.Lines[0].(*ast.DialogueLine)
	if !ok {
		t.Fatalf("expected DialogueLine, got %T", script.Lines[0])
	}

	voices := d.GetVoiceCommands()
	if len(voices) != 1 {
		t.Fatalf("expected 1 voice command, got %d", len(voices))
	}

	v := voices[0]
	if v.Channel != 0 {
		t.Errorf("channel: got %d, want 0", v.Channel)
	}
	if v.CharacterID != "10" {
		t.Errorf("character ID: got %q, want '10'", v.CharacterID)
	}
	if v.AudioID != "10100001" {
		t.Errorf("audio ID: got %q, want '10100001'", v.AudioID)
	}
}

func TestParse_DialogueWithFormatTag(t *testing.T) {
	input := `d ` + "`This is {c:FF0000:red text} here.`"
	script := Parse(input)

	d := script.Lines[0].(*ast.DialogueLine)
	text := getPlainText(d)

	if text != "This is red text here." {
		t.Errorf("plain text: got %q, want 'This is red text here.'", text)
	}

	// Find the format tag
	var foundTag *ast.FormatTag
	for _, elem := range d.Content {
		if ft, ok := elem.(*ast.FormatTag); ok {
			foundTag = ft
			break
		}
	}

	if foundTag == nil {
		t.Fatal("expected to find FormatTag")
	}
	if foundTag.Name != "c" {
		t.Errorf("tag name: got %q, want 'c'", foundTag.Name)
	}
	if foundTag.Param != "FF0000" {
		t.Errorf("tag param: got %q, want 'FF0000'", foundTag.Param)
	}
}

func TestParse_DialogueWithPreset(t *testing.T) {
	input := `d ` + "`{p:1:Red truth text}`"
	script := Parse(input)

	d := script.Lines[0].(*ast.DialogueLine)

	var foundTag *ast.FormatTag
	for _, elem := range d.Content {
		if ft, ok := elem.(*ast.FormatTag); ok {
			foundTag = ft
			break
		}
	}

	if foundTag == nil {
		t.Fatal("expected to find FormatTag")
	}
	if foundTag.Name != "p" {
		t.Errorf("tag name: got %q, want 'p'", foundTag.Name)
	}
	if foundTag.Param != "1" {
		t.Errorf("tag param: got %q, want '1'", foundTag.Param)
	}

	text := getPlainText(d)
	if text != "Red truth text" {
		t.Errorf("plain text: got %q, want 'Red truth text'", text)
	}
}

func TestParse_DialogueWithSpecialChars(t *testing.T) {
	input := `d ` + "`Line one{n}Line two`"
	script := Parse(input)

	d := script.Lines[0].(*ast.DialogueLine)
	text := getPlainText(d)

	if text != "Line one Line two" {
		t.Errorf("plain text: got %q, want 'Line one Line two'", text)
	}
}

func TestParse_PresetDefine(t *testing.T) {
	input := `preset_define 1,1,-1,#FF0000,0,0,0,0,0`
	script := Parse(input)

	if len(script.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(script.Lines))
	}

	p, ok := script.Lines[0].(*ast.PresetDefineLine)
	if !ok {
		t.Fatalf("expected PresetDefineLine, got %T", script.Lines[0])
	}

	if p.ID != 1 {
		t.Errorf("ID: got %d, want 1", p.ID)
	}
	if p.FontID != 1 {
		t.Errorf("FontID: got %d, want 1", p.FontID)
	}
	if p.Size != -1 {
		t.Errorf("Size: got %d, want -1", p.Size)
	}
	if p.Colour != "#FF0000" {
		t.Errorf("Colour: got %q, want '#FF0000'", p.Colour)
	}
}

func TestParse_EpisodeMarker(t *testing.T) {
	tests := []struct {
		input   string
		epType  string
		episode int
	}{
		{"new_episode 1", "episode", 1},
		{"new_tea 5", "tea", 5},
		{"new_ura 3", "ura", 3},
	}

	for _, tt := range tests {
		script := Parse(tt.input)
		if len(script.Lines) != 1 {
			t.Fatalf("%q: expected 1 line, got %d", tt.input, len(script.Lines))
		}

		m, ok := script.Lines[0].(*ast.EpisodeMarkerLine)
		if !ok {
			t.Fatalf("%q: expected EpisodeMarkerLine, got %T", tt.input, script.Lines[0])
		}

		if m.Type != tt.epType {
			t.Errorf("%q: type got %q, want %q", tt.input, m.Type, tt.epType)
		}
		if m.Episode != tt.episode {
			t.Errorf("%q: episode got %d, want %d", tt.input, m.Episode, tt.episode)
		}
	}
}

func TestParse_ClickWaits(t *testing.T) {
	input := `d ` + "`Text`[@]`More`[\\]"
	script := Parse(input)

	d := script.Lines[0].(*ast.DialogueLine)

	var clicks []*ast.ClickWait
	for _, elem := range d.Content {
		if cw, ok := elem.(*ast.ClickWait); ok {
			clicks = append(clicks, cw)
		}
	}

	if len(clicks) != 2 {
		t.Fatalf("expected 2 click waits, got %d", len(clicks))
	}
	if clicks[0].Type != "@" {
		t.Errorf("first click wait: got %q, want '@'", clicks[0].Type)
	}
	if clicks[1].Type != "\\" {
		t.Errorf("second click wait: got %q, want '\\\\'", clicks[1].Type)
	}
}

func TestParse_TimedWait(t *testing.T) {
	input := `d ` + "`Text`[!w500]`More`"
	script := Parse(input)

	d := script.Lines[0].(*ast.DialogueLine)

	var wait *ast.TimedWait
	for _, elem := range d.Content {
		if tw, ok := elem.(*ast.TimedWait); ok {
			wait = tw
			break
		}
	}

	if wait == nil {
		t.Fatal("expected to find TimedWait")
	}
	if wait.Skippable {
		t.Error("!w should not be skippable")
	}
	if wait.Duration != 500 {
		t.Errorf("duration: got %d, want 500", wait.Duration)
	}
}

func TestParse_NestedFormatTags(t *testing.T) {
	input := `d ` + "`{p:1:{c:FF0000:nested red}}`"
	script := Parse(input)

	d := script.Lines[0].(*ast.DialogueLine)
	text := getPlainText(d)

	if text != "nested red" {
		t.Errorf("plain text: got %q, want 'nested red'", text)
	}
}

func TestParse_RealExample(t *testing.T) {
	// A realistic example from Umineko
	input := `new_episode 1
preset_define 1,1,-1,#FF0000,0,0,0,0,0
*ep1_scene1
d2 [lv 0*"10"*"10100001"]` + "`{p:1:I'll make you understand!}`[@]"

	script := Parse(input)

	// Should have: episode marker, preset, label, dialogue
	if len(script.Lines) != 4 {
		t.Fatalf("expected 4 lines, got %d", len(script.Lines))
	}

	// Check episode marker
	em, ok := script.Lines[0].(*ast.EpisodeMarkerLine)
	if !ok {
		t.Fatalf("line 0: expected EpisodeMarkerLine, got %T", script.Lines[0])
	}
	if em.Episode != 1 {
		t.Errorf("episode: got %d, want 1", em.Episode)
	}

	// Check preset
	preset, ok := script.Lines[1].(*ast.PresetDefineLine)
	if !ok {
		t.Fatalf("line 1: expected PresetDefineLine, got %T", script.Lines[1])
	}
	if preset.ID != 1 || preset.Colour != "#FF0000" {
		t.Errorf("preset: got ID=%d, Colour=%s", preset.ID, preset.Colour)
	}

	// Check label
	label, ok := script.Lines[2].(*ast.LabelLine)
	if !ok {
		t.Fatalf("line 2: expected LabelLine, got %T", script.Lines[2])
	}
	if label.Name != "ep1_scene1" {
		t.Errorf("label: got %q, want 'ep1_scene1'", label.Name)
	}

	// Check dialogue
	d, ok := script.Lines[3].(*ast.DialogueLine)
	if !ok {
		t.Fatalf("line 3: expected DialogueLine, got %T", script.Lines[3])
	}

	voices := d.GetVoiceCommands()
	if len(voices) != 1 {
		t.Fatalf("expected 1 voice, got %d", len(voices))
	}
	if voices[0].CharacterID != "10" {
		t.Errorf("character ID: got %q, want '10'", voices[0].CharacterID)
	}

	text := getPlainText(d)
	if text != "I'll make you understand!" {
		t.Errorf("text: got %q", text)
	}
}

func TestParse_VoiceInsideFormatTag(t *testing.T) {
	input := `d2 ` + "`{a:c: `[lv 0*\"30\"*\"53200275\"][#][*]`\"{p:1:Ushiromiya Natsuhi is not the culprit}!\"} `[\\]"

	script := Parse(input)
	if len(script.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(script.Lines))
	}

	d, ok := script.Lines[0].(*ast.DialogueLine)
	if !ok {
		t.Fatalf("expected DialogueLine, got %T", script.Lines[0])
	}

	voices := d.GetVoiceCommands()
	if len(voices) != 1 {
		t.Fatalf("expected 1 voice command, got %d", len(voices))
	}

	if voices[0].CharacterID != "30" {
		t.Errorf("character ID: got %q, want '30'", voices[0].CharacterID)
	}
	if voices[0].AudioID != "53200275" {
		t.Errorf("audio ID: got %q, want '53200275'", voices[0].AudioID)
	}
}

func TestParse_NestedNobrWithCharSpacing(t *testing.T) {
	// {nobr:{m:-5:——}—} should produce "———" not "-5:——}—"
	input := `d ` + "`{nobr:{m:-5:——}—}`"
	script := Parse(input)

	d := script.Lines[0].(*ast.DialogueLine)
	text := getPlainText(d)

	if text != "———" {
		t.Errorf("plain text: got %q, want '———'", text)
	}
}

func TestParse_ConditionalNTag(t *testing.T) {
	// {n:0:human} is a conditional default tag, not a line break.
	// {y:0:Human} is a conditional non-default tag (stripped).
	input := `d ` + "`the {y:0:Human}{n:0:human} side`"
	script := Parse(input)

	d := script.Lines[0].(*ast.DialogueLine)
	text := getPlainText(d)

	if text != "the human side" {
		t.Errorf("plain text: got %q, want 'the human side'", text)
	}
}

func TestParse_LeadingSpaceAfterQuote(t *testing.T) {
	// '" text' — the quote mark should be trimmed and so should the space after it.
	input := `d ` + "`\" Everything I speak in red is the truth!\"`"
	script := Parse(input)

	d := script.Lines[0].(*ast.DialogueLine)
	text := getPlainText(d)

	if text != "Everything I speak in red is the truth!" {
		t.Errorf("plain text: got %q, want 'Everything I speak in red is the truth!'", text)
	}
}
