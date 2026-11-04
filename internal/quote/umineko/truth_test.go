package umineko

import "testing"

func TestTruthParse(t *testing.T) {
	var tr Truth

	tests := []struct {
		input string
		want  Truth
	}{
		{"red", TruthRed},
		{"blue", TruthBlue},
		{"", TruthAll},
		{"unknown", TruthAll},
		{"RED", TruthAll},
		{"Blue", TruthAll},
	}

	for _, tc := range tests {
		got := tr.Parse(tc.input)
		if got != tc.want {
			t.Errorf("Parse(%q): got %q, want %q", tc.input, got, tc.want)
		}
	}
}
