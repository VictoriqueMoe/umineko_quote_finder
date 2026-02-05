package lexar

import (
	"strings"
	"testing"

	"umineko_quote/internal/lexar/transformer"
)

func TestExtractQuotes_PresetColours(t *testing.T) {
	input := `preset_define 0,6,36,#FFFFFF,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1
preset_define 1,1,-1,#FF0000,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1
preset_define 2,1,-1,#39C6FF,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1
preset_define 41,1,-1,#FFAA00,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1
preset_define 42,1,-1,#AA71FF,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1
new_episode 1
d [lv 0*"27"*"10100001"]` + "`\"{p:1:Red truth test line here}.\"`" + `[\]
d [lv 0*"10"*"10100002"]` + "`\"{p:2:Blue truth test line here}.\"`" + `[\]
d [lv 0*"10"*"10100003"]` + "`\"{p:41:Gold truth test line here}.\"`" + `[\]
d [lv 0*"17"*"10100004"]` + "`\"{p:42:Purple truth test line}.\"`" + `[\]
d [lv 0*"10"*"10100005"]` + "`\"Text with {p:0:Japanese font preset} here.\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 5 {
		t.Fatalf("expected 5 quotes, got %d", len(quotes))
	}

	// Create transformer registry for HTML transformation
	registry := transformer.NewFactory(extractor.Presets())
	htmlTransformer := registry.MustGet(transformer.FormatHTML)
	plainTransformer := registry.MustGet(transformer.FormatPlainText)

	html0 := htmlTransformer.Transform(quotes[0].Content)
	html1 := htmlTransformer.Transform(quotes[1].Content)
	html2 := htmlTransformer.Transform(quotes[2].Content)
	html3 := htmlTransformer.Transform(quotes[3].Content)
	html4 := htmlTransformer.Transform(quotes[4].Content)
	text4 := plainTransformer.Transform(quotes[4].Content)

	// Red truth uses semantic class
	if !strings.Contains(html0, `class="red-truth"`) {
		t.Errorf("red truth should use class: %q", html0)
	}
	if strings.Contains(html0, "color:#FF0000") {
		t.Errorf("red truth should NOT use inline colour: %q", html0)
	}

	// Blue truth uses semantic class
	if !strings.Contains(html1, `class="blue-truth"`) {
		t.Errorf("blue truth should use class: %q", html1)
	}

	// Gold uses inline colour from parsed preset
	if !strings.Contains(html2, `color:#FFAA00`) {
		t.Errorf("gold should use parsed colour: %q", html2)
	}

	// Purple uses inline colour from parsed preset
	if !strings.Contains(html3, `color:#AA71FF`) {
		t.Errorf("purple should use parsed colour: %q", html3)
	}

	// White preset (0) should NOT add any colour span
	if strings.Contains(html4, "color:") {
		t.Errorf("white preset should not add colour: %q", html4)
	}
	if !strings.Contains(text4, "Japanese font preset") {
		t.Errorf("content should be preserved: %q", text4)
	}
}

func TestExtractQuotes_EpisodeMarkers(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		wantEpisode     int
		wantContentType string
	}{
		{
			name: "new_episode",
			input: `new_episode 3
d [lv 0*"10"*"30100001"]` + "`\"This is a test line for episode 3.\"`" + `[\]`,
			wantEpisode:     3,
			wantContentType: "",
		},
		{
			name: "new_tea",
			input: `new_tea 5
d [lv 0*"10"*"90100001"]` + "`\"Welcome to the tea party!\"`" + `[\]`,
			wantEpisode:     5,
			wantContentType: "tea",
		},
		{
			name: "new_ura",
			input: `new_ura 7
d [lv 0*"10"*"91000001"]` + "`\"This is the ura content line.\"`" + `[\]`,
			wantEpisode:     7,
			wantContentType: "ura",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := NewQuoteExtractor()
			quotes := extractor.ExtractQuotes(tt.input)

			if len(quotes) != 1 {
				t.Fatalf("expected 1 quote, got %d", len(quotes))
			}

			if quotes[0].Episode != tt.wantEpisode {
				t.Errorf("episode: got %d, want %d", quotes[0].Episode, tt.wantEpisode)
			}
			if quotes[0].ContentType != tt.wantContentType {
				t.Errorf("content type: got %q, want %q", quotes[0].ContentType, tt.wantContentType)
			}
		})
	}
}

func TestExtractQuotes_VoiceMetadata(t *testing.T) {
	input := `new_episode 1
d [lv 0*"19"*"11900001"]` + "`\"First part. `[@][lv 0*\"19\"*\"11900002\"]`Second part.\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	registry := transformer.NewFactory(extractor.Presets())
	plainTransformer := registry.MustGet(transformer.FormatPlainText)

	q := quotes[0]
	text := plainTransformer.Transform(q.Content)

	if q.CharacterID != "19" {
		t.Errorf("character ID: got %q, want '19'", q.CharacterID)
	}
	if q.AudioID != "11900001, 11900002" {
		t.Errorf("audio ID: got %q, want '11900001, 11900002'", q.AudioID)
	}
	if !strings.Contains(text, "First part") || !strings.Contains(text, "Second part") {
		t.Errorf("text should contain both parts: %q", text)
	}
}

func TestExtractQuotes_RedTruth(t *testing.T) {
	input := `preset_define 1,1,-1,#FF0000,0,0,0,0,0
new_episode 4
d [lv 0*"27"*"40700001"]` + "`\"{p:1:I speak the red truth!}\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	registry := transformer.NewFactory(extractor.Presets())
	htmlTransformer := registry.MustGet(transformer.FormatHTML)
	plainTransformer := registry.MustGet(transformer.FormatPlainText)

	q := quotes[0]
	html := htmlTransformer.Transform(q.Content)
	text := plainTransformer.Transform(q.Content)

	if !strings.Contains(html, `class="red-truth"`) {
		t.Errorf("expected red-truth class: %q", html)
	}
	if !strings.Contains(text, "I speak the red truth!") {
		t.Errorf("text content: %q", text)
	}
	if !q.Truth.HasRed {
		t.Errorf("Truth.HasRed: got false, want true")
	}
}

func TestExtractQuotes_BlueTruth(t *testing.T) {
	input := `preset_define 2,1,-1,#39C6FF,0,0,0,0,0
new_episode 5
d [lv 0*"10"*"50100001"]` + "`\"{p:2:Counter with blue truth!}\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	registry := transformer.NewFactory(extractor.Presets())
	htmlTransformer := registry.MustGet(transformer.FormatHTML)

	q := quotes[0]
	html := htmlTransformer.Transform(q.Content)

	if !strings.Contains(html, `class="blue-truth"`) {
		t.Errorf("expected blue-truth class: %q", html)
	}
	if !q.Truth.HasBlue {
		t.Errorf("Truth.HasBlue: got false, want true")
	}
}

func TestExtractQuotes_ColourFormatting(t *testing.T) {
	input := `new_episode 1
d [lv 0*"10"*"10100001"]` + "`\"This is {c:FF0000:red text} here.\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	registry := transformer.NewFactory(extractor.Presets())
	htmlTransformer := registry.MustGet(transformer.FormatHTML)
	plainTransformer := registry.MustGet(transformer.FormatPlainText)

	q := quotes[0]
	html := htmlTransformer.Transform(q.Content)
	text := plainTransformer.Transform(q.Content)

	if !strings.Contains(html, `color:#FF0000`) {
		t.Errorf("expected colour style: %q", html)
	}
	if text != "This is red text here." {
		t.Errorf("plain text: got %q, want 'This is red text here.'", text)
	}
	if q.Truth.HasRed || q.Truth.HasBlue {
		t.Errorf("Truth: got HasRed=%v HasBlue=%v, want both false (colour tag is not truth)", q.Truth.HasRed, q.Truth.HasBlue)
	}
}

func TestExtractQuotes_ItalicFormatting(t *testing.T) {
	input := `new_episode 1
d [lv 0*"10"*"10100001"]` + "`\"This is {i:italic text} here.\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	registry := transformer.NewFactory(extractor.Presets())
	htmlTransformer := registry.MustGet(transformer.FormatHTML)
	plainTransformer := registry.MustGet(transformer.FormatPlainText)

	q := quotes[0]
	html := htmlTransformer.Transform(q.Content)
	text := plainTransformer.Transform(q.Content)

	if !strings.Contains(html, "<em>italic text</em>") {
		t.Errorf("expected italic tags: %q", html)
	}
	if text != "This is italic text here." {
		t.Errorf("plain text: got %q", text)
	}
}

func TestExtractQuotes_RubyAnnotations(t *testing.T) {
	input := `new_episode 1
d [lv 0*"10"*"10100001"]` + "`\"The {ruby:Beatrice:Golden Witch} appears.\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	registry := transformer.NewFactory(extractor.Presets())
	htmlTransformer := registry.MustGet(transformer.FormatHTML)

	q := quotes[0]
	html := htmlTransformer.Transform(q.Content)

	if !strings.Contains(html, "<ruby>") {
		t.Errorf("expected ruby tags: %q", html)
	}
	if !strings.Contains(html, "Golden Witch") {
		t.Errorf("expected main text: %q", html)
	}
	if !strings.Contains(html, "Beatrice") {
		t.Errorf("expected annotation: %q", html)
	}
}

func TestExtractQuotes_LineBreaks(t *testing.T) {
	input := `new_episode 1
d [lv 0*"10"*"10100001"]` + "`\"Line one{n}Line two\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	registry := transformer.NewFactory(extractor.Presets())
	htmlTransformer := registry.MustGet(transformer.FormatHTML)
	plainTransformer := registry.MustGet(transformer.FormatPlainText)

	q := quotes[0]
	html := htmlTransformer.Transform(q.Content)
	text := plainTransformer.Transform(q.Content)

	if !strings.Contains(html, "<br>") {
		t.Errorf("expected <br> in HTML: %q", html)
	}
	if text != "Line one Line two" {
		t.Errorf("plain text: got %q, want 'Line one Line two'", text)
	}
}

func TestExtractQuotes_NestedTags(t *testing.T) {
	input := `preset_define 1,1,-1,#FF0000,0,0,0,0,0
new_episode 4
d [lv 0*"27"*"40700001"]` + "`\"{p:1:{c:FFFFFF:Nested colour in red truth}}\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	registry := transformer.NewFactory(extractor.Presets())
	htmlTransformer := registry.MustGet(transformer.FormatHTML)
	plainTransformer := registry.MustGet(transformer.FormatPlainText)

	q := quotes[0]
	html := htmlTransformer.Transform(q.Content)
	text := plainTransformer.Transform(q.Content)

	if !strings.Contains(html, `class="red-truth"`) {
		t.Errorf("expected red-truth class: %q", html)
	}
	// The nested colour should also be present
	if !strings.Contains(html, `color:#FFFFFF`) {
		t.Errorf("expected nested colour: %q", html)
	}
	if !strings.Contains(text, "Nested colour in red truth") {
		t.Errorf("plain text: %q", text)
	}
}

func TestExtractQuotes_SpecialCharacters(t *testing.T) {
	input := `new_episode 1
d [lv 0*"10"*"10100001"]` + "`\"She said {qt}hello{qt} to me.\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	registry := transformer.NewFactory(extractor.Presets())
	plainTransformer := registry.MustGet(transformer.FormatPlainText)

	q := quotes[0]
	text := plainTransformer.Transform(q.Content)

	if !strings.Contains(text, `"hello"`) {
		t.Errorf("expected quote characters: %q", text)
	}
}

func TestExtractQuotes_YTagStripped(t *testing.T) {
	input := `new_episode 1
d [lv 0*"10"*"10100001"]` + "`\"Visible text{y:1:hidden Japanese} continues.\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	registry := transformer.NewFactory(extractor.Presets())
	plainTransformer := registry.MustGet(transformer.FormatPlainText)

	q := quotes[0]
	text := plainTransformer.Transform(q.Content)

	if strings.Contains(text, "hidden") {
		t.Errorf("y: content should be stripped: %q", text)
	}
	if text != "Visible text continues." {
		t.Errorf("plain text: got %q, want 'Visible text continues.'", text)
	}
}

func TestExtractQuotes_MultipleVoiceChannels(t *testing.T) {
	// Test with channel 1 as well as channel 0
	input := `new_episode 1
d [lv 1*"10"*"10100001"]` + "`\"Channel 1 voice.\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.CharacterID != "10" {
		t.Errorf("character ID: got %q, want '10'", q.CharacterID)
	}
	if q.AudioID != "10100001" {
		t.Errorf("audio ID: got %q, want '10100001'", q.AudioID)
	}
}

func TestExtractQuotes_D2Command(t *testing.T) {
	input := `new_episode 5
d2 [lv 0*"27"*"50700001"]` + "`\"You idiot!! `[@][lv 0*\"27\"*\"50700002\"]`Isn't it obvious?!\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	registry := transformer.NewFactory(extractor.Presets())
	plainTransformer := registry.MustGet(transformer.FormatPlainText)

	q := quotes[0]
	text := plainTransformer.Transform(q.Content)

	if q.CharacterID != "27" {
		t.Errorf("character ID: got %q, want '27'", q.CharacterID)
	}
	if q.AudioID != "50700001, 50700002" {
		t.Errorf("audio ID: got %q, want '50700001, 50700002'", q.AudioID)
	}
	if !strings.Contains(text, "You idiot") {
		t.Errorf("text should contain 'You idiot': %q", text)
	}
}

func TestExtractQuotes_EpisodeFromAudioID(t *testing.T) {
	input := `d [lv 0*"10"*"50100001"]` + "`\"Episode 5 inferred from audio ID.\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	// Episode should be inferred from audio ID starting with 5
	if quotes[0].Episode != 5 {
		t.Errorf("episode: got %d, want 5 (from audio ID)", quotes[0].Episode)
	}
}

func TestExtractQuotes_TruthDetection(t *testing.T) {
	presets := `preset_define 1,1,-1,#FF0000,0,0,0,0,0
preset_define 2,1,-1,#39C6FF,0,0,0,0,0
`

	tests := []struct {
		name     string
		dialogue string
		wantRed  bool
		wantBlue bool
	}{
		{
			name:     "red only",
			dialogue: `d [lv 0*"27"*"10100001"]` + "`\"{p:1:This is red truth}.\"`" + `[\]`,
			wantRed:  true,
			wantBlue: false,
		},
		{
			name:     "blue only",
			dialogue: `d [lv 0*"10"*"50100001"]` + "`\"{p:2:This is blue truth}.\"`" + `[\]`,
			wantRed:  false,
			wantBlue: true,
		},
		{
			name:     "red then blue",
			dialogue: `d [lv 0*"27"*"20700001"]` + "`\"{p:1:Red first}. `[@]`{p:2:Then blue}.\"`" + `[\]`,
			wantRed:  true,
			wantBlue: true,
		},
		{
			name:     "blue then red",
			dialogue: `d [lv 0*"10"*"50100001"]` + "`\"{p:2:Blue first}. `[@]`{p:1:Then red}.\"`" + `[\]`,
			wantRed:  true,
			wantBlue: true,
		},
		{
			name:     "nested blue inside red",
			dialogue: `d [lv 0*"27"*"10100001"]` + "`\"{p:1:Red with {p:2:blue nested} inside}.\"`" + `[\]`,
			wantRed:  true,
			wantBlue: true,
		},
		{
			name:     "nested red inside blue",
			dialogue: `d [lv 0*"10"*"50100001"]` + "`\"{p:2:Blue with {p:1:red nested} inside}.\"`" + `[\]`,
			wantRed:  true,
			wantBlue: true,
		},
		{
			name:     "neither truth",
			dialogue: `d [lv 0*"10"*"10100001"]` + "`\"Just normal dialogue.\"`" + `[\]`,
			wantRed:  false,
			wantBlue: false,
		},
		{
			name:     "colour tag is not truth",
			dialogue: `d [lv 0*"10"*"10100001"]` + "`\"This is {c:FF0000:red coloured} but not truth.\"`" + `[\]`,
			wantRed:  false,
			wantBlue: false,
		},
		{
			name:     "real mixed quote from game (blue first)",
			dialogue: `d2 [lv 0*"47"*"54600077"]` + "`\"{p:2:Was the existence of Natsuhi's blind spot depicted?}. `[@]`{p:1:Knox's 8th}!\"`" + `[\]`,
			wantRed:  true,
			wantBlue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := presets + tt.dialogue

			extractor := NewQuoteExtractor()
			quotes := extractor.ExtractQuotes(input)

			if len(quotes) != 1 {
				t.Fatalf("expected 1 quote, got %d", len(quotes))
			}

			q := quotes[0]
			if q.Truth.HasRed != tt.wantRed {
				t.Errorf("HasRed: got %v, want %v", q.Truth.HasRed, tt.wantRed)
			}
			if q.Truth.HasBlue != tt.wantBlue {
				t.Errorf("HasBlue: got %v, want %v", q.Truth.HasBlue, tt.wantBlue)
			}
		})
	}
}

func TestExtractQuotes_NarrationWithEmbeddedVoice(t *testing.T) {
	// Narration text that has voiced audio clips at the end should be
	// attributed to narrator, not the voiced character.
	input := "preset_define 0,6,36,#FFFFFF,0,0,0,1,-1,#000000,0,-1,-1,#000000,1,-1\n" +
		"*o4_16\n" +
		"d `Furthermore, there are unused voice files left on the disc.`[@]` To our relief, they ultimately found this unnecessary.`[@][lv 0*\"28\"*\"92100173\"]` Mii,`[|][lv 0*\"28\"*\"92100174\"]` nipah~{p:0:" + "\u2606" + "}`[\\]"

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.CharacterID != "narrator" {
		t.Errorf("characterID: got %q, want \"narrator\"", q.CharacterID)
	}
	if q.AudioID != "" {
		t.Errorf("audioID: got %q, want empty", q.AudioID)
	}
}

func TestExtractQuotes_DotsBeforeVoiceIsCharacter(t *testing.T) {
	// Dots/ellipsis before a voice command represent character pauses,
	// not narration. The character should still be detected.
	input := `new_episode 3
d ` + "`\"............ `[@][#][*][lv 0*\"01\"*\"31500076\"]`...I will bring some tea now.`[\\]"

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.CharacterID != "01" {
		t.Errorf("characterID: got %q, want \"01\"", q.CharacterID)
	}
	if q.AudioID != "31500076" {
		t.Errorf("audioID: got %q, want \"31500076\"", q.AudioID)
	}
}

func TestExtractQuotes_MultiCharacterAudioCharMap(t *testing.T) {
	// Multiple voice commands from different characters should populate AudioCharMap.
	input := `new_episode 3
d2 [lv 6*"38"*"32300003"][lv 5*"39"*"32400003"][lv 4*"40"*"32500003"]` + "`\"Eeep!\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.CharacterID != "38" {
		t.Errorf("primary character ID: got %q, want \"38\"", q.CharacterID)
	}
	if q.AudioID != "32300003, 32400003, 32500003" {
		t.Errorf("audio ID: got %q, want \"32300003, 32400003, 32500003\"", q.AudioID)
	}
	if q.AudioCharMap == nil {
		t.Fatal("AudioCharMap should not be nil for multi-character quote")
	}
	if q.AudioCharMap["32300003"] != "38" {
		t.Errorf("AudioCharMap[32300003]: got %q, want \"38\"", q.AudioCharMap["32300003"])
	}
	if q.AudioCharMap["32400003"] != "39" {
		t.Errorf("AudioCharMap[32400003]: got %q, want \"39\"", q.AudioCharMap["32400003"])
	}
	if q.AudioCharMap["32500003"] != "40" {
		t.Errorf("AudioCharMap[32500003]: got %q, want \"40\"", q.AudioCharMap["32500003"])
	}
}

func TestExtractQuotes_SingleCharacterNoAudioCharMap(t *testing.T) {
	// Same character across multiple voice commands should NOT populate AudioCharMap.
	input := `new_episode 1
d [lv 0*"19"*"11900001"]` + "`\"First. `[@][lv 0*\"19\"*\"11900002\"]`Second.\"`" + `[\]`

	extractor := NewQuoteExtractor()
	quotes := extractor.ExtractQuotes(input)

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	if quotes[0].AudioCharMap != nil {
		t.Errorf("AudioCharMap should be nil for single-character quote, got %v", quotes[0].AudioCharMap)
	}
}
