package ciconia

import (
	"log"
	"time"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/store"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
	"github.com/VictoriqueMoe/umineko_script_parser/ciconia"
	cicharacter "github.com/VictoriqueMoe/umineko_script_parser/ciconia/character"
	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
)

type (
	Service interface {
		Search(params params.SearchParams, chapter string) dto.SearchResponse
		Browse(params params.BrowseParams, chapter string) dto.CharacterResponse
		GetByAudioID(lang language.Language, audioID string) *dto.CiconiaQuote
		GetByIndex(lang language.Language, index int) *dto.CiconiaQuote
		GetContext(lang language.Language, audioID string, lines int) *dto.ContextResponse
		Random(lang language.Language, characterID string, chapter string) *dto.CiconiaQuote
		GetCharacters() dto.CharactersResult
		GetStats() store.Stats
		LoadedLanguages() map[language.Language]int
	}

	ciconiaService struct {
		store *store.GameStore[dto.CiconiaQuote]
		stats store.Stats
	}
)

func NewService() Service {
	serviceStart := time.Now()

	f, err := quote.DataFS.Open("data/Ciconia/en.file")
	if err != nil {
		log.Printf("[ciconia] failed to open data/Ciconia/en.file: %v", err)
		return &ciconiaService{
			store: store.NewGameStore(
				map[language.Language][]dto.CiconiaQuote{},
				"",
				func(q *dto.CiconiaQuote) *scriptdto.BaseQuote { return &q.BaseQuote },
			),
			stats: NewCiconiaStats(nil),
		}
	}
	defer f.Close()

	timeStart := time.Now()
	parsed, _, err := scriptparser.ParseReader(f, ciconia.NewParser())
	timeEnd := time.Now()

	log.Printf("[ciconia/en] parsed %d quotes took %v", len(parsed), timeEnd.Sub(timeStart))

	if err != nil {
		log.Printf("[ciconia] failed to parse: %v", err)
		return &ciconiaService{
			store: store.NewGameStore(
				map[language.Language][]dto.CiconiaQuote{},
				"",
				func(q *dto.CiconiaQuote) *scriptdto.BaseQuote { return &q.BaseQuote },
			),
			stats: NewCiconiaStats(nil),
		}
	}

	indexed := make([]dto.CiconiaQuote, len(parsed))
	for i := 0; i < len(parsed); i++ {
		indexed[i] = dto.CiconiaQuote{
			CicroniaQuote: parsed[i],
			Index:         i,
		}
	}

	quotes := map[language.Language][]dto.CiconiaQuote{
		language.English: indexed,
	}

	gameStore := store.NewGameStore(quotes, "", func(q *dto.CiconiaQuote) *scriptdto.BaseQuote {
		return &q.BaseQuote
	})

	log.Printf("[ciconia] initialised in %v", time.Since(serviceStart).Round(time.Millisecond))

	return &ciconiaService{
		store: gameStore,
		stats: NewCiconiaStats(indexed),
	}
}

func (s *ciconiaService) Search(p params.SearchParams, chapter string) dto.SearchResponse {
	var filter func(*dto.CiconiaQuote) bool
	if chapter != "" {
		filter = func(q *dto.CiconiaQuote) bool {
			return q.Chapter == chapter
		}
	}
	return s.store.Search(p, filter)
}

func (s *ciconiaService) Browse(p params.BrowseParams, chapter string) dto.CharacterResponse {
	var filter func(*dto.CiconiaQuote) bool
	if chapter != "" {
		filter = func(q *dto.CiconiaQuote) bool {
			return q.Chapter == chapter
		}
	}
	return s.store.Browse(p, filter, func(id string) string {
		return cicharacter.CharacterNames.GetCharacterName(cicharacter.Character(id))
	})
}

func (s *ciconiaService) Random(lang language.Language, characterID string, chapter string) *dto.CiconiaQuote {
	var filter func(*dto.CiconiaQuote) bool
	if chapter != "" {
		filter = func(q *dto.CiconiaQuote) bool {
			return q.Chapter == chapter
		}
	}
	return s.store.Random(lang, characterID, 0, filter)
}

func (s *ciconiaService) GetByAudioID(lang language.Language, audioID string) *dto.CiconiaQuote {
	return s.store.GetByAudioID(lang, audioID)
}

func (s *ciconiaService) GetByIndex(lang language.Language, index int) *dto.CiconiaQuote {
	return s.store.GetByIndex(lang, index)
}

func (s *ciconiaService) GetContext(lang language.Language, audioID string, lines int) *dto.ContextResponse {
	return s.store.GetContext(lang, audioID, lines)
}

func (s *ciconiaService) GetCharacters() dto.CharactersResult {
	curated := cicharacter.MainCharacterNames.GetAllCharacters()
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

func (s *ciconiaService) GetStats() store.Stats {
	return s.stats
}

func (s *ciconiaService) LoadedLanguages() map[language.Language]int {
	return s.store.LoadedLanguages()
}
