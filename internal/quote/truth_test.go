package quote

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

	for i := 0; i < len(tests); i++ {
		got := tr.Parse(tests[i].input)
		if got != tests[i].want {
			t.Errorf("Parse(%q): got %q, want %q", tests[i].input, got, tests[i].want)
		}
	}
}
