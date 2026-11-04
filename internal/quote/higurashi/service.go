package higurashi

import (
	"fmt"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/quotes"
	"umineko_quote/internal/repository"

	hicharacter "github.com/VictoriqueMoe/umineko_script_parser/higurashi/character"
)

type (
	Service interface {
		quotes.Service[dto.HigurashiQuote]
		GetCharacters() dto.CharactersResult
		GetStats() quotes.Stats
	}

	higurashiService struct {
		quotes.Service[dto.HigurashiQuote]
		stats      quotes.Stats
		characters dto.CharactersResult
	}
)

func NewService(repo repository.QuoteRepository[dto.HigurashiQuote]) (Service, error) {
	tallies, err := repo.SpeakerTallies(language.English)
	if err != nil {
		return nil, fmt.Errorf("higurashi speaker tallies: %w", err)
	}

	speakers, err := repo.Speakers(language.English)
	if err != nil {
		return nil, fmt.Errorf("higurashi speakers: %w", err)
	}

	return &higurashiService{
		Service:    quotes.NewService(repo, resolveName),
		stats:      NewHigurashiStats(tallies),
		characters: quotes.NewCharactersResult(quotes.CuratedCharacters(hicharacter.CharacterNames.GetAllCharacters()), speakers),
	}, nil
}

func (s *higurashiService) GetCharacters() dto.CharactersResult {
	return s.characters
}

func (s *higurashiService) GetStats() quotes.Stats {
	return s.stats
}
