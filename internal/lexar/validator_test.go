package lexar

import (
	"strings"
	"testing"
)

func TestValidate_ValidScript(t *testing.T) {
	input := `preset_define 1,1,-1,#FF0000,0,0,0,0,0
new_episode 3
d [lv 0*"10"*"30100001"]` + "`\"This is valid dialogue.\"`" + `[\]`

	script := Parse(input)
	errors := Validate(script)

	if len(errors) != 0 {
		t.Errorf("expected no errors for valid script, got %d:", len(errors))
		for _, err := range errors {
			t.Errorf("  %s", err)
		}
	}
}

func TestValidate_UnknownFormatTag(t *testing.T) {
	input := `d [lv 0*"10"*"10100001"]` + "`\"{bogus:some content}\"`" + `[\]`

	script := Parse(input)
	errors := Validate(script)

	found := false
	for _, err := range errors {
		if strings.Contains(err.Message, "unknown format tag") && strings.Contains(err.Message, "bogus") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected unknown format tag error for 'bogus', got %v", errors)
	}
}

func TestValidate_KnownFormatTags(t *testing.T) {
	tags := []string{"p", "c", "i", "ruby", "h", "y"}

	for _, tag := range tags {
		t.Run(tag, func(t *testing.T) {
			input := `d [lv 0*"10"*"10100001"]` + "`\"{" + tag + ":param:content}\"`" + `[\]`
			script := Parse(input)
			errors := Validate(script)

			for _, err := range errors {
				if strings.Contains(err.Message, "unknown format tag") {
					t.Errorf("tag %q should be known, got error: %s", tag, err)
				}
			}
		})
	}
}

func TestValidate_VoiceCommandMissingFields(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantMsg string
	}{
		{
			name:    "missing character and audio",
			input:   `d [lv 0]` + "`\"Hello.\"`" + `[\]`,
			wantMsg: "missing character ID",
		},
		{
			name:    "missing audio only",
			input:   `d [lv 0*"10"]` + "`\"Hello.\"`" + `[\]`,
			wantMsg: "missing audio ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := Parse(tt.input)
			errors := Validate(script)

			found := false
			for _, err := range errors {
				if strings.Contains(err.Message, tt.wantMsg) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected error containing %q, got %v", tt.wantMsg, errors)
			}
		})
	}
}

func TestValidate_EpisodeMarkerMissingNumber(t *testing.T) {
	input := "new_episode"

	script := Parse(input)
	errors := Validate(script)

	found := false
	for _, err := range errors {
		if strings.Contains(err.Message, "episode marker") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected episode marker error, got %v", errors)
	}
}

func TestValidate_NestedUnknownTag(t *testing.T) {
	input := `d [lv 0*"10"*"10100001"]` + "`\"{p:1:{fake:nested content}}\"`" + `[\]`

	script := Parse(input)
	errors := Validate(script)

	found := false
	for _, err := range errors {
		if strings.Contains(err.Message, "unknown format tag") && strings.Contains(err.Message, "fake") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected unknown format tag error for nested 'fake', got %v", errors)
	}
}

func TestValidate_RealGameScript(t *testing.T) {
	input := "*o1_1\n" +
		`d [lv 0*"04"*"10200442"]` + "`\"KyaaaaaAAAAAAaaaaaAAaa!!!\"`" + `[\]` + "\n" +
		"*o1_5\n" +
		`d [lv 0*"05"*"11000358"]` + "`\"George, take everyone and return to the mansion!!!`[@][lv 0*\"05\"*\"11000359\"]` Quickly!!`[@][lv 0*\"05\"*\"11000360\"]` Right now!!\"`" + `[\]`

	script := Parse(input)
	errors := Validate(script)

	if len(errors) != 0 {
		t.Errorf("real game script should have no validation errors, got %d:", len(errors))
		for _, err := range errors {
			t.Errorf("  %s", err)
		}
	}
}
