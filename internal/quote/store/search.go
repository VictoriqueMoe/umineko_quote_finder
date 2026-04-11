package store

import (
	"runtime"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"

	"umineko_quote/internal/dto"
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

func concurrentExactSearchGeneric[Q any](
	indices []int,
	lowerTexts []string,
	quotes []Q,
	queryLower string,
	exact bool,
	matchesFilter func(int) bool,
) []dto.SearchResult {
	numWorkers := runtime.NumCPU()
	total := len(indices)
	if total == 0 {
		return nil
	}
	if numWorkers > total {
		numWorkers = total
	}

	type chunk struct {
		start, end int
	}
	chunkSize := (total + numWorkers - 1) / numWorkers
	chunks := make([]chunk, 0, numWorkers)
	for i := 0; i < total; i += chunkSize {
		end := i + chunkSize
		if end > total {
			end = total
		}
		chunks = append(chunks, chunk{i, end})
	}

	resultSlices := make([][]dto.SearchResult, len(chunks))
	var wg sync.WaitGroup

	for w, c := range chunks {
		wg.Go(func() {
			var local []dto.SearchResult
			for j := c.start; j < c.end; j++ {
				idx := indices[j]
				if matchText(lowerTexts[idx], queryLower, exact) {
					if matchesFilter(idx) {
						local = append(local, dto.SearchResult{Quote: quotes[idx], Score: 100})
					}
				}
			}
			resultSlices[w] = local
		})
	}

	wg.Wait()

	var merged []dto.SearchResult
	for _, s := range resultSlices {
		merged = append(merged, s...)
	}
	return merged
}
