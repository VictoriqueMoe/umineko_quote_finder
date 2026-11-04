package memory

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func matchText(text string, query string, exact bool) bool {
	if exact {
		return containsWord(text, query)
	}
	return strings.Contains(text, query)
}

func containsWord(text string, word string) bool {
	if word == "" {
		return true
	}
	start := 0
	for {
		rel := strings.Index(text[start:], word)
		if rel < 0 {
			return false
		}
		idx := start + rel
		if isWordBoundary(text, idx, idx+len(word)) {
			return true
		}
		start = idx + 1
	}
}

func isWordBoundary(text string, lo int, hi int) bool {
	if lo > 0 {
		r, _ := utf8.DecodeLastRuneInString(text[:lo])
		if isWordRune(r) {
			return false
		}
	}
	if hi < len(text) {
		r, _ := utf8.DecodeRuneInString(text[hi:])
		if isWordRune(r) {
			return false
		}
	}
	return true
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
