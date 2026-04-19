package umineko

import (
	"log"
	"sync"
	"time"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/store"
	"umineko_quote/internal/quote/subtitle"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/umineko"
	"github.com/VictoriqueMoe/umineko_script_parser/umineko/character"
)

const uminekoAudioDir = "internal/quote/data/audio"

type (
	Service interface {
		Search(params params.SearchParams, truth string) dto.SearchResponse
		Browse(params params.BrowseParams, truth string) dto.CharacterResponse
		GetByAudioID(lang language.Language, audioID string) *dto.UminekoQuote
		GetByIndex(lang language.Language, index int) *dto.UminekoQuote
		GetContext(lang language.Language, audioID string, lines int) *dto.ContextResponse
		NearestVoicedAudioID(lang language.Language, audioID string, direction string) string
		Random(lang language.Language, characterID string, episode int, truth Truth) *dto.UminekoQuote
		GetCharacters() map[character.Character]string
		AudioFilePath(characterId string, audioId string) string
		GetStats() store.Stats
		HasAudio() bool
		LoadedLanguages() map[language.Language]int
	}

	uminekoService struct {
		store *store.GameStore[dto.UminekoQuote]
		stats store.Stats
	}

	uminekoLangParseResult struct {
		lang   language.Language
		parsed []dto.UminekoQuote
	}
)

func NewService() Service {
	serviceStart := time.Now()

	langFiles := map[language.Language]string{
		language.English:    "data/en.file",
		language.WitchHunt:  "data/wh.file",
		language.Japanese:   "data/ja.file",
		language.Russian:    "data/ru.file",
		language.Spanish:    "data/es.file",
		language.Portuguese: "data/pt.file",
	}

	results := make(chan uminekoLangParseResult, len(langFiles))
	var wg sync.WaitGroup

	for lang, path := range langFiles {
		wg.Go(func() {
			f, err := quote.DataFS.Open(path)
			if err != nil {
				log.Printf("[umineko/%s] failed to open %s: %v", lang, path, err)
				return
			}
			defer f.Close()

			timeStart := time.Now()
			parser := umineko.NewParser()
			parsed, validationErrors, err := scriptparser.ParseReader(f, parser)
			subtitleRefs := parser.SubtitleRefs()
			timeEnd := time.Now()

			log.Printf("[umineko/%s] parsed %d quotes took %v", lang, len(parsed), timeEnd.Sub(timeStart))

			if err != nil {
				log.Printf("[umineko/%s] failed to parse %s: %v", lang, path, err)
				return
			}
			log.Printf("[umineko/%s] parsed %s: %d quotes", lang, path, len(parsed))

			if len(validationErrors) > 0 {
				errorCount := 0
				warningCount := 0
				for _, ve := range validationErrors {
					if ve.Severity == scriptparser.SeverityError {
						errorCount++
					} else {
						warningCount++
					}
				}
				log.Printf("[umineko/%s] validation: %d errors, %d warnings", lang, errorCount, warningCount)
				limit := len(validationErrors)
				if limit > 10 {
					limit = 10
				}
				for i := 0; i < limit; i++ {
					log.Printf("[umineko/%s]   %s", lang, validationErrors[i])
				}
				if len(validationErrors) > 10 {
					log.Printf("[umineko/%s]   ... and %d more", lang, len(validationErrors)-10)
				}
			}
			indexed := make([]dto.UminekoQuote, len(parsed))
			for i, q := range parsed {
				indexed[i] = dto.UminekoQuote{
					UminekoQuote: q,
					Index:        i,
				}
			}

			subQuotes := subtitle.ResolveRefs(quote.DataFS, subtitleRefs)
			if len(subQuotes) > 0 {
				startIdx := len(indexed)
				for i := 0; i < len(subQuotes); i++ {
					indexed = append(indexed, dto.UminekoQuote{
						UminekoQuote: subQuotes[i].UminekoQuote,
						Index:        startIdx + i,
					})
				}
				log.Printf("[umineko/%s] added %d subtitle quotes", lang, len(subQuotes))
			}

			results <- uminekoLangParseResult{
				lang:   lang,
				parsed: indexed,
			}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	quotes := make(map[language.Language][]dto.UminekoQuote)

	for r := range results {
		quotes[r.lang] = r.parsed
	}

	gameStore := store.NewGameStore(quotes, uminekoAudioDir, func(q *dto.UminekoQuote) *scriptdto.BaseQuote {
		return &q.BaseQuote
	})

	if gameStore.HasAudio() {
		log.Printf("[umineko/audio] audio features enabled")
	} else {
		log.Printf("[umineko/audio] no audio files found, disabling audio features")
	}

	log.Printf("[umineko] initialised in %v", time.Since(serviceStart).Round(time.Millisecond))

	return &uminekoService{
		store: gameStore,
		stats: NewUminekoStats(quotes[language.English]),
	}
}

func (s *uminekoService) Search(p params.SearchParams, truth string) dto.SearchResponse {
	truthVal := TruthAll.Parse(truth)
	var filter func(*dto.UminekoQuote) bool
	if truthVal != TruthAll {
		filter = func(q *dto.UminekoQuote) bool {
			return matchesTruth(q, truthVal)
		}
	}
	return s.store.Search(p, filter)
}

func (s *uminekoService) Browse(p params.BrowseParams, truth string) dto.CharacterResponse {
	truthVal := TruthAll.Parse(truth)
	var filter func(*dto.UminekoQuote) bool
	if truthVal != TruthAll {
		filter = func(q *dto.UminekoQuote) bool {
			return matchesTruth(q, truthVal)
		}
	}
	return s.store.Browse(p, filter, func(id string) string {
		return character.CharacterNames.GetCharacterName(character.CharacterFromID(id))
	})
}

func (s *uminekoService) Random(lang language.Language, characterID string, episode int, truth Truth) *dto.UminekoQuote {
	var filter func(*dto.UminekoQuote) bool
	if truth != TruthAll {
		filter = func(q *dto.UminekoQuote) bool {
			return matchesTruth(q, truth)
		}
	}
	return s.store.Random(lang, characterID, episode, filter)
}

func (s *uminekoService) GetByAudioID(lang language.Language, audioID string) *dto.UminekoQuote {
	return s.store.GetByAudioID(lang, audioID)
}

func (s *uminekoService) GetByIndex(lang language.Language, index int) *dto.UminekoQuote {
	return s.store.GetByIndex(lang, index)
}

func (s *uminekoService) GetContext(lang language.Language, audioID string, lines int) *dto.ContextResponse {
	return s.store.GetContext(lang, audioID, lines)
}

func (s *uminekoService) NearestVoicedAudioID(lang language.Language, audioID string, direction string) string {
	return s.store.NearestVoicedAudioID(lang, audioID, direction)
}

func (s *uminekoService) GetCharacters() map[character.Character]string {
	return character.CharacterNames.GetAllCharacters()
}

func (s *uminekoService) AudioFilePath(characterId string, audioId string) string {
	return s.store.AudioFilePath(characterId, audioId)
}

func (s *uminekoService) GetStats() store.Stats {
	return s.stats
}

func (s *uminekoService) HasAudio() bool {
	return s.store.HasAudio()
}

func (s *uminekoService) LoadedLanguages() map[language.Language]int {
	return s.store.LoadedLanguages()
}

func matchesTruth(q *dto.UminekoQuote, truth Truth) bool {
	if truth == TruthRed && !q.HasRedTruth {
		return false
	}
	if truth == TruthBlue && !q.HasBlueTruth {
		return false
	}
	if truth == TruthGold && !q.HasGoldTruth {
		return false
	}
	if truth == TruthPurple && !q.HasPurpleTruth {
		return false
	}
	return true
}
