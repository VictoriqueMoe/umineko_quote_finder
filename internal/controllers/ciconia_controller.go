package controllers

import (
	"strings"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/params"

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
	return handleSearch[dto.CiconiaQuote](ctx, s.CiconiaService, strings.TrimSpace, 0, params.QuoteScope{Segment: ctx.Query("chapter")})
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
	return handleRandom[dto.CiconiaQuote](ctx, s.CiconiaService, strings.TrimSpace, 0, params.QuoteScope{Segment: ctx.Query("chapter")})
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
	return handleBrowse[dto.CiconiaQuote](ctx, s.CiconiaService, strings.TrimSpace, 0, params.QuoteScope{Segment: ctx.Query("chapter")})
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
	return handleByAudioID[dto.CiconiaQuote](ctx, s.CiconiaService, ctx.Params("+"))
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
	return handleByIndex[dto.CiconiaQuote](ctx, s.CiconiaService)
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
	return handleContext[dto.CiconiaQuote](ctx, s.CiconiaService, ctx.Params("+"))
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
