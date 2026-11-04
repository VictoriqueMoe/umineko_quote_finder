package umineko

import (
	"fmt"
	"log"

	"umineko_quote/internal/audio"
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/quotes"
	"umineko_quote/internal/repository"

	"github.com/VictoriqueMoe/umineko_script_parser/umineko/character"
)

type (
	Service interface {
		quotes.Service[dto.UminekoQuote]
		GetCharacters() map[character.Character]string
		AudioFilePath(characterId string, audioId string) string
		GetStats() quotes.Stats
		HasAudio() bool
	}

	uminekoService struct {
		quotes.Service[dto.UminekoQuote]
		audio.Resolver
		stats quotes.Stats
	}
)

func NewService(repo repository.UminekoRepository, resolver audio.Resolver) (Service, error) {
	tallies, err := repo.SpeakerTallies(language.English)
	if err != nil {
		return nil, fmt.Errorf("umineko speaker tallies: %w", err)
	}

	truth, err := repo.TruthPerEpisode(language.English)
	if err != nil {
		return nil, fmt.Errorf("umineko truth per episode: %w", err)
	}

	stats, err := NewUminekoStats(tallies, truth)
	if err != nil {
		return nil, err
	}

	if resolver.HasAudio() {
		log.Printf("[umineko/audio] audio features enabled")
	} else {
		log.Printf("[umineko/audio] no audio files found, disabling audio features")
	}

	return &uminekoService{
		Service:  quotes.NewService(repo, resolveName),
		Resolver: resolver,
		stats:    stats,
	}, nil
}

func (s *uminekoService) GetCharacters() map[character.Character]string {
	return character.CharacterNames.GetAllCharacters()
}

func (s *uminekoService) GetStats() quotes.Stats {
	return s.stats
}
