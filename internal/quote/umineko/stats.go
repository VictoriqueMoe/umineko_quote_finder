package umineko

import (
	"umineko_quote/internal/quote/store"

	"github.com/VictoriqueMoe/umineko_script_parser/umineko/character"

	"umineko_quote/internal/dto"
)

type (
	uminekoStatsComputer struct {
		quotes []dto.UminekoQuote
		cached *dto.StatsResult
	}

	uminekoTallies struct {
		charCounts   map[string]int
		charEpCounts map[string]map[int]int
		epTruth      map[int][4]int
		interactions map[string]int
	}
)

var uminekoEpisodeNames = map[int]string{
	1: "Legend",
	2: "Turn",
	3: "Banquet",
	4: "Alliance",
	5: "End",
	6: "Dawn",
	7: "Requiem",
	8: "Twilight",
}

func NewUminekoStats(quotes []dto.UminekoQuote) store.Stats {
	s := &uminekoStatsComputer{quotes: quotes}
	s.cached = s.compute(store.AllEpisodes)
	return s
}

func (s *uminekoStatsComputer) Compute(episode int) any {
	if episode != store.AllEpisodes {
		return s.compute(episode)
	}
	return s.cached
}

func resolveName(id string) string {
	return character.CharacterNames.GetCharacterName(character.CharacterFromID(id))
}

func (s *uminekoStatsComputer) compute(episode int) *dto.StatsResult {
	t := s.tally(episode)
	ranked := store.RankCharacters(t.charCounts)

	result := &dto.StatsResult{
		TopSpeakers:       store.TopSpeakers(ranked, 20, resolveName),
		Interactions:      store.TopInteractions(t.interactions, 25, resolveName),
		InteractionCounts: t.interactions,
		CharacterNames:    store.BuildNameMap(t.charCounts, resolveName),
		EpisodeNames:      uminekoEpisodeNames,
	}

	if episode == store.AllEpisodes {
		result.LinesPerEpisode = uminekoLinesPerEpisode(t.charEpCounts, ranked, 10)
		result.TruthPerEpisode = uminekoTruthPerEpisode(t.epTruth)
		result.CharacterPresence = uminekoBuildCharacterPresence(ranked, t.charEpCounts, 12)
	}

	return result
}

func (s *uminekoStatsComputer) tally(episode int) uminekoTallies {
	t := uminekoTallies{
		charCounts:   make(map[string]int),
		charEpCounts: make(map[string]map[int]int),
		epTruth:      make(map[int][4]int),
		interactions: make(map[string]int),
	}

	var prevCharID string
	var prevEpisode int

	for _, q := range s.quotes {
		if episode != store.AllEpisodes && q.Episode != episode {
			prevCharID = ""
			continue
		}

		if q.HasRedTruth {
			counts := t.epTruth[q.Episode]
			counts[0]++
			t.epTruth[q.Episode] = counts
		}
		if q.HasBlueTruth {
			counts := t.epTruth[q.Episode]
			counts[1]++
			t.epTruth[q.Episode] = counts
		}
		if q.HasGoldTruth {
			counts := t.epTruth[q.Episode]
			counts[2]++
			t.epTruth[q.Episode] = counts
		}
		if q.HasPurpleTruth {
			counts := t.epTruth[q.Episode]
			counts[3]++
			t.epTruth[q.Episode] = counts
		}

		if q.CharacterID == "narrator" {
			prevCharID = ""
			continue
		}

		t.charCounts[q.CharacterID]++

		if t.charEpCounts[q.CharacterID] == nil {
			t.charEpCounts[q.CharacterID] = make(map[int]int)
		}
		t.charEpCounts[q.CharacterID][q.Episode]++

		if prevCharID != "" && prevCharID != q.CharacterID && prevEpisode == q.Episode {
			if key, ok := store.InteractionPairKey(prevCharID, q.CharacterID); ok {
				t.interactions[key]++
			}
		}

		prevCharID = q.CharacterID
		prevEpisode = q.Episode
	}

	return t
}

func uminekoLinesPerEpisode(charEpCounts map[string]map[int]int, ranked []store.RankedChar, topN int) []dto.EpisodeCharacterLines {
	if len(ranked) < topN {
		topN = len(ranked)
	}
	topSet := make(map[string]bool, topN)
	for i := 0; i < topN; i++ {
		topSet[ranked[i].ID] = true
	}

	result := make([]dto.EpisodeCharacterLines, 8)
	for ep := 1; ep <= 8; ep++ {
		chars := make(map[string]int)
		for id, epMap := range charEpCounts {
			if epMap[ep] > 0 {
				if topSet[id] {
					chars[id] = epMap[ep]
				} else {
					chars["other"] += epMap[ep]
				}
			}
		}
		result[ep-1] = dto.EpisodeCharacterLines{
			Episode:     ep,
			EpisodeName: uminekoEpisodeNames[ep],
			Characters:  chars,
		}
	}
	return result
}

func uminekoTruthPerEpisode(epTruth map[int][4]int) []dto.EpisodeTruth {
	result := make([]dto.EpisodeTruth, 8)
	for ep := 1; ep <= 8; ep++ {
		counts := epTruth[ep]
		result[ep-1] = dto.EpisodeTruth{
			Episode: ep,
			Red:     counts[0],
			Blue:    counts[1],
			Gold:    counts[2],
			Purple:  counts[3],
		}
	}
	return result
}

func uminekoBuildCharacterPresence(ranked []store.RankedChar, charEpCounts map[string]map[int]int, n int) []dto.CharacterPresence {
	if len(ranked) < n {
		n = len(ranked)
	}
	result := make([]dto.CharacterPresence, n)
	for i := 0; i < n; i++ {
		id := ranked[i].ID
		episodes := make([]int, 8)
		for ep := 1; ep <= 8; ep++ {
			episodes[ep-1] = charEpCounts[id][ep]
		}
		result[i] = dto.CharacterPresence{
			CharacterID: id,
			Name:        resolveName(id),
			Episodes:    episodes,
		}
	}
	return result
}
