package store

import (
	"math/rand/v2"
	"strings"
	"sync"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
)

type GameStore[Q any] struct {
	quotes  map[language.Language][]Q
	indexer Indexer
	base    func(*Q) *scriptdto.BaseQuote
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

func NewGameStore[Q any](quotes map[language.Language][]Q, audioDir string, base func(*Q) *scriptdto.BaseQuote) *GameStore[Q] {
	fields := make(map[language.Language]QuoteFields, len(quotes))
	for lang, qs := range quotes {
		n := len(qs)
		f := QuoteFields{
			Texts:        make([]string, n),
			CharacterIDs: make([]string, n),
			Episodes:     make([]int, n),
			AudioIDs:     make([]string, n),
		}
		for i := 0; i < n; i++ {
			bq := base(&qs[i])
			f.Texts[i] = bq.Text
			f.CharacterIDs[i] = bq.CharacterID
			f.Episodes[i] = bq.Episode
			f.AudioIDs[i] = bq.AudioID
		}
		fields[lang] = f
	}

	return &GameStore[Q]{
		quotes:  quotes,
		indexer: NewIndexer(fields, audioDir),
		base:    base,
	}
}

func (gs *GameStore[Q]) Indexer() Indexer {
	return gs.indexer
}

func (gs *GameStore[Q]) Quotes(lang language.Language) []Q {
	return gs.quotes[lang]
}

func (gs *GameStore[Q]) Base(q *Q) *scriptdto.BaseQuote {
	return gs.base(q)
}

func (gs *GameStore[Q]) Search(p params.SearchParams, filter func(*Q) bool) dto.SearchResponse {
	if p.Lang == language.Auto {
		return gs.searchAuto(p, filter)
	}

	query := p.Query
	lang := p.Lang
	limit := p.Limit
	offset := p.Offset
	characterID := p.CharacterID
	episode := p.Episode
	interactionA := p.InteractionA
	interactionB := p.InteractionB

	if limit <= 0 {
		limit = 30
	}
	if offset < 0 {
		offset = 0
	}

	quotes := gs.quotes[lang]
	lowerTexts := gs.indexer.LowerTexts(lang)
	if quotes == nil {
		return NewSearchResponse(nil, limit, offset)
	}

	matchesFilter := func(idx int) bool {
		if filter != nil && !filter(&quotes[idx]) {
			return false
		}
		return true
	}

	queryLower := strings.ToLower(query)

	searchIndices := gs.indexer.FilteredIndices(lang, characterID, episode)
	if interactionA != "" && interactionB != "" {
		interactionIndices := gs.indexer.InteractionIndices(lang, interactionA, interactionB)
		searchIndices = mergeFilteredIndices(searchIndices, interactionIndices)
	}

	var exactMatches []dto.SearchResult
	if searchIndices != nil {
		if len(searchIndices) > 5000 {
			exactMatches = concurrentExactSearchGeneric(searchIndices, lowerTexts, quotes, queryLower, matchesFilter)
		} else {
			for _, idx := range searchIndices {
				if strings.Contains(lowerTexts[idx], queryLower) {
					if matchesFilter(idx) {
						exactMatches = append(exactMatches, dto.SearchResult{Quote: quotes[idx], Score: 100})
					}
				}
			}
		}
	} else {
		allIndices := make([]int, len(quotes))
		for i := range allIndices {
			allIndices[i] = i
		}
		exactMatches = concurrentExactSearchGeneric(allIndices, lowerTexts, quotes, queryLower, matchesFilter)
	}

	return NewSearchResponse(exactMatches, limit, offset)
}

func (gs *GameStore[Q]) searchAuto(p params.SearchParams, filter func(*Q) bool) dto.SearchResponse {
	type langResult struct {
		lang   language.Language
		result dto.SearchResponse
	}

	results := make([]langResult, len(language.All))
	var wg sync.WaitGroup
	for i, lang := range language.All {
		wg.Go(func() {
			sp := p
			sp.Lang = lang
			results[i] = langResult{lang: lang, result: gs.Search(sp, filter)}
		})
	}
	wg.Wait()

	for _, lr := range results {
		if lr.result.Total > 0 {
			lr.result.Lang = string(lr.lang)
			return lr.result
		}
	}

	p.Lang = language.English
	return gs.Search(p, filter)
}

func (gs *GameStore[Q]) Browse(p params.BrowseParams, filter func(*Q) bool, resolveName func(string) string) dto.CharacterResponse {
	lang := p.Lang
	limit := p.Limit
	offset := p.Offset
	episode := p.Episode
	interactionA := p.InteractionA
	interactionB := p.InteractionB
	characterID := p.CharacterID

	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	quotes := gs.quotes[lang]
	if quotes == nil {
		return newCharacterResponse(characterID, []Q{}, limit, offset, resolveName)
	}

	var source []int
	indexed := gs.indexer.FilteredIndices(lang, characterID, episode)
	if interactionA != "" && interactionB != "" {
		interactionIndices := gs.indexer.InteractionIndices(lang, interactionA, interactionB)
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

	var all []Q
	for _, idx := range source {
		if filter != nil && !filter(&quotes[idx]) {
			continue
		}
		all = append(all, quotes[idx])
	}

	if all == nil {
		all = []Q{}
	}

	return newCharacterResponse(characterID, all, limit, offset, resolveName)
}

func (gs *GameStore[Q]) Random(lang language.Language, characterID string, episode int, filter func(*Q) bool) *Q {
	quotes := gs.quotes[lang]
	if quotes == nil || len(quotes) == 0 {
		return nil
	}

	needsFilter := filter != nil

	if characterID == "" && episode <= 0 && !needsFilter {
		indices := gs.indexer.NonNarratorIndices(lang)
		if len(indices) == 0 {
			return nil
		}
		pick := indices[rand.IntN(len(indices))]
		return &quotes[pick]
	}

	var candidates []int

	if needsFilter {
		var source []int
		indexed := gs.indexer.FilteredIndices(lang, characterID, episode)
		if indexed != nil {
			source = indexed
		} else if characterID == "" && episode <= 0 {
			source = gs.indexer.NonNarratorIndices(lang)
		}

		if source != nil {
			for _, idx := range source {
				if filter(&quotes[idx]) {
					candidates = append(candidates, idx)
				}
			}
		} else {
			for i := 0; i < len(quotes); i++ {
				bq := gs.base(&quotes[i])
				if characterID != "" && bq.CharacterID != characterID {
					continue
				}
				if episode > 0 && bq.Episode != episode {
					continue
				}
				if filter(&quotes[i]) {
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

	indices := gs.indexer.FilteredIndices(lang, characterID, episode)
	if indices != nil {
		if len(indices) == 0 {
			return nil
		}
		pick := indices[rand.IntN(len(indices))]
		return &quotes[pick]
	}

	var filtered []int
	for i := 0; i < len(quotes); i++ {
		bq := gs.base(&quotes[i])
		if characterID != "" && bq.CharacterID != characterID {
			continue
		}
		if episode > 0 && bq.Episode != episode {
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

func (gs *GameStore[Q]) GetByAudioID(lang language.Language, audioID string) *Q {
	quotes := gs.quotes[lang]
	if quotes == nil {
		return nil
	}

	idx, ok := gs.indexer.QuoteIndex(lang, audioID)
	if !ok {
		return nil
	}
	return &quotes[idx]
}

func (gs *GameStore[Q]) GetByIndex(lang language.Language, index int) *Q {
	quotes := gs.quotes[lang]
	if quotes == nil {
		return nil
	}
	if index < 0 || index >= len(quotes) {
		return nil
	}
	return &quotes[index]
}

func (gs *GameStore[Q]) GetContext(lang language.Language, audioID string, lines int) *dto.ContextResponse {
	if lines <= 0 {
		lines = 5
	}
	if lines > 20 {
		lines = 20
	}

	quotes := gs.quotes[lang]
	if quotes == nil {
		return nil
	}

	idx, ok := gs.indexer.QuoteIndex(lang, audioID)
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

	before := make([]Q, idx-start)
	copy(before, quotes[start:idx])

	after := make([]Q, end-(idx+1))
	copy(after, quotes[idx+1:end])

	return &dto.ContextResponse{
		Before: before,
		Quote:  quotes[idx],
		After:  after,
	}
}

func (gs *GameStore[Q]) NearestVoicedAudioID(lang language.Language, audioID string, direction string) string {
	quotes := gs.quotes[lang]
	if quotes == nil {
		return ""
	}

	idx, ok := gs.indexer.QuoteIndex(lang, audioID)
	if !ok {
		return ""
	}

	if direction == "prev" {
		for i := idx - 1; i >= 0; i-- {
			bq := gs.base(&quotes[i])
			if bq.AudioID != "" {
				return strings.Split(bq.AudioID, ", ")[0]
			}
		}
		return ""
	}

	for i := idx + 1; i < len(quotes); i++ {
		bq := gs.base(&quotes[i])
		if bq.AudioID != "" {
			return strings.Split(bq.AudioID, ", ")[0]
		}
	}
	return ""
}

func (gs *GameStore[Q]) AudioFilePath(characterId string, audioId string) string {
	return gs.indexer.AudioFilePath(characterId, audioId)
}

func (gs *GameStore[Q]) HasAudio() bool {
	return gs.indexer.HasAudio()
}

func (gs *GameStore[Q]) LoadedLanguages() map[language.Language]int {
	result := make(map[language.Language]int, len(gs.quotes))
	for lang, quotes := range gs.quotes {
		result[lang] = len(quotes)
	}
	return result
}

func newCharacterResponse[Q any](characterID string, quotes []Q, limit int, offset int, resolveName func(string) string) dto.CharacterResponse {
	var characterName string
	if characterID != "" && resolveName != nil {
		characterName = resolveName(characterID)
	}

	total := len(quotes)

	if offset >= total {
		return dto.CharacterResponse{
			CharacterID: characterID,
			Character:   characterName,
			Quotes:      []Q{},
			Total:       total,
			Limit:       limit,
			Offset:      offset,
		}
	}

	end := offset + limit
	if end > total {
		end = total
	}

	return dto.CharacterResponse{
		CharacterID: characterID,
		Character:   characterName,
		Quotes:      quotes[offset:end],
		Total:       total,
		Limit:       limit,
		Offset:      offset,
	}
}
