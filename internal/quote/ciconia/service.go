package ciconia

import (
	"fmt"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/quotes"
	"umineko_quote/internal/repository"

	cicharacter "github.com/VictoriqueMoe/umineko_script_parser/ciconia/character"
)

type (
	Service interface {
		quotes.Service[dto.CiconiaQuote]
		GetCharacters() dto.CharactersResult
		GetStats() quotes.Stats
	}

	ciconiaService struct {
		quotes.Service[dto.CiconiaQuote]
		stats      quotes.Stats
		characters dto.CharactersResult
	}
)

func NewService(repo repository.QuoteRepository[dto.CiconiaQuote]) (Service, error) {
	tallies, err := repo.SpeakerTallies(language.English)
	if err != nil {
		return nil, fmt.Errorf("ciconia speaker tallies: %w", err)
	}

	speakers, err := repo.Speakers(language.English)
	if err != nil {
		return nil, fmt.Errorf("ciconia speakers: %w", err)
	}

	return &ciconiaService{
		Service:    quotes.NewService(repo, resolveName),
		stats:      NewCiconiaStats(tallies),
		characters: quotes.NewCharactersResult(quotes.CuratedCharacters(cicharacter.MainCharacterNames.GetAllCharacters()), speakers),
	}, nil
}

func (s *ciconiaService) GetCharacters() dto.CharactersResult {
	return s.characters
}

func (s *ciconiaService) GetStats() quotes.Stats {
	return s.stats
}
