package quote

import (
	"embed"
	"math/rand/v2"
	"strings"
	"sync"

	"github.com/sahilm/fuzzy"
)

const audioDir = "internal/quote/data/audio"

//go:embed data/*.txt
var dataFS embed.FS

type (
	Service interface {
		Search(query string, lang string, limit int, offset int, characterID string, episode int, forceFuzzy bool, truth Truth) SearchResponse
		Browse(lang string, characterID string, limit int, offset int, episode int, truth Truth) CharacterResponse
		GetByCharacter(lang string, characterID string, limit int, offset int, episode int, truth Truth) CharacterResponse
		GetByAudioID(lang string, audioID string) *ParsedQuote
		Random(lang string, characterID string, episode int, truth Truth) *ParsedQuote
		GetCharacters() map[string]string
		AudioFilePath(characterId string, audioId string) string
		GetStats() Stats
	}

	service struct {
		quotes     map[string][]ParsedQuote
		quoteTexts map[string][]string
		indexer    Indexer
		stats      Stats
	}

	langParseResult struct {
		lang   string
		parsed []ParsedQuote
		texts  []string
	}
)

func NewService() Service {
	p := NewParser()

	langFiles := map[string]string{
		"en": "data/english.txt",
		"ja": "data/japanese.txt",
	}

	results := make(chan langParseResult, len(langFiles))
	var wg sync.WaitGroup

	for lang, path := range langFiles {
		wg.Add(1)
		go func(lang, path string) {
			defer wg.Done()

			data, err := dataFS.ReadFile(path)
			if err != nil {
				return
			}
			lines := strings.Split(string(data), "\n")

			parsed := p.ParseAll(lines)

			texts := make([]string, len(parsed))
			for i := 0; i < len(parsed); i++ {
				texts[i] = parsed[i].Text
			}

			results <- langParseResult{
				lang:   lang,
				parsed: parsed,
				texts:  texts,
			}
		}(lang, path)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	quotes := make(map[string][]ParsedQuote)
	texts := make(map[string][]string)

	for r := range results {
		quotes[r.lang] = r.parsed
		texts[r.lang] = r.texts
	}

	return &service{
		quotes:     quotes,
		quoteTexts: texts,
		indexer:    NewIndexer(quotes, audioDir),
		stats:      NewStats(quotes["en"]),
	}
}

func (s *service) Search(query string, lang string, limit int, offset int, characterID string, episode int, forceFuzzy bool, truth Truth) SearchResponse {
	if limit <= 0 {
		limit = 30
	}
	if offset < 0 {
		offset = 0
	}
	if lang == "" {
		lang = "en"
	}

	quotes := s.quotes[lang]
	lowerTexts := s.indexer.LowerTexts(lang)
	if quotes == nil {
		return NewSearchResponse(nil, limit, offset)
	}

	matchesFilter := func(q ParsedQuote) bool {
		if characterID != "" && q.CharacterID != characterID {
			return false
		}
		if episode > 0 && q.Episode != episode {
			return false
		}
		if truth == TruthRed && !strings.Contains(q.TextHtml, "red-truth") {
			return false
		}
		if truth == TruthBlue && !strings.Contains(q.TextHtml, "blue-truth") {
			return false
		}
		return true
	}

	if !forceFuzzy {
		queryLower := strings.ToLower(query)

		searchIndices := s.indexer.FilteredIndices(lang, characterID, episode)

		var exactMatches []SearchResult
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

		if len(exactMatches) > 0 {
			return NewSearchResponse(exactMatches, limit, offset)
		}
	}

	quoteTexts := s.quoteTexts[lang]
	matches := fuzzy.Find(query, quoteTexts)
	if len(matches) == 0 {
		return NewSearchResponse(nil, limit, offset)
	}

	topScore := matches[0].Score
	relativeThreshold := topScore / 10
	minFuzzyScore := len(query) * 100

	var fuzzyResults []SearchResult
	for i := 0; i < len(matches); i++ {
		if matches[i].Score >= relativeThreshold && matches[i].Score >= minFuzzyScore {
			if matchesFilter(quotes[matches[i].Index]) {
				fuzzyResults = append(fuzzyResults, NewSearchResult(quotes[matches[i].Index], matches[i].Score))
			}
		}
	}

	return NewSearchResponse(fuzzyResults, limit, offset)
}

func (s *service) Browse(lang string, characterID string, limit int, offset int, episode int, truth Truth) CharacterResponse {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	if lang == "" {
		lang = "en"
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

	var all []ParsedQuote
	for _, idx := range source {
		q := quotes[idx]
		if truth == TruthRed && !strings.Contains(q.TextHtml, "red-truth") {
			continue
		}
		if truth == TruthBlue && !strings.Contains(q.TextHtml, "blue-truth") {
			continue
		}
		all = append(all, q)
	}

	return NewCharacterResponse(characterID, all, limit, offset)
}

func (s *service) GetByCharacter(lang string, characterID string, limit int, offset int, episode int, truth Truth) CharacterResponse {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	if lang == "" {
		lang = "en"
	}

	quotes := s.quotes[lang]
	if quotes == nil {
		return NewCharacterResponse(characterID, nil, limit, offset)
	}

	indices := s.indexer.CharacterIndices(lang, characterID)
	if len(indices) == 0 {
		return NewCharacterResponse(characterID, nil, limit, offset)
	}

	var all []ParsedQuote
	for _, idx := range indices {
		q := quotes[idx]
		if episode > 0 && q.Episode != episode {
			continue
		}
		if truth == TruthRed && !strings.Contains(q.TextHtml, "red-truth") {
			continue
		}
		if truth == TruthBlue && !strings.Contains(q.TextHtml, "blue-truth") {
			continue
		}
		all = append(all, q)
	}

	return NewCharacterResponse(characterID, all, limit, offset)
}

func (s *service) Random(lang string, characterID string, episode int, truth Truth) *ParsedQuote {
	if lang == "" {
		lang = "en"
	}

	quotes := s.quotes[lang]
	if quotes == nil || len(quotes) == 0 {
		return nil
	}

	matchesTruth := func(q ParsedQuote) bool {
		if truth == TruthRed && !strings.Contains(q.TextHtml, "red-truth") {
			return false
		}
		if truth == TruthBlue && !strings.Contains(q.TextHtml, "blue-truth") {
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

func (s *service) GetByAudioID(lang string, audioID string) *ParsedQuote {
	if lang == "" {
		lang = "en"
	}

	quotes := s.quotes[lang]
	if quotes == nil {
		return nil
	}

	for i := range quotes {
		if quotes[i].AudioID == audioID {
			return &quotes[i]
		}
		for _, id := range strings.Split(quotes[i].AudioID, ", ") {
			if id == audioID {
				return &quotes[i]
			}
		}
	}
	return nil
}

func (s *service) GetCharacters() map[string]string {
	return CharacterNames.GetAllCharacters()
}

func (s *service) AudioFilePath(characterId string, audioId string) string {
	return s.indexer.AudioFilePath(characterId, audioId)
}

func (s *service) GetStats() Stats {
	return s.stats
}
