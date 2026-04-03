package higurashi

import (
	"log"
	"time"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/store"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/higurashi"
	hicharacter "github.com/VictoriqueMoe/umineko_script_parser/higurashi/character"
)

type (
	Service interface {
		Search(params params.SearchParams, arc string) dto.SearchResponse
		Browse(params params.BrowseParams, arc string) dto.CharacterResponse
		GetByAudioID(lang language.Language, audioID string) *dto.HigurashiQuote
		GetByIndex(lang language.Language, index int) *dto.HigurashiQuote
		GetContext(lang language.Language, audioID string, lines int) *dto.ContextResponse
		NearestVoicedAudioID(lang language.Language, audioID string, direction string) string
		Random(lang language.Language, characterID string, episode int, arc string) *dto.HigurashiQuote
		GetCharacters() dto.CharactersResult
		GetStats() store.Stats
		LoadedLanguages() map[language.Language]int
	}

	higurashiService struct {
		store *store.GameStore[dto.HigurashiQuote]
		stats store.Stats
	}
)

func NewService() Service {
	serviceStart := time.Now()

	f, err := quote.DataFS.Open("data/higurashi/en.file")
	if err != nil {
		log.Printf("[higurashi] failed to open data/higurashi/en.file: %v", err)
		return &higurashiService{
			store: store.NewGameStore(
				map[language.Language][]dto.HigurashiQuote{},
				"",
				func(q *dto.HigurashiQuote) *scriptdto.BaseQuote { return &q.BaseQuote },
			),
			stats: NewHigurashiStats(nil),
		}
	}
	defer f.Close()

	timeStart := time.Now()
	parsed, _, err := higurashi.ParseFile(f)
	timeEnd := time.Now()

	log.Printf("[higurashi/en] parsed %d quotes took %v", len(parsed), timeEnd.Sub(timeStart))

	if err != nil {
		log.Printf("[higurashi] failed to parse: %v", err)
		return &higurashiService{
			store: store.NewGameStore(
				map[language.Language][]dto.HigurashiQuote{},
				"",
				func(q *dto.HigurashiQuote) *scriptdto.BaseQuote { return &q.BaseQuote },
			),
			stats: NewHigurashiStats(nil),
		}
	}

	indexed := make([]dto.HigurashiQuote, len(parsed))
	for i := 0; i < len(parsed); i++ {
		indexed[i] = dto.HigurashiQuote{
			HigurashiQuote: parsed[i],
			Index:          i,
		}
	}

	quotes := map[language.Language][]dto.HigurashiQuote{
		language.English: indexed,
	}

	gameStore := store.NewGameStore(quotes, "", func(q *dto.HigurashiQuote) *scriptdto.BaseQuote {
		return &q.BaseQuote
	})

	log.Printf("[higurashi] initialised in %v", time.Since(serviceStart).Round(time.Millisecond))

	return &higurashiService{
		store: gameStore,
		stats: NewHigurashiStats(indexed),
	}
}

func (s *higurashiService) Search(p params.SearchParams, arc string) dto.SearchResponse {
	var filter func(*dto.HigurashiQuote) bool
	if arc != "" {
		filter = func(q *dto.HigurashiQuote) bool {
			return q.Arc == arc
		}
	}
	return s.store.Search(p, filter)
}

func (s *higurashiService) Browse(p params.BrowseParams, arc string) dto.CharacterResponse {
	var filter func(*dto.HigurashiQuote) bool
	if arc != "" {
		filter = func(q *dto.HigurashiQuote) bool {
			return q.Arc == arc
		}
	}
	return s.store.Browse(p, filter, func(id string) string {
		return hicharacter.CharacterNames.GetCharacterName(hicharacter.CharacterFromID(id))
	})
}

func (s *higurashiService) Random(lang language.Language, characterID string, episode int, arc string) *dto.HigurashiQuote {
	var filter func(*dto.HigurashiQuote) bool
	if arc != "" {
		filter = func(q *dto.HigurashiQuote) bool {
			return q.Arc == arc
		}
	}
	return s.store.Random(lang, characterID, episode, filter)
}

func (s *higurashiService) GetByAudioID(lang language.Language, audioID string) *dto.HigurashiQuote {
	return s.store.GetByAudioID(lang, audioID)
}

func (s *higurashiService) GetByIndex(lang language.Language, index int) *dto.HigurashiQuote {
	return s.store.GetByIndex(lang, index)
}

func (s *higurashiService) GetContext(lang language.Language, audioID string, lines int) *dto.ContextResponse {
	return s.store.GetContext(lang, audioID, lines)
}

func (s *higurashiService) NearestVoicedAudioID(lang language.Language, audioID string, direction string) string {
	return s.store.NearestVoicedAudioID(lang, audioID, direction)
}

func (s *higurashiService) GetCharacters() dto.CharactersResult {
	curated := hicharacter.CharacterNames.GetAllCharacters()
	curatedIDs := make(map[string]bool, len(curated))
	characters := make(map[string]string, len(curated))
	for ch, name := range curated {
		id := ch.ID()
		characters[id] = name
		curatedIDs[id] = true
	}

	additional := map[string]string{}
	quotes := s.store.Quotes(language.English)
	for i := 0; i < len(quotes); i++ {
		base := s.store.Base(&quotes[i])
		cid := base.CharacterID
		if cid == "" || curatedIDs[cid] {
			continue
		}
		if _, exists := additional[cid]; !exists {
			additional[cid] = base.Character
		}
	}

	return dto.CharactersResult{
		Characters: characters,
		Additional: additional,
	}
}

func (s *higurashiService) GetStats() store.Stats {
	return s.stats
}

func (s *higurashiService) LoadedLanguages() map[language.Language]int {
	return s.store.LoadedLanguages()
}
