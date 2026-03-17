package quote

import (
	"slices"

	"umineko_quote/internal/dto"
)

func interactionPairKey(charA string, charB string) (string, bool) {
	if charA == "" || charB == "" || charA == charB {
		return "", false
	}
	if charA > charB {
		charA, charB = charB, charA
	}
	return charA + "|" + charB, true
}

func buildInteractionQuoteIndex(quotes []dto.ParsedQuote) map[string][]int {
	sets := make(map[string]map[int]any)

	var prevCharID string
	var prevEpisode int
	var prevIndex int

	for i, q := range quotes {
		if q.CharacterID == "narrator" {
			prevCharID = ""
			continue
		}

		if prevCharID != "" && prevEpisode == q.Episode {
			key, ok := interactionPairKey(prevCharID, q.CharacterID)
			if ok {
				if sets[key] == nil {
					sets[key] = make(map[int]any)
				}
				sets[key][prevIndex] = struct{}{}
				sets[key][i] = struct{}{}
			}
		}

		prevCharID = q.CharacterID
		prevEpisode = q.Episode
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
