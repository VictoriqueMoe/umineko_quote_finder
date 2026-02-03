package quote

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

type (
	Stats interface {
		Compute(episode int) any
	}

	speakerStat struct {
		CharacterID string `json:"characterId"`
		Name        string `json:"name"`
		Count       int    `json:"count"`
	}

	episodeTruth struct {
		Episode int `json:"episode"`
		Red     int `json:"red"`
		Blue    int `json:"blue"`
	}

	interactionPair struct {
		CharA string `json:"charA"`
		CharB string `json:"charB"`
		NameA string `json:"nameA"`
		NameB string `json:"nameB"`
		Count int    `json:"count"`
	}

	episodeCharacterLines struct {
		Episode     int            `json:"episode"`
		EpisodeName string         `json:"episodeName"`
		Characters  map[string]int `json:"characters"`
	}

	characterPresence struct {
		CharacterID string `json:"characterId"`
		Name        string `json:"name"`
		Episodes    []int  `json:"episodes"`
	}

	statsResult struct {
		TopSpeakers       []speakerStat           `json:"topSpeakers"`
		LinesPerEpisode   []episodeCharacterLines `json:"linesPerEpisode"`
		TruthPerEpisode   []episodeTruth          `json:"truthPerEpisode"`
		Interactions      []interactionPair       `json:"interactions"`
		CharacterPresence []characterPresence     `json:"characterPresence"`
		CharacterNames    map[string]string       `json:"characterNames"`
		EpisodeNames      map[int]string          `json:"episodeNames"`
	}

	statsComputer struct {
		quotes []ParsedQuote
		cached *statsResult
	}

	tallies struct {
		charCounts   map[string]int
		charEpCounts map[string]map[int]int
		epTruth      map[int][2]int
		interactions map[string]int
	}

	rankedChar struct {
		id    string
		count int
	}
)

const (
	AllEpisodes = iota
)

var (
	episodeNames = map[int]string{
		1: "Legend",
		2: "Turn",
		3: "Banquet",
		4: "Alliance",
		5: "End",
		6: "Dawn",
		7: "Requiem",
		8: "Twilight",
	}
)

func NewStats(quotes []ParsedQuote) Stats {
	s := &statsComputer{quotes: quotes}
	s.cached = s.compute(AllEpisodes)
	return s
}

func (s *statsComputer) Compute(episode int) any {
	if episode != AllEpisodes {
		return s.compute(episode)
	}
	return s.cached
}

func (s *statsComputer) compute(episode int) *statsResult {
	t := s.tally(episode)
	ranked := s.rankCharacters(t.charCounts)

	result := &statsResult{
		TopSpeakers:    s.topSpeakers(ranked, 20),
		Interactions:   s.topInteractions(t.interactions, 25),
		CharacterNames: s.buildNameMap(t.charCounts),
		EpisodeNames:   episodeNames,
	}

	if episode == AllEpisodes {
		result.LinesPerEpisode = s.linesPerEpisode(t.charEpCounts, ranked, 10)
		result.TruthPerEpisode = s.truthPerEpisode(t.epTruth)
		result.CharacterPresence = s.buildCharacterPresence(ranked, t.charEpCounts, 12)
	}

	return result
}

func (s *statsComputer) tally(episode int) tallies {
	t := tallies{
		charCounts:   make(map[string]int),
		charEpCounts: make(map[string]map[int]int),
		epTruth:      make(map[int][2]int),
		interactions: make(map[string]int),
	}

	var prevCharID string
	var prevEpisode int

	for _, q := range s.quotes {
		if episode != AllEpisodes && q.Episode != episode {
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
			a, b := prevCharID, q.CharacterID
			if a > b {
				a, b = b, a
			}
			t.interactions[fmt.Sprintf("%s|%s", a, b)]++
		}

		prevCharID = q.CharacterID
		prevEpisode = q.Episode
	}

	return t
}

func (*statsComputer) rankCharacters(charCounts map[string]int) []rankedChar {
	ranked := make([]rankedChar, 0, len(charCounts))
	for id, count := range charCounts {
		ranked = append(ranked, rankedChar{id, count})
	}
	slices.SortFunc(ranked, func(a, b rankedChar) int {
		return cmp.Compare(b.count, a.count)
	})
	return ranked
}

func (*statsComputer) topSpeakers(ranked []rankedChar, n int) []speakerStat {
	if len(ranked) < n {
		n = len(ranked)
	}
	result := make([]speakerStat, n)
	for i := 0; i < n; i++ {
		result[i] = speakerStat{
			CharacterID: ranked[i].id,
			Name:        CharacterNames.GetCharacterName(ranked[i].id),
			Count:       ranked[i].count,
		}
	}
	return result
}

func (*statsComputer) linesPerEpisode(charEpCounts map[string]map[int]int, ranked []rankedChar, topN int) []episodeCharacterLines {
	if len(ranked) < topN {
		topN = len(ranked)
	}
	topSet := make(map[string]bool, topN)
	for i := 0; i < topN; i++ {
		topSet[ranked[i].id] = true
	}

	result := make([]episodeCharacterLines, 8)
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
		result[ep-1] = episodeCharacterLines{
			Episode:     ep,
			EpisodeName: episodeNames[ep],
			Characters:  chars,
		}
	}
	return result
}

func (*statsComputer) truthPerEpisode(epTruth map[int][2]int) []episodeTruth {
	result := make([]episodeTruth, 8)
	for ep := 1; ep <= 8; ep++ {
		counts := epTruth[ep]
		result[ep-1] = episodeTruth{
			Episode: ep,
			Red:     counts[0],
			Blue:    counts[1],
		}
	}
	return result
}

func (*statsComputer) topInteractions(interactionCounts map[string]int, n int) []interactionPair {
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

	result := make([]interactionPair, n)
	for i := 0; i < n; i++ {
		parts := strings.SplitN(sorted[i].key, "|", 2)
		result[i] = interactionPair{
			CharA: parts[0],
			CharB: parts[1],
			NameA: CharacterNames.GetCharacterName(parts[0]),
			NameB: CharacterNames.GetCharacterName(parts[1]),
			Count: sorted[i].count,
		}
	}
	return result
}

func (*statsComputer) buildCharacterPresence(ranked []rankedChar, charEpCounts map[string]map[int]int, n int) []characterPresence {
	if len(ranked) < n {
		n = len(ranked)
	}
	result := make([]characterPresence, n)
	for i := 0; i < n; i++ {
		id := ranked[i].id
		episodes := make([]int, 8)
		for ep := 1; ep <= 8; ep++ {
			episodes[ep-1] = charEpCounts[id][ep]
		}
		result[i] = characterPresence{
			CharacterID: id,
			Name:        CharacterNames.GetCharacterName(id),
			Episodes:    episodes,
		}
	}
	return result
}

func (*statsComputer) buildNameMap(charCounts map[string]int) map[string]string {
	nameMap := make(map[string]string, len(charCounts))
	for id := range charCounts {
		nameMap[id] = CharacterNames.GetCharacterName(id)
	}
	return nameMap
}
