package controllers

import (
	"strings"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"

	_ "umineko_quote/internal/dto"

	"github.com/gofiber/fiber/v3"
)

func (s *Service) getAllCiconiaQuoteRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupCiconiaSearchRoute,
		s.setupCiconiaRandomRoute,
		s.setupCiconiaBrowseRoute,
		s.setupCiconiaByIndexRoute,
		s.setupCiconiaByAudioIDRoute,
		s.setupCiconiaContextRoute,
		s.setupCiconiaCharactersRoute,
		s.setupCiconiaStatsRoute,
	}
}

func (s *Service) setupCiconiaSearchRoute(routeGroup fiber.Router) {
	routeGroup.Get("/search", s.ciconiaSearch)
}

func (s *Service) setupCiconiaRandomRoute(routeGroup fiber.Router) {
	routeGroup.Get("/random", s.ciconiaRandom)
}

func (s *Service) setupCiconiaBrowseRoute(routeGroup fiber.Router) {
	routeGroup.Get("/browse", s.ciconiaBrowse)
}

func (s *Service) setupCiconiaByAudioIDRoute(routeGroup fiber.Router) {
	routeGroup.Get("/quote/+", s.ciconiaByAudioID)
}

func (s *Service) setupCiconiaByIndexRoute(routeGroup fiber.Router) {
	routeGroup.Get("/quote/index/:index<int>", s.ciconiaByIndex)
}

func (s *Service) setupCiconiaContextRoute(routeGroup fiber.Router) {
	routeGroup.Get("/context/+", s.ciconiaContext)
}

func (s *Service) setupCiconiaCharactersRoute(routeGroup fiber.Router) {
	routeGroup.Get("/characters", s.ciconiaCharacters)
}

func (s *Service) setupCiconiaStatsRoute(routeGroup fiber.Router) {
	routeGroup.Get("/stats", s.ciconiaStats)
}

// ciconiaSearch godoc
//
//	@Summary		Search Ciconia quotes
//	@Description	Search for Ciconia no Naku Koro ni Phase 1 quotes by text query with optional filters
//	@Tags			ciconia
//	@Produce		json
//	@Param			q				query		string	true	"Search query"
//	@Param			lang			query		string	false	"Language"	default(en)
//	@Param			limit			query		int		false	"Maximum results"	default(30)
//	@Param			offset			query		int		false	"Offset for pagination"	default(0)
//	@Param			character		query		string	false	"Filter by character ID (e.g. miyao, jayden, keropoyo)"
//	@Param			interactionA	query		string	false	"Interaction filter: first character ID (requires interactionB)"
//	@Param			interactionB	query		string	false	"Interaction filter: second character ID (requires interactionA)"
//	@Param			chapter			query		string	false	"Filter by chapter ID (00 for prologue, 01-25 for main acts, 25b for finale, df01-df16 for data fragments)"
//	@Param			exact			query		bool	false	"Match whole words only"	default(false)
//	@Success		200				{object}	dto.SearchAPIResponse
//	@Failure		400				{object}	dto.ErrorResponse
//	@Router			/ciconia/search [get]
func (s *Service) ciconiaSearch(ctx fiber.Ctx) error {
	query := ctx.Query("q")
	if query == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query parameter 'q' is required",
		})
	}

	lang := language.English.Parse(ctx.Query("lang"))
	limit := fiber.Query[int](ctx, "limit", 30)
	offset := fiber.Query[int](ctx, "offset", 0)
	characterParam := strings.TrimSpace(ctx.Query("character"))
	interactionAParam := ctx.Query("interactionA")
	interactionBParam := ctx.Query("interactionB")
	exact := fiber.Query[bool](ctx, "exact", false)

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
		0,
		strings.TrimSpace(interactionAParam),
		strings.TrimSpace(interactionBParam),
		exact,
	)

	response := s.CiconiaService.Search(searchParams, ctx.Query("chapter"))
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

// ciconiaRandom godoc
//
//	@Summary		Get random Ciconia quote
//	@Description	Returns a random Ciconia quote with optional filters
//	@Tags			ciconia
//	@Produce		json
//	@Param			lang		query		string	false	"Language"	default(en)
//	@Param			character	query		string	false	"Filter by character ID"
//	@Param			chapter		query		string	false	"Filter by chapter ID"
//	@Success		200			{object}	dto.CiconiaQuote
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/ciconia/random [get]
func (s *Service) ciconiaRandom(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	characterID := strings.TrimSpace(ctx.Query("character"))
	chapter := ctx.Query("chapter")
	q := s.CiconiaService.Random(lang, characterID, chapter)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no quotes available",
		})
	}
	return ctx.JSON(q)
}

// ciconiaBrowse godoc
//
//	@Summary		Browse Ciconia quotes
//	@Description	Browse all Ciconia quotes with optional filters and pagination
//	@Tags			ciconia
//	@Produce		json
//	@Param			lang			query		string	false	"Language"	default(en)
//	@Param			character		query		string	false	"Filter by character ID"
//	@Param			interactionA	query		string	false	"Interaction filter: first character ID (requires interactionB)"
//	@Param			interactionB	query		string	false	"Interaction filter: second character ID (requires interactionA)"
//	@Param			limit			query		int		false	"Maximum results"	default(50)
//	@Param			offset			query		int		false	"Offset for pagination"	default(0)
//	@Param			chapter			query		string	false	"Filter by chapter ID"
//	@Success		200				{object}	dto.CharacterResponse
//	@Router			/ciconia/browse [get]
func (s *Service) ciconiaBrowse(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	limit := fiber.Query[int](ctx, "limit", 50)
	offset := fiber.Query[int](ctx, "offset", 0)
	characterParam := strings.TrimSpace(ctx.Query("character"))
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
		0,
		strings.TrimSpace(interactionAParam),
		strings.TrimSpace(interactionBParam),
	)

	response := s.CiconiaService.Browse(browseParams, ctx.Query("chapter"))
	return ctx.JSON(response)
}

// ciconiaByAudioID godoc
//
//	@Summary		Get Ciconia quote by synthetic ID
//	@Description	Returns a specific Ciconia quote identified by its synthetic audio ID (e.g. c01:a3f2b81c, pro:xxxxxxxx, df03:xxxxxxxx)
//	@Tags			ciconia
//	@Produce		json
//	@Param			audioId		path		string	true	"Synthetic audio ID (may contain colons)"
//	@Param			lang		query		string	false	"Language"	default(en)
//	@Success		200			{object}	dto.CiconiaQuote
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/ciconia/quote/{audioId} [get]
func (s *Service) ciconiaByAudioID(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	audioID := ctx.Params("+")

	q := s.CiconiaService.GetByAudioID(lang, audioID)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(q)
}

// ciconiaByIndex godoc
//
//	@Summary		Get Ciconia quote by index
//	@Description	Returns a specific Ciconia quote identified by its position index in the parsed script
//	@Tags			ciconia
//	@Produce		json
//	@Param			index		path		int		true	"Index of the quote in the parsed script"
//	@Param			lang		query		string	false	"Language"	default(en)
//	@Success		200			{object}	dto.CiconiaQuote
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/ciconia/quote/index/{index} [get]
func (s *Service) ciconiaByIndex(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	index := fiber.Params[int](ctx, "index")

	q := s.CiconiaService.GetByIndex(lang, index)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(q)
}

// ciconiaContext godoc
//
//	@Summary		Get Ciconia quote context
//	@Description	Returns surrounding dialogue lines for a specific Ciconia quote
//	@Tags			ciconia
//	@Produce		json
//	@Param			audioId		path		string	true	"Synthetic audio ID"
//	@Param			lang		query		string	false	"Language"	default(en)
//	@Param			lines		query		int		false	"Number of context lines before and after"	default(5)
//	@Success		200			{object}	dto.ContextResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/ciconia/context/{audioId} [get]
func (s *Service) ciconiaContext(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	audioID := ctx.Params("+")

	lines := fiber.Query[int](ctx, "lines", 5)
	result := s.CiconiaService.GetContext(lang, audioID, lines)
	if result == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(result)
}

// ciconiaCharacters godoc
//
//	@Summary		List all Ciconia characters
//	@Description	Returns main cast (curated roster of named speakers) and additional (ensemble roles and other script markers) character maps
//	@Tags			ciconia
//	@Produce		json
//	@Success		200	{object}	dto.CharactersResult
//	@Router			/ciconia/characters [get]
func (s *Service) ciconiaCharacters(ctx fiber.Ctx) error {
	ctx.Set("Cache-Control", "public, max-age=86400")
	return ctx.JSON(s.CiconiaService.GetCharacters())
}

// ciconiaStats godoc
//
//	@Summary		Get Ciconia quote statistics
//	@Description	Returns statistics about Ciconia quotes including top speakers, lines per chapter, and character interactions
//	@Tags			ciconia
//	@Produce		json
//	@Success		200	{object}	dto.CiconiaStatsResult
//	@Router			/ciconia/stats [get]
func (s *Service) ciconiaStats(ctx fiber.Ctx) error {
	return ctx.JSON(s.CiconiaService.GetStats().Compute(0))
}
