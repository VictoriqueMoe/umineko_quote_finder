package quote

import (
	"embed"
	"log"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"
	quoteparams "umineko_quote/internal/quote/params"

	"umineko_quote/internal/quote/subtitle"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
	"github.com/VictoriqueMoe/umineko_script_parser/lexer"
	"github.com/VictoriqueMoe/umineko_script_parser/quote/character"
)

const audioDir = "internal/quote/data/audio"

//go:embed data/*.file data/sub/*.ass
var dataFS embed.FS

type (
	Service interface {
		Search(params quoteparams.SearchParams) dto.SearchResponse
		Browse(params quoteparams.BrowseParams) dto.CharacterResponse
		GetByAudioID(lang language.Language, audioID string) *dto.ParsedQuote
		GetContext(lang language.Language, audioID string, lines int) *dto.ContextResponse
		NearestVoicedAudioID(lang language.Language, audioID string, direction string) string
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

	langFiles := map[language.Language]string{
		language.English:    "data/en.file",
		language.WitchHunt:  "data/wh.file",
		language.Japanese:   "data/ja.file",
		language.Russian:    "data/ru.file",
		language.Spanish:    "data/es.file",
		language.Portuguese: "data/pt.file",
	}

	results := make(chan langParseResult, len(langFiles))
	var wg sync.WaitGroup

	for lang, path := range langFiles {
		wg.Go(func() {
			f, err := dataFS.Open(path)
			if err != nil {
				log.Printf("[%s] failed to open %s: %v", lang, path, err)
				return
			}
			defer f.Close()

			timeStart := time.Now()
			parsed, subtitleRefs, validationErrors, err := scriptparser.ParseFile(f)
			timeEnd := time.Now()

			log.Printf("[%s] parsed %d quotes took %v", lang, len(parsed), timeEnd.Sub(timeStart))

			if err != nil {
				log.Printf("[%s] failed to parse %s: %v", lang, path, err)
				return
			}
			log.Printf("[%s] parsed %s: %d quotes", lang, path, len(parsed))

			if len(validationErrors) > 0 {
				errorCount := 0
				warningCount := 0
				for _, ve := range validationErrors {
					if ve.Severity == lexer.SeverityError {
						errorCount++
					} else {
						warningCount++
					}
				}
				log.Printf("[%s] validation: %d errors, %d warnings", lang, errorCount, warningCount)
				limit := len(validationErrors)
				if limit > 10 {
					limit = 10
				}
				for i := 0; i < limit; i++ {
					log.Printf("[%s]   %s", lang, validationErrors[i])
				}
				if len(validationErrors) > 10 {
					log.Printf("[%s]   ... and %d more", lang, len(validationErrors)-10)
				}
			}
			subQuotes := subtitle.ResolveRefs(dataFS, subtitleRefs)
			if len(subQuotes) > 0 {
				parsed = append(parsed, subQuotes...)
				log.Printf("[%s] added %d subtitle quotes", lang, len(subQuotes))
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

func (s *service) Search(params quoteparams.SearchParams) dto.SearchResponse {
	if params.Lang == language.Auto {
		return s.searchAuto(params)
	}

	characterID := params.Character.ID()
	query := params.Query
	lang := params.Lang
	limit := params.Limit
	offset := params.Offset
	episode := params.Episode
	truth := TruthAll.Parse(params.Truth)
	interactionA := params.InteractionA
	interactionB := params.InteractionB

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
		if truth == TruthGold && !q.HasGoldTruth {
			return false
		}
		if truth == TruthPurple && !q.HasPurpleTruth {
			return false
		}
		return true
	}

	queryLower := strings.ToLower(query)

	searchIndices := s.indexer.FilteredIndices(lang, characterID, episode)
	if interactionA != "" && interactionB != "" {
		interactionIndices := s.indexer.InteractionIndices(lang, interactionA, interactionB)
		searchIndices = mergeFilteredIndices(searchIndices, interactionIndices)
	}

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

func (s *service) searchAuto(params quoteparams.SearchParams) dto.SearchResponse {
	type langResult struct {
		lang   language.Language
		result dto.SearchResponse
	}

	results := make([]langResult, len(language.All))
	var wg sync.WaitGroup
	for i, lang := range language.All {
		wg.Go(func() {
			p := params
			p.Lang = lang
			results[i] = langResult{lang: lang, result: s.Search(p)}
		})
	}
	wg.Wait()

	for _, lr := range results {
		if lr.result.Total > 0 {
			lr.result.Lang = string(lr.lang)
			return lr.result
		}
	}

	params.Lang = language.English
	return s.Search(params)
}

func mergeFilteredIndices(baseIndices []int, interactionIndices []int) []int {
	if baseIndices == nil {
		return interactionIndices
	}
	if len(baseIndices) == 0 || len(interactionIndices) == 0 {
		return []int{}
	}

	interactionSet := make(map[int]struct{}, len(interactionIndices))
	for _, idx := range interactionIndices {
		interactionSet[idx] = struct{}{}
	}

	out := make([]int, 0, len(baseIndices))
	for _, idx := range baseIndices {
		if _, ok := interactionSet[idx]; ok {
			out = append(out, idx)
		}
	}
	return out
}

func (s *service) Browse(params quoteparams.BrowseParams) dto.CharacterResponse {
	lang := params.Lang
	limit := params.Limit
	offset := params.Offset
	episode := params.Episode
	truth := TruthAll.Parse(params.Truth)
	interactionA := params.InteractionA
	interactionB := params.InteractionB
	characterID := params.Character.ID()

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
	if interactionA != "" && interactionB != "" {
		interactionIndices := s.indexer.InteractionIndices(lang, interactionA, interactionB)
		indexed = mergeFilteredIndices(indexed, interactionIndices)
	}
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
		if truth == TruthGold && !q.HasGoldTruth {
			continue
		}
		if truth == TruthPurple && !q.HasPurpleTruth {
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
		if truth == TruthGold && !q.HasGoldTruth {
			return false
		}
		if truth == TruthPurple && !q.HasPurpleTruth {
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

func (s *service) NearestVoicedAudioID(lang language.Language, audioID string, direction string) string {
	quotes := s.quotes[lang]
	if quotes == nil {
		return ""
	}

	idx, ok := s.indexer.QuoteIndex(lang, audioID)
	if !ok {
		return ""
	}

	if direction == "prev" {
		for i := idx - 1; i >= 0; i-- {
			if quotes[i].AudioID != "" {
				return strings.Split(quotes[i].AudioID, ", ")[0]
			}
		}
		return ""
	}

	for i := idx + 1; i < len(quotes); i++ {
		if quotes[i].AudioID != "" {
			return strings.Split(quotes[i].AudioID, ", ")[0]
		}
	}
	return ""
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
