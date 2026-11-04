package ciconia

import (
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/quotes"

	cicharacter "github.com/VictoriqueMoe/umineko_script_parser/ciconia/character"
)

type ciconiaStatsComputer struct {
	cached *dto.CiconiaStatsResult
}

func NewCiconiaStats(t params.SpeakerTallies) quotes.Stats {
	charCounts := t.CharacterTotals()
	interactions := t.InteractionTotals()
	ranked := quotes.RankCharacters(charCounts)

	return &ciconiaStatsComputer{cached: &dto.CiconiaStatsResult{
		TopSpeakers:       quotes.TopSpeakers(ranked, 20, resolveName),
		Interactions:      quotes.TopInteractions(interactions, 25, resolveName),
		InteractionCounts: interactions,
		CharacterNames:    quotes.BuildNameMap(charCounts, resolveName),
		LinesPerChapter:   quotes.LinesPerSegment(t.CharacterSegments, ranked, 10),
	}}
}

func (s *ciconiaStatsComputer) Compute(_ int) any {
	return s.cached
}

func resolveName(id string) string {
	return cicharacter.CharacterNames.GetCharacterName(cicharacter.Character(id))
}
