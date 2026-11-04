package higurashi

import (
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/quotes"

	hicharacter "github.com/VictoriqueMoe/umineko_script_parser/higurashi/character"
)

type higurashiStatsComputer struct {
	cached *dto.HigurashiStatsResult
}

func NewHigurashiStats(t params.SpeakerTallies) quotes.Stats {
	charCounts := t.CharacterTotals()
	interactions := t.InteractionTotals()
	ranked := quotes.RankCharacters(charCounts)

	return &higurashiStatsComputer{cached: &dto.HigurashiStatsResult{
		TopSpeakers:       quotes.TopSpeakers(ranked, 20, resolveName),
		Interactions:      quotes.TopInteractions(interactions, 25, resolveName),
		InteractionCounts: interactions,
		CharacterNames:    quotes.BuildNameMap(charCounts, resolveName),
		LinesPerArc:       quotes.LinesPerSegment(t.CharacterSegments, ranked, 10),
	}}
}

func (s *higurashiStatsComputer) Compute(_ int) any {
	return s.cached
}

func resolveName(id string) string {
	return hicharacter.CharacterNames.GetCharacterName(hicharacter.CharacterFromID(id))
}
