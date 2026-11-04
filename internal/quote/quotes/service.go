package quotes

import (
	"strings"
	"sync"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"
	"umineko_quote/internal/repository"
)

type (
	Service[Q any] interface {
		Search(p params.SearchParams, scope params.QuoteScope) (dto.SearchResponse, error)
		Browse(p params.BrowseParams, scope params.QuoteScope) (dto.CharacterResponse, error)
		Random(lang language.Language, characterID string, episode int, scope params.QuoteScope) (*Q, error)
		GetByAudioID(lang language.Language, audioID string) (*Q, error)
		GetByIndex(lang language.Language, index int) (*Q, error)
		GetContext(lang language.Language, audioID string, lines int) (*dto.ContextResponse, error)
		NearestVoicedAudioID(lang language.Language, audioID string, direction string) (string, error)
		LoadedLanguages() (map[language.Language]int, error)
	}

	service[Q any] struct {
		repo        repository.QuoteRepository[Q]
		resolveName func(string) string
	}
)

const (
	defaultSearchLimit  = 30
	defaultBrowseLimit  = 50
	defaultContextLines = 5
	maxContextLines     = 20
)

func NewService[Q any](repo repository.QuoteRepository[Q], resolveName func(string) string) Service[Q] {
	return &service[Q]{
		repo:        repo,
		resolveName: resolveName,
	}
}

func (s *service[Q]) Search(p params.SearchParams, scope params.QuoteScope) (dto.SearchResponse, error) {
	if p.Lang == language.Auto {
		return s.searchAuto(p, scope)
	}

	limit := p.Limit
	if limit <= 0 {
		limit = defaultSearchLimit
	}
	offset := max(p.Offset, 0)

	found, total, err := s.repo.FindQuotes(params.QuoteFilter{
		QuoteScope:   scope,
		Lang:         p.Lang,
		Text:         p.Query,
		Exact:        p.Exact,
		CharacterID:  p.CharacterID,
		Episode:      p.Episode,
		InteractionA: p.InteractionA,
		InteractionB: p.InteractionB,
	}, limit, offset)
	if err != nil {
		return dto.SearchResponse{}, err
	}

	results := make([]dto.SearchResult, len(found))
	for i, q := range found {
		results[i] = dto.SearchResult{Quote: q, Score: 100}
	}

	return dto.SearchResponse{
		Results: results,
		Total:   total,
		Limit:   limit,
		Offset:  offset,
	}, nil
}

func (s *service[Q]) searchAuto(p params.SearchParams, scope params.QuoteScope) (dto.SearchResponse, error) {
	type langResult struct {
		result dto.SearchResponse
		err    error
	}

	results := make([]langResult, len(language.All))
	var wg sync.WaitGroup
	for i, lang := range language.All {
		wg.Go(func() {
			sp := p
			sp.Lang = lang
			result, err := s.Search(sp, scope)
			results[i] = langResult{result: result, err: err}
		})
	}
	wg.Wait()

	for i, lr := range results {
		if lr.err != nil {
			return dto.SearchResponse{}, lr.err
		}
		if lr.result.Total > 0 {
			lr.result.Lang = string(language.All[i])
			return lr.result, nil
		}
	}

	return results[0].result, nil
}

func (s *service[Q]) Browse(p params.BrowseParams, scope params.QuoteScope) (dto.CharacterResponse, error) {
	limit := p.Limit
	if limit <= 0 {
		limit = defaultBrowseLimit
	}
	offset := max(p.Offset, 0)

	found, total, err := s.repo.FindQuotes(params.QuoteFilter{
		QuoteScope:   scope,
		Lang:         p.Lang,
		CharacterID:  p.CharacterID,
		Episode:      p.Episode,
		InteractionA: p.InteractionA,
		InteractionB: p.InteractionB,
	}, limit, offset)
	if err != nil {
		return dto.CharacterResponse{}, err
	}

	var characterName string
	if p.CharacterID != "" {
		characterName = s.resolveName(p.CharacterID)
	}

	return dto.CharacterResponse{
		CharacterID: p.CharacterID,
		Character:   characterName,
		Quotes:      found,
		Total:       total,
		Limit:       limit,
		Offset:      offset,
	}, nil
}

func (s *service[Q]) Random(lang language.Language, characterID string, episode int, scope params.QuoteScope) (*Q, error) {
	return s.repo.RandomQuote(params.QuoteFilter{
		QuoteScope:      scope,
		Lang:            lang,
		CharacterID:     characterID,
		Episode:         episode,
		ExcludeNarrator: characterID == "" && episode <= 0,
	})
}

func (s *service[Q]) GetByAudioID(lang language.Language, audioID string) (*Q, error) {
	idx, ok, err := s.repo.QuoteIndexByAudioID(lang, audioID)
	if err != nil || !ok {
		return nil, err
	}

	return s.repo.QuoteByIndex(lang, idx)
}

func (s *service[Q]) GetByIndex(lang language.Language, index int) (*Q, error) {
	return s.repo.QuoteByIndex(lang, index)
}

func (s *service[Q]) GetContext(lang language.Language, audioID string, lines int) (*dto.ContextResponse, error) {
	if lines <= 0 {
		lines = defaultContextLines
	}
	lines = min(lines, maxContextLines)

	idx, ok, err := s.repo.QuoteIndexByAudioID(lang, audioID)
	if err != nil || !ok {
		return nil, err
	}

	from := max(idx-lines, 0)

	between, err := s.repo.QuotesBetween(lang, from, idx+lines)
	if err != nil {
		return nil, err
	}

	pos := idx - from
	if pos >= len(between) {
		return nil, nil
	}

	return &dto.ContextResponse{
		Before: between[:pos],
		Quote:  between[pos],
		After:  between[pos+1:],
	}, nil
}

func (s *service[Q]) NearestVoicedAudioID(lang language.Language, audioID string, direction string) (string, error) {
	idx, ok, err := s.repo.QuoteIndexByAudioID(lang, audioID)
	if err != nil || !ok {
		return "", err
	}

	voiced, err := s.repo.NearestVoiced(lang, idx, direction == "prev")
	if err != nil || voiced == "" {
		return "", err
	}

	return strings.Split(voiced, ", ")[0], nil
}

func (s *service[Q]) LoadedLanguages() (map[language.Language]int, error) {
	return s.repo.LanguageCounts()
}
