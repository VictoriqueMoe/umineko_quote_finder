package controllers

import (
	"os"
	"regexp"
	"strings"
	"umineko_quote/internal/quote/character"
	"umineko_quote/internal/quote/language"
	quoteparams "umineko_quote/internal/quote/params"

	"umineko_quote/internal/audio"
	_ "umineko_quote/internal/dto"
	"umineko_quote/internal/quote"
	"umineko_quote/internal/utils"

	"github.com/gofiber/fiber/v3"
)

var audioIdPattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func (s *Service) getAllQuoteRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupSearchRoute,
		s.setupRandomRoute,
		s.setupBrowseRoute,
		s.setupByAudioIDRoute,
		s.setupContextRoute,
		s.setupCharactersRoute,
		s.setupCombinedAudioRoute,
		s.setupAudioRoute,
		s.setupStatsRoute,
	}
}

func (s *Service) setupSearchRoute(routeGroup fiber.Router) {
	routeGroup.Get("/search", s.search)
}

func (s *Service) setupRandomRoute(routeGroup fiber.Router) {
	routeGroup.Get("/random", s.random)
}

func (s *Service) setupBrowseRoute(routeGroup fiber.Router) {
	routeGroup.Get("/browse", s.browse)
}

func (s *Service) setupByAudioIDRoute(routeGroup fiber.Router) {
	routeGroup.Get("/quote/:audioId", s.byAudioID)
}

func (s *Service) setupCharactersRoute(routeGroup fiber.Router) {
	routeGroup.Get("/characters", s.characters)
}

// search godoc
//
//	@Summary		Search quotes
//	@Description	Search for quotes by text query with optional filters
//	@Tags			quotes
//	@Produce		json
//	@Param			q			query		string	true	"Search query"
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(auto, en, wh, ja, ru, es, pt)
//	@Param			limit		query		int		false	"Maximum results"	default(30)
//	@Param			offset		query		int		false	"Offset for pagination"	default(0)
//	@Param			character	query		character.Character	false	"Filter by character ID"
//	@Param			interactionA	query		character.Character	false	"Interaction filter: first character ID (requires interactionB)"
//	@Param			interactionB	query		character.Character	false	"Interaction filter: second character ID (requires interactionA)"
//	@Param			episode		query		int		false	"Filter by episode (1-8)"
//	@Param			truth		query		string	false	"Filter by truth type"	Enums(red, blue)
//	@Success		200			{object}	dto.SearchAPIResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Router			/search [get]
func (s *Service) search(ctx fiber.Ctx) error {
	query := ctx.Query("q")
	if query == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query parameter 'q' is required",
		})
	}

	lang := language.English.Parse(ctx.Query("lang"))
	limit := fiber.Query[int](ctx, "limit", 30)
	offset := fiber.Query[int](ctx, "offset", 0)
	characterParam := ctx.Query("character")
	interactionAParam := ctx.Query("interactionA")
	interactionBParam := ctx.Query("interactionB")
	episode := fiber.Query[int](ctx, "episode", 0)

	if errMsg, invalid := validateInteractionQueryParams(characterParam, interactionAParam, interactionBParam); invalid {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errMsg,
		})
	}

	searchParams := quoteparams.NewSearchParams(
		query,
		lang,
		limit,
		offset,
		characterParam,
		episode,
		ctx.Query("truth"),
		interactionAParam,
		interactionBParam,
	)

	response := s.QuoteService.Search(searchParams)
	result := fiber.Map{
		"query":   query,
		"results": response.Results,
		"total":   response.Total,
		"limit":   response.Limit,
		"offset":  response.Offset,
	}
	if response.Lang != "" {
		result["lang"] = response.Lang
	}
	return ctx.JSON(result)
}

// random godoc
//
//	@Summary		Get random quote
//	@Description	Returns a random quote with optional filters
//	@Tags			quotes
//	@Produce		json
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(en, wh, ja, ru, es, pt)
//	@Param			character	query		character.Character	false	"Filter by character ID"
//	@Param			episode		query		int		false	"Filter by episode (1-8)"
//	@Param			truth		query		string	false	"Filter by truth type"	Enums(red, blue)
//	@Success		200			{object}	dto.ParsedQuote
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/random [get]
func (s *Service) random(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	characterStruct := character.Character(ctx.Query("character"))
	episode := fiber.Query[int](ctx, "episode", 0)
	truth := quote.TruthAll.Parse(ctx.Query("truth"))
	q := s.QuoteService.Random(lang, characterStruct, episode, truth)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no quotes available",
		})
	}
	return ctx.JSON(q)
}

// browse godoc
//
//	@Summary		Browse quotes
//	@Description	Browse all quotes with optional filters and pagination
//	@Tags			quotes
//	@Produce		json
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(en, wh, ja, ru, es, pt)
//	@Param			character	query		character.Character	false	"Filter by character ID"
//	@Param			interactionA	query		character.Character	false	"Interaction filter: first character ID (requires interactionB)"
//	@Param			interactionB	query		character.Character	false	"Interaction filter: second character ID (requires interactionA)"
//	@Param			limit		query		int		false	"Maximum results"	default(50)
//	@Param			offset		query		int		false	"Offset for pagination"	default(0)
//	@Param			episode		query		int		false	"Filter by episode (1-8)"
//	@Param			truth		query		string	false	"Filter by truth type"	Enums(red, blue)
//	@Success		200			{object}	dto.CharacterResponse
//	@Router			/browse [get]
func (s *Service) browse(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	limit := fiber.Query[int](ctx, "limit", 50)
	offset := fiber.Query[int](ctx, "offset", 0)
	episode := fiber.Query[int](ctx, "episode", 0)
	characterParam := ctx.Query("character")
	interactionAParam := ctx.Query("interactionA")
	interactionBParam := ctx.Query("interactionB")

	if errMsg, invalid := validateInteractionQueryParams(characterParam, interactionAParam, interactionBParam); invalid {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errMsg,
		})
	}

	browseParams := quoteparams.NewBrowseParams(
		lang,
		limit,
		offset,
		characterParam,
		episode,
		ctx.Query("truth"),
		interactionAParam,
		interactionBParam,
	)

	response := s.QuoteService.Browse(browseParams)
	return ctx.JSON(response)
}

// byAudioID godoc
//
//	@Summary		Get quote by audio ID
//	@Description	Returns a specific quote identified by its audio ID
//	@Tags			quotes
//	@Produce		json
//	@Param			audioId		path		string	true	"Audio ID of the quote"
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(en, wh, ja, ru, es, pt)
//	@Success		200			{object}	dto.ParsedQuote
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/quote/{audioId} [get]
func (s *Service) byAudioID(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	audioID := ctx.Params("audioId")

	q := s.QuoteService.GetByAudioID(lang, audioID)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(q)
}

func (s *Service) setupContextRoute(routeGroup fiber.Router) {
	routeGroup.Get("/context/:audioId", s.context)
}

// context godoc
//
//	@Summary		Get quote context
//	@Description	Returns surrounding dialogue lines for a specific quote
//	@Tags			quotes
//	@Produce		json
//	@Param			audioId		path		string	true	"Audio ID of the quote"
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(en, wh, ja, ru, es, pt)
//	@Param			lines		query		int		false	"Number of context lines before and after"	default(5)
//	@Success		200			{object}	dto.ContextResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/context/{audioId} [get]
func (s *Service) context(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	audioID := ctx.Params("audioId")
	if !audioIdPattern.MatchString(audioID) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid audio ID",
		})
	}

	lines := fiber.Query[int](ctx, "lines", 5)
	result := s.QuoteService.GetContext(lang, audioID, lines)
	if result == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(result)
}

// characters godoc
//
//	@Summary		List all characters
//	@Description	Returns a map of character IDs to character names
//	@Tags			quotes
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/characters [get]
func (s *Service) characters(ctx fiber.Ctx) error {
	ctx.Set("Cache-Control", "public, max-age=86400")
	return ctx.JSON(s.QuoteService.GetCharacters())
}

func (s *Service) setupStatsRoute(routeGroup fiber.Router) {
	routeGroup.Get("/stats", s.stats)
}

// stats godoc
//
//	@Summary		Get quote statistics
//	@Description	Returns statistics about quotes including top speakers, truth per episode, and character interactions
//	@Tags			quotes
//	@Produce		json
//	@Param			episode		query		int		false	"Filter by episode (1-8), 0 for all"
//	@Success		200			{object}	dto.StatsResult
//	@Router			/stats [get]
func (s *Service) stats(ctx fiber.Ctx) error {
	episode := fiber.Query[int](ctx, "episode", 0)
	return ctx.JSON(s.QuoteService.GetStats().Compute(episode))
}

func (s *Service) setupCombinedAudioRoute(routeGroup fiber.Router) {
	routeGroup.Get("/audio/combined", s.combinedAudioSegments)
	routeGroup.Get("/audio/:charId/combined", s.combinedAudioLegacy)
}

func (s *Service) setupAudioRoute(routeGroup fiber.Router) {
	routeGroup.Get("/audio/:charId/:audioId", s.audio)
}

func (s *Service) audio(ctx fiber.Ctx) error {
	charId := ctx.Params("charId")
	audioId := ctx.Params("audioId")
	if !audioIdPattern.MatchString(charId) || !audioIdPattern.MatchString(audioId) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid audio ID",
		})
	}

	filePath := s.QuoteService.AudioFilePath(charId, audioId)
	if filePath == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "audio file not found",
		})
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read audio file",
		})
	}

	return utils.ServeAudio(ctx, data)
}

func (s *Service) combinedAudioSegments(ctx fiber.Ctx) error {
	segmentsParam := ctx.Query("segments")
	if segmentsParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query parameter 'segments' is required",
		})
	}

	parts := strings.Split(segmentsParam, ",")
	if len(parts) > 20 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "maximum 20 audio segments allowed",
		})
	}

	segments := make([]audio.AudioSegment, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		colonIdx := strings.IndexByte(part, ':')
		if colonIdx < 1 || colonIdx >= len(part)-1 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid segment format: " + part + " (expected charId:audioId)",
			})
		}
		charId := part[:colonIdx]
		audioId := part[colonIdx+1:]
		if !audioIdPattern.MatchString(charId) || !audioIdPattern.MatchString(audioId) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid segment: " + part,
			})
		}
		segments = append(segments, audio.AudioSegment{CharID: charId, AudioID: audioId})
	}

	data, err := s.AudioCombiner.CombineOgg(segments, s.QuoteService.AudioFilePath)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.ServeAudio(ctx, data)
}

func (s *Service) combinedAudioLegacy(ctx fiber.Ctx) error {
	charId := ctx.Params("charId")
	if !audioIdPattern.MatchString(charId) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid character ID",
		})
	}

	idsParam := ctx.Query("ids")
	if idsParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query parameter 'ids' is required",
		})
	}

	ids := strings.Split(idsParam, ",")
	if len(ids) > 20 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "maximum 20 audio IDs allowed",
		})
	}

	segments := make([]audio.AudioSegment, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if !audioIdPattern.MatchString(id) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid audio ID: " + id,
			})
		}
		segments = append(segments, audio.AudioSegment{CharID: charId, AudioID: id})
	}

	data, err := s.AudioCombiner.CombineOgg(segments, s.QuoteService.AudioFilePath)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.ServeAudio(ctx, data)
}

func validateInteractionQueryParams(characterParam string, interactionAParam string, interactionBParam string) (string, bool) {
	characterParam = strings.TrimSpace(characterParam)
	interactionAParam = strings.TrimSpace(interactionAParam)
	interactionBParam = strings.TrimSpace(interactionBParam)

	hasInteractionA := interactionAParam != ""
	hasInteractionB := interactionBParam != ""
	if hasInteractionA != hasInteractionB {
		return "interactionA and interactionB must both be provided", true
	}

	if characterParam != "" && hasInteractionA && hasInteractionB {
		return "character cannot be combined with interactionA/interactionB", true
	}

	return "", false
}
