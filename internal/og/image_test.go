package og

import (
	"bytes"
	"strings"
	"testing"

	"umineko_quote/internal/quote/language"
)

func longText(n int) string {
	var b strings.Builder
	for i := range n {
		b.WriteByte(byte('a' + (i % 26)))
	}
	return b.String()
}

func TestTruncateText(t *testing.T) {
	short := "hello"
	if got := truncateText(short, 300); got != short {
		t.Errorf("short text was modified: %q", got)
	}

	long := longText(400)
	got := truncateText(long, 300)
	runes := []rune(got)
	if len(runes) != 300 {
		t.Errorf("expected 300 runes, got %d", len(runes))
	}
	if !strings.HasSuffix(got, "...") {
		t.Error("truncated text should end with ...")
	}
}

func TestTruncateTextNoOpForShort(t *testing.T) {
	text := longText(299)
	if got := truncateText(text, 300); got != text {
		t.Error("text under limit should not be modified")
	}

	exact := longText(300)
	if got := truncateText(exact, 300); got != exact {
		t.Error("text at exact limit should not be modified")
	}
}

func TestTruncateSegments(t *testing.T) {
	short := []textSegment{{Text: "hello", Color: creamColor}}
	if got := truncateSegments(short, 300); len(got) != 1 || got[0].Text != "hello" {
		t.Errorf("short segments were modified: %v", got)
	}

	long := []textSegment{{Text: longText(400), Color: creamColor}}
	got := truncateSegments(long, 300)

	total := 0
	for _, seg := range got {
		total += len([]rune(seg.Text))
	}
	if total != 300 {
		t.Errorf("expected 300 total runes, got %d", total)
	}

	lastSeg := got[len(got)-1]
	if !strings.HasSuffix(lastSeg.Text, "...") {
		t.Error("truncated segments should end with ...")
	}
}

func TestTruncateSegmentsMultiple(t *testing.T) {
	segs := []textSegment{
		{Text: longText(200), Color: creamColor},
		{Text: longText(200), Color: redTruthColor},
	}
	got := truncateSegments(segs, 300)

	total := 0
	for _, seg := range got {
		total += len([]rune(seg.Text))
	}
	if total != 300 {
		t.Errorf("expected 300 total runes across segments, got %d", total)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 segments, got %d", len(got))
	}
}

func TestTruncateSegmentsNoOpForShort(t *testing.T) {
	segs := []textSegment{
		{Text: "hello", Color: creamColor},
		{Text: " world", Color: redTruthColor},
	}
	got := truncateSegments(segs, 300)
	if len(got) != 2 || got[0].Text != "hello" || got[1].Text != " world" {
		t.Error("short segments should not be modified")
	}
}

func TestGenerateFullWithHTMLDoesNotTruncate(t *testing.T) {
	gen := NewImageGenerator()
	html := longText(400)

	truncated, err := gen.Generate("test_html_trunc", language.English, "", html, "Beatrice", false, "Test Brand", "Episode 1")
	if err != nil {
		t.Fatalf("Generate(full=false) failed: %v", err)
	}

	full, err := gen.Generate("test_html_trunc", language.English, "", html, "Beatrice", true, "Test Brand", "Episode 1")
	if err != nil {
		t.Fatalf("Generate(full=true) failed: %v", err)
	}

	if bytes.Equal(truncated, full) {
		t.Error("expected full and truncated HTML images to differ")
	}
}

func TestGenerateCachesSeparately(t *testing.T) {
	gen := NewImageGenerator()
	html := longText(400)

	first, err := gen.Generate("test_cache", language.English, "", html, "Beatrice", false, "Test Brand", "Episode 1")
	if err != nil {
		t.Fatalf("first Generate failed: %v", err)
	}

	second, err := gen.Generate("test_cache", language.English, "", html, "Beatrice", false, "Test Brand", "Episode 1")
	if err != nil {
		t.Fatalf("second Generate failed: %v", err)
	}

	if !bytes.Equal(first, second) {
		t.Error("same params should return cached (identical) result")
	}

	full, err := gen.Generate("test_cache", language.English, "", html, "Beatrice", true, "Test Brand", "Episode 1")
	if err != nil {
		t.Fatalf("full Generate failed: %v", err)
	}

	if bytes.Equal(first, full) {
		t.Error("full=true should not return the truncated cached image")
	}
}

func TestGenerateShortQuoteIdenticalFullOrNot(t *testing.T) {
	gen := NewImageGenerator()
	html := "Short quote"

	truncated, err := gen.Generate("test_short", language.English, "", html, "Beatrice", false, "Test Brand", "Episode 1")
	if err != nil {
		t.Fatalf("Generate(full=false) failed: %v", err)
	}

	full, err := gen.Generate("test_short", language.English, "", html, "Beatrice", true, "Test Brand", "Episode 1")
	if err != nil {
		t.Fatalf("Generate(full=true) failed: %v", err)
	}

	if !bytes.Equal(truncated, full) {
		t.Error("short quote should produce identical images regardless of full flag")
	}
}
