package transformer

import (
	"strings"
	"testing"

	"umineko_quote/internal/lexar/ast"
)

func TestDefaultSemanticPresets(t *testing.T) {
	presets := DefaultSemanticPresets()

	if presets["1"] != "red-truth" {
		t.Errorf("preset 1 = %q, want %q", presets["1"], "red-truth")
	}
	if presets["2"] != "blue-truth" {
		t.Errorf("preset 2 = %q, want %q", presets["2"], "blue-truth")
	}
	if len(presets) != 2 {
		t.Errorf("expected 2 default presets, got %d", len(presets))
	}
}

func TestNewPresetContext(t *testing.T) {
	ctx := NewPresetContext()

	if ctx.SemanticPresets == nil {
		t.Fatal("SemanticPresets is nil")
	}
	if ctx.DynamicColours == nil {
		t.Fatal("DynamicColours is nil")
	}
	if ctx.GetSemanticClass("1") != "red-truth" {
		t.Error("missing default red-truth preset")
	}
	if ctx.GetSemanticClass("2") != "blue-truth" {
		t.Error("missing default blue-truth preset")
	}
}

func TestGetSemanticClass(t *testing.T) {
	ctx := NewPresetContext()

	if got := ctx.GetSemanticClass("1"); got != "red-truth" {
		t.Errorf("GetSemanticClass(1) = %q, want %q", got, "red-truth")
	}
	if got := ctx.GetSemanticClass("99"); got != "" {
		t.Errorf("GetSemanticClass(99) = %q, want empty", got)
	}
}

func TestGetDynamicColour(t *testing.T) {
	ctx := NewPresetContext()
	ctx.DynamicColours["42"] = "#AA71FF"

	if got := ctx.GetDynamicColour("42"); got != "#AA71FF" {
		t.Errorf("GetDynamicColour(42) = %q, want %q", got, "#AA71FF")
	}
	if got := ctx.GetDynamicColour("99"); got != "" {
		t.Errorf("GetDynamicColour(99) = %q, want empty", got)
	}
}

func TestCollectFromScript(t *testing.T) {
	ctx := NewPresetContext()
	script := &ast.Script{
		Lines: []ast.Line{
			&ast.PresetDefineLine{ID: 41, Colour: "#FFAA00"},
			&ast.PresetDefineLine{ID: 42, Colour: "#aa71ff"},
		},
	}

	ctx.CollectFromScript(script)

	if got := ctx.GetDynamicColour("41"); got != "#FFAA00" {
		t.Errorf("preset 41 colour = %q, want %q", got, "#FFAA00")
	}
	if got := ctx.GetDynamicColour("42"); got != "#AA71FF" {
		t.Errorf("preset 42 colour = %q, want %q (should be uppercased)", got, "#AA71FF")
	}
}

func TestCollectFromScript_SkipsSemanticPresets(t *testing.T) {
	ctx := NewPresetContext()
	script := &ast.Script{
		Lines: []ast.Line{
			&ast.PresetDefineLine{ID: 1, Colour: "#FF0000"},
			&ast.PresetDefineLine{ID: 2, Colour: "#39C6FF"},
		},
	}

	ctx.CollectFromScript(script)

	if got := ctx.GetDynamicColour("1"); got != "" {
		t.Errorf("semantic preset 1 should not be in DynamicColours, got %q", got)
	}
	if got := ctx.GetDynamicColour("2"); got != "" {
		t.Errorf("semantic preset 2 should not be in DynamicColours, got %q", got)
	}
}

func TestCollectFromScript_SkipsWhiteAndEmpty(t *testing.T) {
	ctx := NewPresetContext()
	script := &ast.Script{
		Lines: []ast.Line{
			&ast.PresetDefineLine{ID: 0, Colour: "#FFFFFF"},
			&ast.PresetDefineLine{ID: 3, Colour: "#ffffff"},
			&ast.PresetDefineLine{ID: 5, Colour: ""},
		},
	}

	ctx.CollectFromScript(script)

	if len(ctx.DynamicColours) != 0 {
		t.Errorf("expected 0 dynamic colours, got %d: %v", len(ctx.DynamicColours), ctx.DynamicColours)
	}
}

func TestCollectFromScript_ResetsOnEachCall(t *testing.T) {
	ctx := NewPresetContext()

	script1 := &ast.Script{
		Lines: []ast.Line{
			&ast.PresetDefineLine{ID: 41, Colour: "#FFAA00"},
		},
	}
	ctx.CollectFromScript(script1)

	if ctx.GetDynamicColour("41") == "" {
		t.Fatal("expected preset 41 after first collection")
	}

	script2 := &ast.Script{Lines: []ast.Line{}}
	ctx.CollectFromScript(script2)

	if got := ctx.GetDynamicColour("41"); got != "" {
		t.Errorf("expected preset 41 to be cleared after second collection, got %q", got)
	}
}

func TestCollectFromScript_IgnoresNonPresetLines(t *testing.T) {
	ctx := NewPresetContext()
	script := &ast.Script{
		Lines: []ast.Line{
			&ast.CommentLine{Text: "this is a comment"},
			&ast.PresetDefineLine{ID: 42, Colour: "#AA71FF"},
			&ast.EpisodeMarkerLine{Episode: 1},
			&ast.DialogueLine{Command: "d"},
		},
	}

	ctx.CollectFromScript(script)

	if got := ctx.GetDynamicColour("42"); got != "#AA71FF" {
		t.Errorf("preset 42 = %q, want %q", got, "#AA71FF")
	}
	if len(ctx.DynamicColours) != 1 {
		t.Errorf("expected 1 dynamic colour, got %d", len(ctx.DynamicColours))
	}
}

// HtmlTransformer tests

func newTestHtmlTransformer() *HtmlTransformer {
	ctx := NewPresetContext()
	ctx.DynamicColours["41"] = "#FFAA00"
	ctx.DynamicColours["42"] = "#AA71FF"
	return NewHtmlTransformer(ctx)
}

func TestHtml_PlainText(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: "Hello, world."},
	}

	got := tr.Transform(elements)
	if got != "Hello, world." {
		t.Errorf("got %q, want %q", got, "Hello, world.")
	}
}

func TestHtml_PlainText_HtmlEscaping(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: `She said <script>alert("xss")</script>.`},
	}

	got := tr.Transform(elements)
	if strings.Contains(got, "<script>") {
		t.Errorf("HTML should be escaped: %q", got)
	}
	if !strings.Contains(got, "&lt;script&gt;") {
		t.Errorf("expected escaped tags: %q", got)
	}
}

func TestHtml_SpecialChars(t *testing.T) {
	tests := []struct {
		name string
		char string
		want string
	}{
		{"newline", "n", "<br>"},
		{"quote", "qt", `"`},
		{"open bracket", "os", "["},
		{"close bracket", "es", "]"},
	}

	tr := newTestHtmlTransformer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := []ast.DialogueElement{
				&ast.PlainText{Text: "A"},
				&ast.SpecialChar{Name: tt.char},
				&ast.PlainText{Text: "B"},
			}
			got := tr.Transform(elements)
			if !strings.Contains(got, tt.want) {
				t.Errorf("got %q, expected to contain %q", got, tt.want)
			}
		})
	}
}

func TestHtml_SpecialChar_Ignored(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: "A"},
		&ast.SpecialChar{Name: "0"},
		&ast.PlainText{Text: "B"},
	}

	got := tr.Transform(elements)
	if got != "AB" {
		t.Errorf("got %q, want %q", got, "AB")
	}
}

func TestHtml_Italic(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:    "i",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "emphasis"}},
		},
	}

	got := tr.Transform(elements)
	if got != "<em>emphasis</em>" {
		t.Errorf("got %q, want %q", got, "<em>emphasis</em>")
	}
}

func TestHtml_ItalicLongForm(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:    "italic",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "text"}},
		},
	}

	got := tr.Transform(elements)
	if got != "<em>text</em>" {
		t.Errorf("got %q, want %q", got, "<em>text</em>")
	}
}

func TestHtml_Colour(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:    "c",
			Param:   "FF0000",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "red"}},
		},
	}

	got := tr.Transform(elements)
	if got != `<span style="color:#FF0000">red</span>` {
		t.Errorf("got %q", got)
	}
}

func TestHtml_ColourLongForms(t *testing.T) {
	for _, name := range []string{"color", "colour"} {
		t.Run(name, func(t *testing.T) {
			tr := newTestHtmlTransformer()
			elements := []ast.DialogueElement{
				&ast.FormatTag{
					Name:    name,
					Param:   "00FF00",
					Content: []ast.DialogueElement{&ast.PlainText{Text: "green"}},
				},
			}

			got := tr.Transform(elements)
			if !strings.Contains(got, `color:#00FF00`) {
				t.Errorf("got %q, expected colour style", got)
			}
		})
	}
}

func TestHtml_Ruby(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:    "ruby",
			Param:   "Beatrice",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "Golden Witch"}},
		},
	}

	got := tr.Transform(elements)
	if !strings.Contains(got, "<ruby>") {
		t.Errorf("expected <ruby> tag: %q", got)
	}
	if !strings.Contains(got, "<rt>Beatrice</rt>") {
		t.Errorf("expected annotation: %q", got)
	}
	if !strings.Contains(got, "Golden Witch") {
		t.Errorf("expected main text: %q", got)
	}
}

func TestHtml_RubyShortForm(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:    "h",
			Param:   "reading",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "kanji"}},
		},
	}

	got := tr.Transform(elements)
	if !strings.Contains(got, "<ruby>") || !strings.Contains(got, "<rt>reading</rt>") {
		t.Errorf("h tag should produce ruby: %q", got)
	}
}

func TestHtml_Ruby_EscapesAnnotation(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:    "ruby",
			Param:   `<b>bold</b>`,
			Content: []ast.DialogueElement{&ast.PlainText{Text: "text"}},
		},
	}

	got := tr.Transform(elements)
	if strings.Contains(got, "<b>bold</b>") {
		t.Errorf("ruby annotation should be escaped: %q", got)
	}
}

func TestHtml_PresetRedTruth(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:    "p",
			Param:   "1",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "I speak the red truth!"}},
		},
	}

	got := tr.Transform(elements)
	if !strings.Contains(got, `class="red-truth"`) {
		t.Errorf("expected red-truth class: %q", got)
	}
	if strings.Contains(got, "color:") {
		t.Errorf("semantic preset should not use inline colour: %q", got)
	}
}

func TestHtml_PresetBlueTruth(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:    "preset",
			Param:   "2",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "Counter!"}},
		},
	}

	got := tr.Transform(elements)
	if !strings.Contains(got, `class="blue-truth"`) {
		t.Errorf("expected blue-truth class: %q", got)
	}
}

func TestHtml_PresetDynamicColour(t *testing.T) {
	tests := []struct {
		name   string
		preset string
		colour string
	}{
		{"gold", "41", "#FFAA00"},
		{"purple", "42", "#AA71FF"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := newTestHtmlTransformer()
			elements := []ast.DialogueElement{
				&ast.FormatTag{
					Name:    "p",
					Param:   tt.preset,
					Content: []ast.DialogueElement{&ast.PlainText{Text: "styled text"}},
				},
			}

			got := tr.Transform(elements)
			if !strings.Contains(got, "color:"+tt.colour) {
				t.Errorf("expected colour %s: %q", tt.colour, got)
			}
		})
	}
}

func TestHtml_PresetUnknown_FallsThrough(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:    "p",
			Param:   "99",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "unstyled"}},
		},
	}

	got := tr.Transform(elements)
	if got != "unstyled" {
		t.Errorf("unknown preset should render content only: got %q", got)
	}
	if strings.Contains(got, "span") {
		t.Errorf("unknown preset should not produce span: %q", got)
	}
}

func TestHtml_YTag_Stripped(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: "Visible"},
		&ast.FormatTag{
			Name:    "y",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "hidden"}},
		},
		&ast.PlainText{Text: " continues."},
	}

	got := tr.Transform(elements)
	if strings.Contains(got, "hidden") {
		t.Errorf("y tag content should be stripped: %q", got)
	}
	if got != "Visible continues." {
		t.Errorf("got %q, want %q", got, "Visible continues.")
	}
}

func TestHtml_NTag_WithContent(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:    "n",
			Param:   "1",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "conditional text"}},
		},
	}

	got := tr.Transform(elements)
	if got != "conditional text" {
		t.Errorf("got %q, want %q", got, "conditional text")
	}
}

func TestHtml_UnknownTag_PassesContent(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:    "nobr",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "no break"}},
		},
	}

	got := tr.Transform(elements)
	if got != "no break" {
		t.Errorf("got %q, want %q", got, "no break")
	}
}

func TestHtml_NestedTags(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:  "p",
			Param: "1",
			Content: []ast.DialogueElement{
				&ast.PlainText{Text: "Red with "},
				&ast.FormatTag{
					Name:    "i",
					Content: []ast.DialogueElement{&ast.PlainText{Text: "italic"}},
				},
			},
		},
	}

	got := tr.Transform(elements)
	if !strings.Contains(got, `class="red-truth"`) {
		t.Errorf("expected red-truth class: %q", got)
	}
	if !strings.Contains(got, "<em>italic</em>") {
		t.Errorf("expected nested italic: %q", got)
	}
}

func TestHtml_TrimQuotes(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: `"Hello world."`},
	}

	got := tr.Transform(elements)
	if got != "Hello world." {
		t.Errorf("got %q, want %q", got, "Hello world.")
	}
}

func TestHtml_StrayBracesRemoved(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: "Text with { stray } braces."},
	}

	got := tr.Transform(elements)
	if strings.Contains(got, "{") || strings.Contains(got, "}") {
		t.Errorf("stray braces should be removed: %q", got)
	}
}

func TestHtml_EmptyElements(t *testing.T) {
	tr := newTestHtmlTransformer()
	got := tr.Transform(nil)
	if got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestHtml_ClickWaitIgnored(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: "Part one. "},
		&ast.ClickWait{Type: "@"},
		&ast.PlainText{Text: "Part two."},
	}

	got := tr.Transform(elements)
	if got != "Part one. Part two." {
		t.Errorf("got %q, want %q", got, "Part one. Part two.")
	}
}

func TestHtml_VoiceCommandIgnored(t *testing.T) {
	tr := newTestHtmlTransformer()
	elements := []ast.DialogueElement{
		&ast.VoiceCommand{CharacterID: "10", AudioID: "10100001"},
		&ast.PlainText{Text: "Dialogue."},
	}

	got := tr.Transform(elements)
	if got != "Dialogue." {
		t.Errorf("got %q, want %q", got, "Dialogue.")
	}
}

// PlainTextTransformer tests

func TestPlain_PlainText(t *testing.T) {
	tr := NewPlainTextTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: "Hello, world."},
	}

	got := tr.Transform(elements)
	if got != "Hello, world." {
		t.Errorf("got %q, want %q", got, "Hello, world.")
	}
}

func TestPlain_NoHtmlEscaping(t *testing.T) {
	tr := NewPlainTextTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: "A & B < C"},
	}

	got := tr.Transform(elements)
	if got != "A & B < C" {
		t.Errorf("plain text should not HTML-escape: got %q", got)
	}
}

func TestPlain_SpecialChars(t *testing.T) {
	tests := []struct {
		name string
		char string
		want string
	}{
		{"newline becomes space", "n", " "},
		{"quote", "qt", `"`},
		{"open bracket", "os", "["},
		{"close bracket", "es", "]"},
	}

	tr := NewPlainTextTransformer()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			elements := []ast.DialogueElement{
				&ast.PlainText{Text: "A"},
				&ast.SpecialChar{Name: tt.char},
				&ast.PlainText{Text: "B"},
			}
			got := tr.Transform(elements)
			want := "A" + tt.want + "B"
			if got != want {
				t.Errorf("got %q, want %q", got, want)
			}
		})
	}
}

func TestPlain_YTag_Stripped(t *testing.T) {
	tr := NewPlainTextTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: "Visible"},
		&ast.FormatTag{
			Name:    "y",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "hidden"}},
		},
		&ast.PlainText{Text: " continues."},
	}

	got := tr.Transform(elements)
	if strings.Contains(got, "hidden") {
		t.Errorf("y tag should be stripped: %q", got)
	}
	if got != "Visible continues." {
		t.Errorf("got %q, want %q", got, "Visible continues.")
	}
}

func TestPlain_Ruby(t *testing.T) {
	tr := NewPlainTextTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: "The "},
		&ast.FormatTag{
			Name:    "ruby",
			Param:   "Beatrice",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "Golden Witch"}},
		},
		&ast.PlainText{Text: " appears."},
	}

	got := tr.Transform(elements)
	if got != "The Golden Witch (Beatrice) appears." {
		t.Errorf("got %q, want %q", got, "The Golden Witch (Beatrice) appears.")
	}
}

func TestPlain_RubyShortForm(t *testing.T) {
	tr := NewPlainTextTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:    "h",
			Param:   "reading",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "text"}},
		},
	}

	got := tr.Transform(elements)
	if got != "text (reading)" {
		t.Errorf("got %q, want %q", got, "text (reading)")
	}
}

func TestPlain_RubyNoParam(t *testing.T) {
	tr := NewPlainTextTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name:    "ruby",
			Content: []ast.DialogueElement{&ast.PlainText{Text: "text"}},
		},
	}

	got := tr.Transform(elements)
	if got != "text" {
		t.Errorf("got %q, want %q (no annotation expected)", got, "text")
	}
}

func TestPlain_FormatTagsPassContent(t *testing.T) {
	tags := []string{"i", "italic", "c", "color", "p", "preset", "nobr", "bold", "f"}

	tr := NewPlainTextTransformer()
	for _, tag := range tags {
		t.Run(tag, func(t *testing.T) {
			elements := []ast.DialogueElement{
				&ast.FormatTag{
					Name:    tag,
					Content: []ast.DialogueElement{&ast.PlainText{Text: "content"}},
				},
			}
			got := tr.Transform(elements)
			if got != "content" {
				t.Errorf("tag %q: got %q, want %q", tag, got, "content")
			}
		})
	}
}

func TestPlain_NestedTags(t *testing.T) {
	tr := NewPlainTextTransformer()
	elements := []ast.DialogueElement{
		&ast.FormatTag{
			Name: "p",
			Content: []ast.DialogueElement{
				&ast.PlainText{Text: "Red with "},
				&ast.FormatTag{
					Name:    "i",
					Content: []ast.DialogueElement{&ast.PlainText{Text: "italic"}},
				},
			},
		},
	}

	got := tr.Transform(elements)
	if got != "Red with italic" {
		t.Errorf("got %q, want %q", got, "Red with italic")
	}
}

func TestPlain_TrimQuotes(t *testing.T) {
	tr := NewPlainTextTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: `"Hello world."`},
	}

	got := tr.Transform(elements)
	if got != "Hello world." {
		t.Errorf("got %q, want %q", got, "Hello world.")
	}
}

func TestPlain_StrayBracesRemoved(t *testing.T) {
	tr := NewPlainTextTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: "Text with { stray } braces."},
	}

	got := tr.Transform(elements)
	if strings.Contains(got, "{") || strings.Contains(got, "}") {
		t.Errorf("stray braces should be removed: %q", got)
	}
}

func TestPlain_EmptyElements(t *testing.T) {
	tr := NewPlainTextTransformer()
	got := tr.Transform(nil)
	if got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestPlain_ClickWaitIgnored(t *testing.T) {
	tr := NewPlainTextTransformer()
	elements := []ast.DialogueElement{
		&ast.PlainText{Text: "Part one. "},
		&ast.ClickWait{Type: "@"},
		&ast.PlainText{Text: "Part two."},
	}

	got := tr.Transform(elements)
	if got != "Part one. Part two." {
		t.Errorf("got %q, want %q", got, "Part one. Part two.")
	}
}

// Factory tests

func TestFactory_DefaultTransformers(t *testing.T) {
	f := NewFactory(NewPresetContext())

	plain, err := f.Get(FormatPlainText)
	if err != nil {
		t.Fatalf("Get(FormatPlainText) error: %v", err)
	}
	if plain == nil {
		t.Fatal("plain text transformer is nil")
	}

	html, err := f.Get(FormatHTML)
	if err != nil {
		t.Fatalf("Get(FormatHTML) error: %v", err)
	}
	if html == nil {
		t.Fatal("html transformer is nil")
	}
}

func TestFactory_Get_UnknownFormat(t *testing.T) {
	f := NewFactory(NewPresetContext())

	_, err := f.Get(Format(999))
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestFactory_MustGet_Panics(t *testing.T) {
	f := NewFactory(NewPresetContext())

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustGet should panic for unknown format")
		}
	}()

	f.MustGet(Format(999))
}

func TestFactory_MustGet_Success(t *testing.T) {
	f := NewFactory(NewPresetContext())

	tr := f.MustGet(FormatPlainText)
	if tr == nil {
		t.Error("MustGet returned nil")
	}
}

type mockTransformer struct{}

func (m *mockTransformer) Transform(elements []ast.DialogueElement) string {
	return "mock"
}

func TestFactory_Register(t *testing.T) {
	f := NewFactory(NewPresetContext())

	customFormat := Format(100)
	f.Register(customFormat, &mockTransformer{})

	tr, err := f.Get(customFormat)
	if err != nil {
		t.Fatalf("Get after Register: %v", err)
	}

	got := tr.Transform(nil)
	if got != "mock" {
		t.Errorf("custom transformer: got %q, want %q", got, "mock")
	}
}

func TestFactory_Register_OverridesExisting(t *testing.T) {
	f := NewFactory(NewPresetContext())

	f.Register(FormatPlainText, &mockTransformer{})

	tr := f.MustGet(FormatPlainText)
	got := tr.Transform(nil)
	if got != "mock" {
		t.Errorf("overridden transformer: got %q, want %q", got, "mock")
	}
}

func TestFactory_Presets(t *testing.T) {
	ctx := NewPresetContext()
	f := NewFactory(ctx)

	if f.Presets() != ctx {
		t.Error("Presets() should return the same PresetContext")
	}
}
