package ciconia

import (
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/store"

	cicharacter "github.com/VictoriqueMoe/umineko_script_parser/ciconia/character"
)

type (
	ciconiaStatsComputer struct {
		quotes []dto.CiconiaQuote
		cached *dto.CiconiaStatsResult
	}

	ciconiaTallies struct {
		charCounts        map[string]int
		charChapterCounts map[string]map[string]int
		interactions      map[string]int
	}
)

func NewCiconiaStats(quotes []dto.CiconiaQuote) store.Stats {
	s := &ciconiaStatsComputer{quotes: quotes}
	s.cached = s.compute("")
	return s
}

func (s *ciconiaStatsComputer) Compute(_ int) any {
	return s.cached
}

func resolveName(id string) string {
	return cicharacter.CharacterNames.GetCharacterName(cicharacter.Character(id))
}

func (s *ciconiaStatsComputer) compute(chapter string) *dto.CiconiaStatsResult {
	t := s.tally(chapter)
	ranked := store.RankCharacters(t.charCounts)

	return &dto.CiconiaStatsResult{
		TopSpeakers:       store.TopSpeakers(ranked, 20, resolveName),
		Interactions:      store.TopInteractions(t.interactions, 25, resolveName),
		InteractionCounts: t.interactions,
		CharacterNames:    store.BuildNameMap(t.charCounts, resolveName),
		LinesPerChapter:   ciconiaLinesPerChapter(t.charChapterCounts, ranked, 10),
	}
}

func (s *ciconiaStatsComputer) tally(chapter string) ciconiaTallies {
	t := ciconiaTallies{
		charCounts:        make(map[string]int),
		charChapterCounts: make(map[string]map[string]int),
		interactions:      make(map[string]int),
	}

	var prevCharID string
	var prevChapter string

	for _, q := range s.quotes {
		if chapter != "" && q.Chapter != chapter {
			prevCharID = ""
			continue
		}

		if q.CharacterID == "narrator" {
			prevCharID = ""
			continue
		}

		t.charCounts[q.CharacterID]++

		if t.charChapterCounts[q.CharacterID] == nil {
			t.charChapterCounts[q.CharacterID] = make(map[string]int)
		}
		t.charChapterCounts[q.CharacterID][q.Chapter]++

		if prevCharID != "" && prevCharID != q.CharacterID && prevChapter == q.Chapter {
			if key, ok := store.InteractionPairKey(prevCharID, q.CharacterID); ok {
				t.interactions[key]++
			}
		}

		prevCharID = q.CharacterID
		prevChapter = q.Chapter
	}

	return t
}

func ciconiaLinesPerChapter(charChapterCounts map[string]map[string]int, ranked []store.RankedChar, topN int) map[string]map[string]int {
	if len(ranked) < topN {
		topN = len(ranked)
	}
	topSet := make(map[string]bool, topN)
	for i := 0; i < topN; i++ {
		topSet[ranked[i].ID] = true
	}

	chapters := make(map[string]map[string]int)
	for id, chMap := range charChapterCounts {
		label := id
		if !topSet[id] {
			label = "other"
		}
		for ch, count := range chMap {
			if chapters[ch] == nil {
				chapters[ch] = make(map[string]int)
			}
			chapters[ch][label] += count
		}
	}

	return chapters
}
