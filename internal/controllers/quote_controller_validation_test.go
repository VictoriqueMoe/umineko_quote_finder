package controllers

import "testing"

func TestValidateInteractionQueryParams(t *testing.T) {
	tests := []struct {
		name         string
		character    string
		interactionA string
		interactionB string
		wantInvalid  bool
		wantError    string
	}{
		{
			name:         "valid no interaction",
			character:    "battler",
			interactionA: "",
			interactionB: "",
			wantInvalid:  false,
		},
		{
			name:         "valid interaction pair",
			character:    "",
			interactionA: "beatrice",
			interactionB: "dlanor",
			wantInvalid:  false,
		},
		{
			name:         "invalid only interactionA",
			character:    "",
			interactionA: "beatrice",
			interactionB: "",
			wantInvalid:  true,
			wantError:    "interactionA and interactionB must both be provided",
		},
		{
			name:         "invalid only interactionB",
			character:    "",
			interactionA: "",
			interactionB: "dlanor",
			wantInvalid:  true,
			wantError:    "interactionA and interactionB must both be provided",
		},
		{
			name:         "invalid character with interaction pair",
			character:    "battler",
			interactionA: "beatrice",
			interactionB: "dlanor",
			wantInvalid:  true,
			wantError:    "character cannot be combined with interactionA/interactionB",
		},
		{
			name:         "invalid character with interaction pair and whitespace",
			character:    " battler ",
			interactionA: " beatrice ",
			interactionB: " dlanor ",
			wantInvalid:  true,
			wantError:    "character cannot be combined with interactionA/interactionB",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			errMsg, invalid := validateInteractionQueryParams(tc.character, tc.interactionA, tc.interactionB)
			if invalid != tc.wantInvalid {
				t.Fatalf("invalid = %v, want %v (err: %q)", invalid, tc.wantInvalid, errMsg)
			}
			if errMsg != tc.wantError {
				t.Fatalf("errMsg = %q, want %q", errMsg, tc.wantError)
			}
		})
	}
}
