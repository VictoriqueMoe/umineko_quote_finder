package higurashi

import (
	"umineko_quote/internal/quote/store"

	hicharacter "github.com/VictoriqueMoe/umineko_script_parser/higurashi/character"

	"umineko_quote/internal/dto"
)

type (
	higurashiStatsComputer struct {
		quotes []dto.HigurashiQuote
		cached *dto.HigurashiStatsResult
	}

	higurashiTallies struct {
		charCounts    map[string]int
		charArcCounts map[string]map[string]int
		interactions  map[string]int
	}
)

func NewHigurashiStats(quotes []dto.HigurashiQuote) store.Stats {
	s := &higurashiStatsComputer{quotes: quotes}
	s.cached = s.compute("")
	return s
}

func (s *higurashiStatsComputer) Compute(episode int) any {
	return s.cached
}

func resolveName(id string) string {
	return hicharacter.CharacterNames.GetCharacterName(hicharacter.CharacterFromID(id))
}

func (s *higurashiStatsComputer) compute(arc string) *dto.HigurashiStatsResult {
	t := s.tally(arc)
	ranked := store.RankCharacters(t.charCounts)

	return &dto.HigurashiStatsResult{
		TopSpeakers:       store.TopSpeakers(ranked, 20, resolveName),
		Interactions:      store.TopInteractions(t.interactions, 25, resolveName),
		InteractionCounts: t.interactions,
		CharacterNames:    store.BuildNameMap(t.charCounts, resolveName),
		LinesPerArc:       higurashiLinesPerArc(t.charArcCounts, ranked, 10),
	}
}

func (s *higurashiStatsComputer) tally(arc string) higurashiTallies {
	t := higurashiTallies{
		charCounts:    make(map[string]int),
		charArcCounts: make(map[string]map[string]int),
		interactions:  make(map[string]int),
	}

	var prevCharID string
	var prevArc string

	for _, q := range s.quotes {
		if arc != "" && q.Arc != arc {
			prevCharID = ""
			continue
		}

		if q.CharacterID == "narrator" {
			prevCharID = ""
			continue
		}

		t.charCounts[q.CharacterID]++

		if t.charArcCounts[q.CharacterID] == nil {
			t.charArcCounts[q.CharacterID] = make(map[string]int)
		}
		t.charArcCounts[q.CharacterID][q.Arc]++

		if prevCharID != "" && prevCharID != q.CharacterID && prevArc == q.Arc {
			if key, ok := store.InteractionPairKey(prevCharID, q.CharacterID); ok {
				t.interactions[key]++
			}
		}

		prevCharID = q.CharacterID
		prevArc = q.Arc
	}

	return t
}

func higurashiLinesPerArc(charArcCounts map[string]map[string]int, ranked []store.RankedChar, topN int) map[string]map[string]int {
	if len(ranked) < topN {
		topN = len(ranked)
	}
	topSet := make(map[string]bool, topN)
	for i := 0; i < topN; i++ {
		topSet[ranked[i].ID] = true
	}

	arcs := make(map[string]map[string]int)
	for id, arcMap := range charArcCounts {
		label := id
		if !topSet[id] {
			label = "other"
		}
		for arc, count := range arcMap {
			if arcs[arc] == nil {
				arcs[arc] = make(map[string]int)
			}
			arcs[arc][label] += count
		}
	}

	return arcs
}
