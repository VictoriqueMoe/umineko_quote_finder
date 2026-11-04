package memory

import (
	"slices"
	"testing"
)

func keepAll(int) bool {
	return true
}

func TestConcurrentMatch_EmptyIndices(t *testing.T) {
	results := concurrentMatch([]int{}, []string{}, "test", false, keepAll)

	if results != nil {
		t.Errorf("empty indices: got %v, want nil", results)
	}
}

func TestConcurrentMatch_FindsMatches(t *testing.T) {
	lowerTexts := []string{
		"hello world",
		"goodbye world",
		"hello again",
		"something else",
	}

	results := concurrentMatch([]int{0, 1, 2, 3}, lowerTexts, "hello", false, keepAll)

	if !slices.Equal(results, []int{0, 2}) {
		t.Errorf("matches: got %v, want [0 2]", results)
	}
}

func TestConcurrentMatch_RespectsKeep(t *testing.T) {
	characterIDs := []string{"10", "27", "10"}
	lowerTexts := []string{
		"hello world",
		"hello again",
		"hello there",
	}

	results := concurrentMatch([]int{0, 1, 2}, lowerTexts, "hello", false, func(idx int) bool {
		return characterIDs[idx] == "10"
	})

	if !slices.Equal(results, []int{0, 2}) {
		t.Errorf("filtered matches: got %v, want [0 2]", results)
	}
}

func TestConcurrentMatch_NoMatches(t *testing.T) {
	lowerTexts := []string{
		"hello world",
		"goodbye world",
	}

	results := concurrentMatch([]int{0, 1}, lowerTexts, "beatrice", false, keepAll)

	if len(results) != 0 {
		t.Errorf("expected 0 matches, got %d", len(results))
	}
}

func TestConcurrentMatch_SubsetIndices(t *testing.T) {
	lowerTexts := []string{
		"hello world",
		"hello again",
		"hello there",
	}

	results := concurrentMatch([]int{0, 2}, lowerTexts, "hello", false, keepAll)

	if !slices.Equal(results, []int{0, 2}) {
		t.Errorf("subset matches: got %v, want [0 2]", results)
	}
}

func TestConcurrentMatch_ExactMode(t *testing.T) {
	lowerTexts := []string{
		"i have a brother",
		"hot broth on the table",
		"brothers in arms",
	}
	indices := []int{0, 1, 2}

	loose := concurrentMatch(indices, lowerTexts, "broth", false, keepAll)
	if len(loose) != 3 {
		t.Fatalf("loose mode: expected 3 matches, got %d", len(loose))
	}

	exact := concurrentMatch(indices, lowerTexts, "broth", true, keepAll)
	if !slices.Equal(exact, []int{1}) {
		t.Fatalf("exact mode: got %v, want [1]", exact)
	}
}
