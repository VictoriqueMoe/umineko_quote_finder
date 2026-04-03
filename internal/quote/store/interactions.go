package store

import (
	"slices"
)

func InteractionPairKey(charA string, charB string) (string, bool) {
	if charA == "" || charB == "" || charA == charB {
		return "", false
	}
	if charA > charB {
		charA, charB = charB, charA
	}
	return charA + "|" + charB, true
}

func buildInteractionIndex(characterIDs []string, episodes []int) map[string][]int {
	sets := make(map[string]map[int]any)

	var prevCharID string
	var prevEpisode int
	var prevIndex int

	for i := 0; i < len(characterIDs); i++ {
		if characterIDs[i] == "narrator" {
			prevCharID = ""
			continue
		}

		if prevCharID != "" && prevEpisode == episodes[i] {
			key, ok := InteractionPairKey(prevCharID, characterIDs[i])
			if ok {
				if sets[key] == nil {
					sets[key] = make(map[int]any)
				}
				sets[key][prevIndex] = struct{}{}
				sets[key][i] = struct{}{}
			}
		}

		prevCharID = characterIDs[i]
		prevEpisode = episodes[i]
		prevIndex = i
	}

	result := make(map[string][]int, len(sets))
	for key, set := range sets {
		indices := make([]int, 0, len(set))
		for idx := range set {
			indices = append(indices, idx)
		}
		slices.Sort(indices)
		result[key] = indices
	}

	return result
}
