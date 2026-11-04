package umineko

import (
	"fmt"
	"maps"
	"strconv"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/quotes"

	"github.com/VictoriqueMoe/umineko_script_parser/umineko/character"
)

type (
	uminekoStatsComputer struct {
		cached    *dto.StatsResult
		byEpisode map[int]*dto.StatsResult
		fallback  *dto.StatsResult
	}

	uminekoTallies struct {
		charCounts   map[string]int
		charEpCounts map[string]map[int]int
		epTruth      map[int]params.TruthCounts
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

func NewUminekoStats(t params.SpeakerTallies, truth map[int]params.TruthCounts) (quotes.Stats, error) {
	all, err := allEpisodeTallies(t, truth)
	if err != nil {
		return nil, err
	}

	s := &uminekoStatsComputer{
		cached:    buildStats(all, quotes.AllEpisodes),
		byEpisode: make(map[int]*dto.StatsResult, len(uminekoEpisodeNames)),
		fallback:  buildStats(episodeTallies(t, -1), -1),
	}

	for ep := 1; ep <= len(uminekoEpisodeNames); ep++ {
		s.byEpisode[ep] = buildStats(episodeTallies(t, ep), ep)
	}

	return s, nil
}

func allEpisodeTallies(t params.SpeakerTallies, truth map[int]params.TruthCounts) (uminekoTallies, error) {
	charEpCounts := make(map[string]map[int]int, len(t.CharacterSegments))
	for id, segments := range t.CharacterSegments {
		charEpCounts[id] = make(map[int]int, len(segments))
		for segment, count := range segments {
			ep, err := strconv.Atoi(segment)
			if err != nil {
				return uminekoTallies{}, fmt.Errorf("umineko segment %q is not an episode: %w", segment, err)
			}
			charEpCounts[id][ep] = count
		}
	}

	return uminekoTallies{
		charCounts:   t.CharacterTotals(),
		charEpCounts: charEpCounts,
		epTruth:      truth,
		interactions: t.InteractionTotals(),
	}, nil
}

func episodeTallies(t params.SpeakerTallies, episode int) uminekoTallies {
	segment := strconv.Itoa(episode)

	charCounts := make(map[string]int)
	for id, segments := range t.CharacterSegments {
		if count, ok := segments[segment]; ok {
			charCounts[id] = count
		}
	}

	interactions := make(map[string]int, len(t.SegmentInteractions[segment]))
	maps.Copy(interactions, t.SegmentInteractions[segment])

	return uminekoTallies{
		charCounts:   charCounts,
		charEpCounts: make(map[string]map[int]int),
		epTruth:      make(map[int]params.TruthCounts),
		interactions: interactions,
	}
}

func (s *uminekoStatsComputer) Compute(episode int) any {
	if episode == quotes.AllEpisodes {
		return s.cached
	}
	if result, ok := s.byEpisode[episode]; ok {
		return result
	}
	return s.fallback
}

func resolveName(id string) string {
	return character.CharacterNames.GetCharacterName(character.CharacterFromID(id))
}

func buildStats(t uminekoTallies, episode int) *dto.StatsResult {
	ranked := quotes.RankCharacters(t.charCounts)

	result := &dto.StatsResult{
		TopSpeakers:       quotes.TopSpeakers(ranked, 20, resolveName),
		Interactions:      quotes.TopInteractions(t.interactions, 25, resolveName),
		InteractionCounts: t.interactions,
		CharacterNames:    quotes.BuildNameMap(t.charCounts, resolveName),
		EpisodeNames:      uminekoEpisodeNames,
	}

	if episode == quotes.AllEpisodes {
		result.LinesPerEpisode = uminekoLinesPerEpisode(t.charEpCounts, ranked, 10)
		result.TruthPerEpisode = uminekoTruthPerEpisode(t.epTruth)
		result.CharacterPresence = uminekoBuildCharacterPresence(ranked, t.charEpCounts, 12)
	}

	return result
}

func uminekoLinesPerEpisode(charEpCounts map[string]map[int]int, ranked []quotes.RankedChar, topN int) []dto.EpisodeCharacterLines {
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

func uminekoTruthPerEpisode(epTruth map[int]params.TruthCounts) []dto.EpisodeTruth {
	result := make([]dto.EpisodeTruth, 8)
	for ep := 1; ep <= 8; ep++ {
		counts := epTruth[ep]
		result[ep-1] = dto.EpisodeTruth{
			Episode: ep,
			Red:     counts.Red,
			Blue:    counts.Blue,
			Gold:    counts.Gold,
			Purple:  counts.Purple,
		}
	}
	return result
}

func uminekoBuildCharacterPresence(ranked []quotes.RankedChar, charEpCounts map[string]map[int]int, n int) []dto.CharacterPresence {
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
