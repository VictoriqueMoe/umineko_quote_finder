package store

import (
	"cmp"
	"slices"
	"strings"

	"umineko_quote/internal/dto"
)

type Stats interface {
	Compute(episode int) any
}

const AllEpisodes = 0

type RankedChar struct {
	ID    string
	Count int
}

func RankCharacters(charCounts map[string]int) []RankedChar {
	ranked := make([]RankedChar, 0, len(charCounts))
	for id, count := range charCounts {
		ranked = append(ranked, RankedChar{id, count})
	}
	slices.SortFunc(ranked, func(a, b RankedChar) int {
		return cmp.Compare(b.Count, a.Count)
	})
	return ranked
}

func TopSpeakers(ranked []RankedChar, n int, resolveName func(string) string) []dto.SpeakerStat {
	if len(ranked) < n {
		n = len(ranked)
	}
	result := make([]dto.SpeakerStat, n)
	for i := 0; i < n; i++ {
		result[i] = dto.SpeakerStat{
			CharacterID: ranked[i].ID,
			Name:        resolveName(ranked[i].ID),
			Count:       ranked[i].Count,
		}
	}
	return result
}

func TopInteractions(interactionCounts map[string]int, n int, resolveName func(string) string) []dto.InteractionPair {
	type pairCount struct {
		key   string
		count int
	}

	sorted := make([]pairCount, 0, len(interactionCounts))
	for key, count := range interactionCounts {
		sorted = append(sorted, pairCount{key, count})
	}

	slices.SortFunc(sorted, func(a, b pairCount) int {
		return cmp.Compare(b.count, a.count)
	})

	if len(sorted) < n {
		n = len(sorted)
	}

	result := make([]dto.InteractionPair, n)
	for i := 0; i < n; i++ {
		parts := strings.SplitN(sorted[i].key, "|", 2)
		result[i] = dto.InteractionPair{
			CharA: parts[0],
			CharB: parts[1],
			NameA: resolveName(parts[0]),
			NameB: resolveName(parts[1]),
			Count: sorted[i].count,
		}
	}
	return result
}

func BuildNameMap(charCounts map[string]int, resolveName func(string) string) map[string]string {
	nameMap := make(map[string]string, len(charCounts))
	for id := range charCounts {
		nameMap[id] = resolveName(id)
	}
	return nameMap
}
