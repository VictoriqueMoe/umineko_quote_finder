package memory

import (
	"runtime"
	"sync"
)

func concurrentMatch(
	indices []int,
	lowerTexts []string,
	queryLower string,
	exact bool,
	keep func(int) bool,
) []int {
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
		end := min(i+chunkSize, total)
		chunks = append(chunks, chunk{i, end})
	}

	resultSlices := make([][]int, len(chunks))
	var wg sync.WaitGroup

	for w, c := range chunks {
		wg.Go(func() {
			var local []int
			for _, idx := range indices[c.start:c.end] {
				if matchText(lowerTexts[idx], queryLower, exact) && keep(idx) {
					local = append(local, idx)
				}
			}
			resultSlices[w] = local
		})
	}

	wg.Wait()

	var merged []int
	for _, s := range resultSlices {
		merged = append(merged, s...)
	}
	return merged
}
