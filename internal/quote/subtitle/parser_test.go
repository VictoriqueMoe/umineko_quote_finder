package subtitle

import "testing"

func TestParseASS_BasicDialogue(t *testing.T) {
	input := `[Script Info]
Title: Test

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 0,0:00:00.73,0:00:02.60,Speech,,0,0,0,,Welcome back, sir.
Dialogue: 0,0:00:02.60,0:00:04.40,Battler,,0,0,0,,Sorry. It took me a while.
`

	entries := ParseASS([]byte(input))

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].Style != "Speech" {
		t.Errorf("entry 0 style: got %q, want 'Speech'", entries[0].Style)
	}
	if entries[0].Text != "Welcome back, sir." {
		t.Errorf("entry 0 text: got %q, want 'Welcome back, sir.'", entries[0].Text)
	}

	if entries[1].Style != "Battler" {
		t.Errorf("entry 1 style: got %q, want 'Battler'", entries[1].Style)
	}
	if entries[1].Text != "Sorry. It took me a while." {
		t.Errorf("entry 1 text: got %q, want 'Sorry. It took me a while.'", entries[1].Text)
	}
}

func TestParseASS_StripsOverrideTags(t *testing.T) {
	input := `[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 0,0:00:00.73,0:00:02.60,Speech,,0,0,0,,{\blur14}{\fad(300,0)}Welcome back, sir.
`

	entries := ParseASS([]byte(input))

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Text != "Welcome back, sir." {
		t.Errorf("text: got %q, want 'Welcome back, sir.'", entries[0].Text)
	}
}

func TestParseASS_StripsInlineFontChanges(t *testing.T) {
	input := `[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 0,0:00:25.40,0:00:26.33,Speech,,0,0,0,,{\blur14}Lord B{\fs29}ATTLER{\fs40}{\blur0\shad0\bord0}.
`

	entries := ParseASS([]byte(input))

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Text != "Lord BATTLER." {
		t.Errorf("text: got %q, want 'Lord BATTLER.'", entries[0].Text)
	}
}

func TestParseASS_HardSpace(t *testing.T) {
	input := `[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 0,0:01:07.63,0:01:11.20,Battler,,0,0,0,,{\pos(626,633)}\hSorry. It took me a while to get here.\h
`

	entries := ParseASS([]byte(input))

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Text != "Sorry. It took me a while to get here." {
		t.Errorf("text: got %q, want 'Sorry. It took me a while to get here.'", entries[0].Text)
	}
}

func TestParseASS_SkipsEmptyLines(t *testing.T) {
	input := `[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 0,0:00:00.00,0:00:01.00,Speech,,0,0,0,,{\blur14}
Dialogue: 0,0:00:01.00,0:00:02.00,Speech,,0,0,0,,Real text here.
`

	entries := ParseASS([]byte(input))

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Text != "Real text here." {
		t.Errorf("text: got %q", entries[0].Text)
	}
}

func TestParseASS_StopsAtNextSection(t *testing.T) {
	input := `[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 0,0:00:00.00,0:00:01.00,Speech,,0,0,0,,First line.

[Extra Section]
Dialogue: 0,0:00:00.00,0:00:01.00,Speech,,0,0,0,,Should not appear.
`

	entries := ParseASS([]byte(input))

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestParseASS_NewlineChars(t *testing.T) {
	input := `[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 0,0:00:00.00,0:00:02.00,Speech,,0,0,0,,Line one\NLine two
Dialogue: 0,0:00:02.00,0:00:04.00,Speech,,0,0,0,,Soft\nwrap here
`

	entries := ParseASS([]byte(input))

	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Text != "Line one Line two" {
		t.Errorf("entry 0 text: got %q, want 'Line one Line two'", entries[0].Text)
	}
	if entries[1].Text != "Soft wrap here" {
		t.Errorf("entry 1 text: got %q, want 'Soft wrap here'", entries[1].Text)
	}
}

func TestParseASS_CommasInText(t *testing.T) {
	input := `[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 0,0:00:00.00,0:00:02.00,Speech,,0,0,0,,Hello, how are you, sir?
`

	entries := ParseASS([]byte(input))

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Text != "Hello, how are you, sir?" {
		t.Errorf("text: got %q, want 'Hello, how are you, sir?'", entries[0].Text)
	}
}

func TestParseASS_MultipleOverrideTags(t *testing.T) {
	input := `[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 0,0:00:00.00,0:00:02.00,Speech,,0,0,0,,{\an8}{\pos(960,100)}{\fad(500,500)}Centered text.
`

	entries := ParseASS([]byte(input))

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Text != "Centered text." {
		t.Errorf("text: got %q, want 'Centered text.'", entries[0].Text)
	}
}

func TestParseASS_NoEventsSection(t *testing.T) {
	input := `[Script Info]
Title: No events
ScriptType: v4.00+
`

	entries := ParseASS([]byte(input))

	if len(entries) != 0 {
		t.Errorf("expected 0 entries for no events section, got %d", len(entries))
	}
}

func TestParseASS_WindowsLineEndings(t *testing.T) {
	input := "[Events]\r\nFormat: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text\r\nDialogue: 0,0:00:00.00,0:00:02.00,Speech,,0,0,0,,Windows line endings.\r\n"

	entries := ParseASS([]byte(input))

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Text != "Windows line endings." {
		t.Errorf("text: got %q, want 'Windows line endings.'", entries[0].Text)
	}
}
