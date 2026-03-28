package controllers

import (
	"regexp"
	"strings"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"

	"github.com/VictoriqueMoe/umineko_script_parser/quote/character"

	_ "umineko_quote/internal/dto"
	"umineko_quote/internal/quote"

	"github.com/gofiber/fiber/v3"
)

var audioIdPattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func (s *Service) getAllQuoteRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupSearchRoute,
		s.setupRandomRoute,
		s.setupBrowseRoute,
		s.setupByAudioIDRoute,
		s.setupByIndexRoute,
		s.setupContextRoute,
		s.setupNearestVoicedRoute,
		s.setupCharactersRoute,
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

func (s *Service) setupByIndexRoute(routeGroup fiber.Router) {
	routeGroup.Get("/quote/index/:index<int>", s.byIndex)
}

func (s *Service) setupContextRoute(routeGroup fiber.Router) {
	routeGroup.Get("/context/:audioId", s.context)
}

func (s *Service) setupNearestVoicedRoute(routeGroup fiber.Router) {
	routeGroup.Get("/nearest-voiced/:audioId", s.nearestVoiced)
}

func (s *Service) setupCharactersRoute(routeGroup fiber.Router) {
	routeGroup.Get("/characters", s.characters)
}

func (s *Service) setupStatsRoute(routeGroup fiber.Router) {
	routeGroup.Get("/stats", s.stats)
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
//	@Param			truth		query		string	false	"Filter by truth type"	Enums(red, blue, gold, purple)
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

	searchParams := params.NewSearchParams(
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
//	@Param			truth		query		string	false	"Filter by truth type"	Enums(red, blue, gold, purple)
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
//	@Param			truth		query		string	false	"Filter by truth type"	Enums(red, blue, gold, purple)
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

	browseParams := params.NewBrowseParams(
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

// byIndex godoc
//
//	@Summary		Get quote by index
//	@Description	Returns a specific quote identified by its position index in the parsed script
//	@Tags			quotes
//	@Produce		json
//	@Param			index		path		int		true	"Index of the quote in the parsed script"
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(en, wh, ja, ru, es, pt)
//	@Success		200			{object}	dto.ParsedQuote
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/quote/index/{index} [get]
func (s *Service) byIndex(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	index := fiber.Params[int](ctx, "index")

	q := s.QuoteService.GetByIndex(lang, index)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(q)
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

// nearestVoiced godoc
//
//	@Summary		Find nearest voiced quote
//	@Description	Returns the audio ID of the nearest quote with voice audio in a given direction
//	@Tags			quotes
//	@Produce		json
//	@Param			audioId		path		string	true	"Audio ID of the reference quote"
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(en, wh, ja, ru, es, pt)
//	@Param			direction	query		string	false	"Direction to search"	default(next)	Enums(next, prev)
//	@Success		200			{object}	map[string]string
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/nearest-voiced/{audioId} [get]
func (s *Service) nearestVoiced(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	audioID := ctx.Params("audioId")
	if !audioIdPattern.MatchString(audioID) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid audio ID",
		})
	}

	direction := ctx.Query("direction", "next")
	if direction != "next" && direction != "prev" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "direction must be 'next' or 'prev'",
		})
	}

	result := s.QuoteService.NearestVoicedAudioID(lang, audioID, direction)
	if result == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no voiced quote found",
		})
	}
	return ctx.JSON(fiber.Map{"audioId": result})
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
