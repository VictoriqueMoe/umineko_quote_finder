package quote

import (
	"embed"
	"log"
	"math/rand/v2"
	"strings"
	"sync"
	"time"
	"umineko_quote/internal/quote/character"
	"umineko_quote/internal/quote/language"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/lexar"
	"umineko_quote/internal/quote/scriptloader"
)

const audioDir = "internal/quote/data/audio"

//go:embed data/*.file data/sub/*.ass
var dataFS embed.FS

type (
	Service interface {
		Search(query string, lang language.Language, limit int, offset int, character character.Character, episode int, truth Truth) dto.SearchResponse
		Browse(lang language.Language, character character.Character, limit int, offset int, episode int, truth Truth) dto.CharacterResponse
		GetByAudioID(lang language.Language, audioID string) *dto.ParsedQuote
		GetContext(lang language.Language, audioID string, lines int) *dto.ContextResponse
		Random(lang language.Language, character character.Character, episode int, truth Truth) *dto.ParsedQuote
		GetCharacters() map[character.Character]string
		AudioFilePath(characterId string, audioId string) string
		GetStats() Stats
		HasAudio() bool
		LoadedLanguages() map[language.Language]int
	}

	service struct {
		quotes  map[language.Language][]dto.ParsedQuote
		indexer Indexer
		stats   Stats
	}

	langParseResult struct {
		lang   language.Language
		parsed []dto.ParsedQuote
	}
)

func NewService() Service {
	serviceStart := time.Now()

	parse := func(lines []string) ([]dto.ParsedQuote, []lexar.SubtitleRef) {
		p := NewParser()
		return p.ParseAll(lines), p.SubtitleRefs()
	}

	loader := scriptloader.New(dataFS, parse)

	langFiles := map[language.Language]string{
		language.English:    "data/en.file",
		language.Japanese:   "data/ja.file",
		language.Spanish:    "data/es.file",
		language.Portuguese: "data/pt.file",
	}

	results := make(chan langParseResult, len(langFiles))
	var wg sync.WaitGroup

	for lang, path := range langFiles {
		wg.Go(func() {
			parsed := loader.Load(string(lang), path)
			if parsed == nil {
				return
			}
			results <- langParseResult{
				lang:   lang,
				parsed: parsed,
			}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	quotes := make(map[language.Language][]dto.ParsedQuote)

	for r := range results {
		quotes[r.lang] = r.parsed
	}

	indexer := NewIndexer(quotes, audioDir)

	if indexer.HasAudio() {
		log.Printf("[audio] audio features enabled")
	} else {
		log.Printf("[audio] no audio files found, disabling audio features")
	}

	log.Printf("[system] initialised in %v", time.Since(serviceStart).Round(time.Millisecond))

	return &service{
		quotes:  quotes,
		indexer: indexer,
		stats:   NewStats(quotes[language.English]),
	}
}

func (s *service) Search(query string, lang language.Language, limit int, offset int, character character.Character, episode int, truth Truth) dto.SearchResponse {
	characterID := character.ID()
	if limit <= 0 {
		limit = 30
	}
	if offset < 0 {
		offset = 0
	}

	quotes := s.quotes[lang]
	lowerTexts := s.indexer.LowerTexts(lang)
	if quotes == nil {
		return NewSearchResponse(nil, limit, offset)
	}

	matchesFilter := func(q dto.ParsedQuote) bool {
		if characterID != "" && q.CharacterID != characterID {
			return false
		}
		if episode > 0 && q.Episode != episode {
			return false
		}
		if truth == TruthRed && !q.HasRedTruth {
			return false
		}
		if truth == TruthBlue && !q.HasBlueTruth {
			return false
		}
		return true
	}

	queryLower := strings.ToLower(query)

	searchIndices := s.indexer.FilteredIndices(lang, characterID, episode)

	var exactMatches []dto.SearchResult
	if searchIndices != nil {
		if len(searchIndices) > 5000 {
			exactMatches = concurrentExactSearch(searchIndices, lowerTexts, quotes, queryLower, matchesFilter)
		} else {
			for _, idx := range searchIndices {
				if strings.Contains(lowerTexts[idx], queryLower) {
					if matchesFilter(quotes[idx]) {
						exactMatches = append(exactMatches, NewSearchResult(quotes[idx], 100))
					}
				}
			}
		}
	} else {
		allIndices := make([]int, len(quotes))
		for i := range allIndices {
			allIndices[i] = i
		}
		exactMatches = concurrentExactSearch(allIndices, lowerTexts, quotes, queryLower, matchesFilter)
	}

	return NewSearchResponse(exactMatches, limit, offset)
}

func (s *service) Browse(lang language.Language, character character.Character, limit int, offset int, episode int, truth Truth) dto.CharacterResponse {
	characterID := character.ID()
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	quotes := s.quotes[lang]
	if quotes == nil {
		return NewCharacterResponse(characterID, nil, limit, offset)
	}

	var source []int
	indexed := s.indexer.FilteredIndices(lang, characterID, episode)
	if indexed != nil {
		source = indexed
	} else {
		source = make([]int, len(quotes))
		for i := range source {
			source[i] = i
		}
	}

	var all []dto.ParsedQuote
	for _, idx := range source {
		q := quotes[idx]
		if truth == TruthRed && !q.HasRedTruth {
			continue
		}
		if truth == TruthBlue && !q.HasBlueTruth {
			continue
		}
		all = append(all, q)
	}

	return NewCharacterResponse(characterID, all, limit, offset)
}

func (s *service) Random(lang language.Language, character character.Character, episode int, truth Truth) *dto.ParsedQuote {
	characterID := character.ID()

	quotes := s.quotes[lang]
	if quotes == nil || len(quotes) == 0 {
		return nil
	}

	matchesTruth := func(q dto.ParsedQuote) bool {
		if truth == TruthRed && !q.HasRedTruth {
			return false
		}
		if truth == TruthBlue && !q.HasBlueTruth {
			return false
		}
		return true
	}

	if characterID == "" && episode <= 0 && truth == TruthAll {
		indices := s.indexer.NonNarratorIndices(lang)
		if len(indices) == 0 {
			return nil
		}
		pick := indices[rand.IntN(len(indices))]
		return &quotes[pick]
	}

	var candidates []int

	if truth != TruthAll {
		var source []int
		indexed := s.indexer.FilteredIndices(lang, characterID, episode)
		if indexed != nil {
			source = indexed
		} else if characterID == "" && episode <= 0 {
			source = s.indexer.NonNarratorIndices(lang)
		}

		if source != nil {
			for _, idx := range source {
				if matchesTruth(quotes[idx]) {
					candidates = append(candidates, idx)
				}
			}
		} else {
			for i := 0; i < len(quotes); i++ {
				if characterID != "" && quotes[i].CharacterID != characterID {
					continue
				}
				if episode > 0 && quotes[i].Episode != episode {
					continue
				}
				if matchesTruth(quotes[i]) {
					candidates = append(candidates, i)
				}
			}
		}

		if len(candidates) == 0 {
			return nil
		}
		pick := candidates[rand.IntN(len(candidates))]
		return &quotes[pick]
	}

	indices := s.indexer.FilteredIndices(lang, characterID, episode)
	if indices != nil {
		if len(indices) == 0 {
			return nil
		}
		pick := indices[rand.IntN(len(indices))]
		return &quotes[pick]
	}

	var filtered []int
	for i := 0; i < len(quotes); i++ {
		if characterID != "" && quotes[i].CharacterID != characterID {
			continue
		}
		if episode > 0 && quotes[i].Episode != episode {
			continue
		}
		filtered = append(filtered, i)
	}

	if len(filtered) == 0 {
		return nil
	}

	pick := filtered[rand.IntN(len(filtered))]
	return &quotes[pick]
}

func (s *service) GetByAudioID(lang language.Language, audioID string) *dto.ParsedQuote {
	quotes := s.quotes[lang]
	if quotes == nil {
		return nil
	}

	idx, ok := s.indexer.QuoteIndex(lang, audioID)
	if !ok {
		return nil
	}
	return &quotes[idx]
}

func (s *service) GetContext(lang language.Language, audioID string, lines int) *dto.ContextResponse {
	if lines <= 0 {
		lines = 5
	}
	if lines > 20 {
		lines = 20
	}

	quotes := s.quotes[lang]
	if quotes == nil {
		return nil
	}

	idx, ok := s.indexer.QuoteIndex(lang, audioID)
	if !ok {
		return nil
	}

	start := idx - lines
	if start < 0 {
		start = 0
	}
	end := idx + lines + 1
	if end > len(quotes) {
		end = len(quotes)
	}

	return &dto.ContextResponse{
		Before: quotes[start:idx],
		Quote:  quotes[idx],
		After:  quotes[idx+1 : end],
	}
}

func (s *service) GetCharacters() map[character.Character]string {
	return character.CharacterNames.GetAllCharacters()
}

func (s *service) AudioFilePath(characterId string, audioId string) string {
	return s.indexer.AudioFilePath(characterId, audioId)
}

func (s *service) GetStats() Stats {
	return s.stats
}

func (s *service) HasAudio() bool {
	return s.indexer.HasAudio()
}

func (s *service) LoadedLanguages() map[language.Language]int {
	result := make(map[language.Language]int, len(s.quotes))
	for lang, quotes := range s.quotes {
		result[lang] = len(quotes)
	}
	return result
}
