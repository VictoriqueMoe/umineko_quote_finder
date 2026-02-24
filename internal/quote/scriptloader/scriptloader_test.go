package scriptloader

import (
	"testing"
	"testing/fstest"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/lexar"
)

func testParseFunc(quotes []dto.ParsedQuote, refs []lexar.SubtitleRef) ParseFunc {
	return func(lines []string) ([]dto.ParsedQuote, []lexar.SubtitleRef) {
		return quotes, refs
	}
}

func buildTestFS(path string, plaintext []byte) fstest.MapFS {
	return fstest.MapFS{
		path: &fstest.MapFile{Data: encodeTestPayload(plaintext)},
	}
}

func TestLoad_ReturnsQuotes(t *testing.T) {
	expected := []dto.ParsedQuote{
		{Text: "Hello", CharacterID: "10", Episode: 1},
		{Text: "World", CharacterID: "27", Episode: 1},
	}
	fs := buildTestFS("data/test.file", []byte("line1\nline2"))
	loader := New(fs, testParseFunc(expected, nil))

	result := loader.Load("en", "data/test.file")

	if len(result) != 2 {
		t.Fatalf("expected 2 quotes, got %d", len(result))
	}
	if result[0].Text != "Hello" {
		t.Errorf("quote 0 text: got %q, want %q", result[0].Text, "Hello")
	}
	if result[1].Text != "World" {
		t.Errorf("quote 1 text: got %q, want %q", result[1].Text, "World")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	fs := fstest.MapFS{}
	loader := New(fs, testParseFunc(nil, nil))

	result := loader.Load("en", "data/missing.file")

	if result != nil {
		t.Errorf("expected nil for missing file, got %d quotes", len(result))
	}
}

func TestLoad_InvalidEncodedData(t *testing.T) {
	fs := fstest.MapFS{
		"data/bad.file": &fstest.MapFile{Data: []byte("not valid ONS2 data here!")},
	}
	loader := New(fs, testParseFunc(nil, nil))

	result := loader.Load("en", "data/bad.file")

	if result != nil {
		t.Errorf("expected nil for invalid data, got %d quotes", len(result))
	}
}

func TestLoad_PassesDecodedLinesToParser(t *testing.T) {
	plaintext := []byte("first line\nsecond line\nthird line")
	fs := buildTestFS("data/test.file", plaintext)

	var capturedLines []string
	parse := func(lines []string) ([]dto.ParsedQuote, []lexar.SubtitleRef) {
		capturedLines = lines
		return nil, nil
	}
	loader := New(fs, parse)

	loader.Load("en", "data/test.file")

	if len(capturedLines) != 3 {
		t.Fatalf("expected 3 lines passed to parser, got %d", len(capturedLines))
	}
	if capturedLines[0] != "first line" {
		t.Errorf("line 0: got %q, want %q", capturedLines[0], "first line")
	}
	if capturedLines[2] != "third line" {
		t.Errorf("line 2: got %q, want %q", capturedLines[2], "third line")
	}
}

func TestLoad_ResolvesSubtitleRefs(t *testing.T) {
	assContent := "[Script Info]\nTitle: Test\n\n[V4+ Styles]\nFormat: Name\nStyle: Default\n\n[Events]\nFormat: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\nDialogue: 0,0:00:00.00,0:00:05.00,Default,,0,0,0,,Welcome back.\nDialogue: 0,0:00:05.00,0:00:10.00,Default,,0,0,0,,Goodbye.\n"

	fs := fstest.MapFS{
		"data/test.file":      &fstest.MapFile{Data: encodeTestPayload([]byte("test"))},
		"data/sub/ending.ass": &fstest.MapFile{Data: []byte(assContent)},
	}

	refs := []lexar.SubtitleRef{
		{SubPath: `sub\ending.ass`, CharacterID: "00", AudioID: "end_test", Episode: 8},
	}
	loader := New(fs, testParseFunc(nil, refs))

	result := loader.Load("en", "data/test.file")

	if len(result) != 2 {
		t.Fatalf("expected 2 subtitle quotes, got %d", len(result))
	}
	if result[0].AudioID != "end_test_s0" {
		t.Errorf("quote 0 audioID: got %q, want %q", result[0].AudioID, "end_test_s0")
	}
	if result[1].AudioID != "end_test_s1" {
		t.Errorf("quote 1 audioID: got %q, want %q", result[1].AudioID, "end_test_s1")
	}
	if result[0].Episode != 8 {
		t.Errorf("quote 0 episode: got %d, want 8", result[0].Episode)
	}
}

func TestLoad_SubtitleRefMissingFile(t *testing.T) {
	fs := buildTestFS("data/test.file", []byte("test"))

	refs := []lexar.SubtitleRef{
		{SubPath: `sub\missing.ass`, CharacterID: "00", AudioID: "end_test", Episode: 8},
	}
	loader := New(fs, testParseFunc(nil, refs))

	result := loader.Load("en", "data/test.file")

	if len(result) != 0 {
		t.Errorf("expected 0 quotes when subtitle file missing, got %d", len(result))
	}
}

func TestLoad_CombinesParsedAndSubtitleQuotes(t *testing.T) {
	assContent := "[Script Info]\nTitle: Test\n\n[V4+ Styles]\nFormat: Name\nStyle: Default\n\n[Events]\nFormat: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\nDialogue: 0,0:00:00.00,0:00:05.00,Default,,0,0,0,,Sub line.\n"

	fs := fstest.MapFS{
		"data/test.file":      &fstest.MapFile{Data: encodeTestPayload([]byte("test"))},
		"data/sub/ending.ass": &fstest.MapFile{Data: []byte(assContent)},
	}

	parsed := []dto.ParsedQuote{
		{Text: "Parsed quote", CharacterID: "10", Episode: 1},
	}
	refs := []lexar.SubtitleRef{
		{SubPath: `sub\ending.ass`, CharacterID: "00", AudioID: "end_test", Episode: 8},
	}
	loader := New(fs, testParseFunc(parsed, refs))

	result := loader.Load("en", "data/test.file")

	if len(result) != 2 {
		t.Fatalf("expected 2 total quotes (1 parsed + 1 subtitle), got %d", len(result))
	}
	if result[0].Text != "Parsed quote" {
		t.Errorf("quote 0: got %q, want %q", result[0].Text, "Parsed quote")
	}
	if result[1].Text != "Sub line." {
		t.Errorf("quote 1: got %q, want %q", result[1].Text, "Sub line.")
	}
}
