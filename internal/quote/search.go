package quote

import (
	"runtime"
	"strings"
	"sync"
)

func concurrentExactSearch(indices []int, lowerTexts []string, quotes []ParsedQuote, queryLower string, matchesFilter func(ParsedQuote) bool) []SearchResult {
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

	resultSlices := make([][]SearchResult, len(chunks))
	var wg sync.WaitGroup

	for w, c := range chunks {
		wg.Go(func() {
			var local []SearchResult
			for j := c.start; j < c.end; j++ {
				idx := indices[j]
				if strings.Contains(lowerTexts[idx], queryLower) {
					if matchesFilter(quotes[idx]) {
						local = append(local, NewSearchResult(quotes[idx], 100))
					}
				}
			}
			resultSlices[w] = local
		})
	}

	wg.Wait()

	var merged []SearchResult
	for _, s := range resultSlices {
		merged = append(merged, s...)
	}
	return merged
}
