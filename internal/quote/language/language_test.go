package language

import "testing"

func TestParse_AllLanguages(t *testing.T) {
	cases := []struct {
		input string
		want  Language
	}{
		{"auto", Auto},
		{"en", English},
		{"wh", WitchHunt},
		{"ja", Japanese},
		{"ru", Russian},
		{"es", Spanish},
		{"pt", Portuguese},
	}
	for _, tc := range cases {
		got := Auto.Parse(tc.input)
		if got != tc.want {
			t.Errorf("Parse(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParse_UnknownDefaultsToEnglish(t *testing.T) {
	cases := []string{"fr", "de", "", "unknown"}
	for _, input := range cases {
		got := Auto.Parse(input)
		if got != English {
			t.Errorf("Parse(%q) = %q, want %q", input, got, English)
		}
	}
}

func TestAll_DoesNotContainAuto(t *testing.T) {
	for _, lang := range All {
		if lang == Auto {
			t.Error("All should not contain Auto")
		}
	}
}

func TestAll_ContainsAllRealLanguages(t *testing.T) {
	expected := []Language{English, WitchHunt, Japanese, Russian, Spanish, Portuguese}
	if len(All) != len(expected) {
		t.Fatalf("All has %d languages, want %d", len(All), len(expected))
	}
	for i, lang := range expected {
		if All[i] != lang {
			t.Errorf("All[%d] = %q, want %q", i, All[i], lang)
		}
	}
}

func TestString(t *testing.T) {
	if English.String() != "en" {
		t.Errorf("English.String() = %q, want %q", English.String(), "en")
	}
	if Auto.String() != "auto" {
		t.Errorf("Auto.String() = %q, want %q", Auto.String(), "auto")
	}
}
