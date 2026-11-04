package memory

import (
	"testing"
)

func TestMatchText(t *testing.T) {
	cases := []struct {
		text  string
		query string
		exact bool
		want  bool
	}{
		{"hello world", "hello", false, true},
		{"hello world", "beatrice", false, false},
		{"i have a brother", "broth", false, true},
		{"i have a brother", "broth", true, false},
		{"hot broth on the table", "broth", true, true},
		{"anything", "", false, true},
		{"anything", "", true, true},
	}
	for _, c := range cases {
		got := matchText(c.text, c.query, c.exact)
		if got != c.want {
			t.Errorf("matchText(%q, %q, %v) = %v, want %v", c.text, c.query, c.exact, got, c.want)
		}
	}
}

func TestContainsWord(t *testing.T) {
	cases := []struct {
		text string
		word string
		want bool
	}{
		{"i have a brother", "broth", false},
		{"hot broth on the table", "broth", true},
		{"broth is good", "broth", true},
		{"all broth", "broth", true},
		{"brothers in arms", "broth", false},
		{"un-broth-like", "broth", true},
		{"café au lait", "café", true},
		{"cafés", "café", false},
		{"真実は一つ", "真実", false},
	}
	for _, c := range cases {
		got := containsWord(c.text, c.word)
		if got != c.want {
			t.Errorf("containsWord(%q, %q) = %v, want %v", c.text, c.word, got, c.want)
		}
	}
}
