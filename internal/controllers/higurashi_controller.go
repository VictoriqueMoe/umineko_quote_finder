package controllers

import (
	"strings"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/params"

	"github.com/gofiber/fiber/v3"
)

func (s *Service) getAllHigurashiQuoteRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupHigurashiSearchRoute,
		s.setupHigurashiRandomRoute,
		s.setupHigurashiBrowseRoute,
		s.setupHigurashiByIndexRoute,
		s.setupHigurashiByAudioIDRoute,
		s.setupHigurashiContextRoute,
		s.setupHigurashiNearestVoicedRoute,
		s.setupHigurashiCharactersRoute,
		s.setupHigurashiStatsRoute,
	}
}

func (s *Service) setupHigurashiSearchRoute(routeGroup fiber.Router) {
	routeGroup.Get("/search", s.higurashiSearch)
}

func (s *Service) setupHigurashiRandomRoute(routeGroup fiber.Router) {
	routeGroup.Get("/random", s.higurashiRandom)
}

func (s *Service) setupHigurashiBrowseRoute(routeGroup fiber.Router) {
	routeGroup.Get("/browse", s.higurashiBrowse)
}

func (s *Service) setupHigurashiByAudioIDRoute(routeGroup fiber.Router) {
	routeGroup.Get("/quote/+", s.higurashiByAudioID)
}

func (s *Service) setupHigurashiByIndexRoute(routeGroup fiber.Router) {
	routeGroup.Get("/quote/index/:index<int>", s.higurashiByIndex)
}

func (s *Service) setupHigurashiContextRoute(routeGroup fiber.Router) {
	routeGroup.Get("/context/+", s.higurashiContext)
}

func (s *Service) setupHigurashiNearestVoicedRoute(routeGroup fiber.Router) {
	routeGroup.Get("/nearest-voiced/+", s.higurashiNearestVoiced)
}

func (s *Service) setupHigurashiCharactersRoute(routeGroup fiber.Router) {
	routeGroup.Get("/characters", s.higurashiCharacters)
}

func (s *Service) setupHigurashiStatsRoute(routeGroup fiber.Router) {
	routeGroup.Get("/stats", s.higurashiStats)
}

// higurashiSearch godoc
//
//	@Summary		Search quotes
//	@Description	Search for Higurashi quotes by text query with optional filters
//	@Tags			higurashi
//	@Produce		json
//	@Param			q			query		string	true	"Search query"
//	@Param			lang		query		string	false	"Language"	default(en)
//	@Param			limit		query		int		false	"Maximum results"	default(30)
//	@Param			offset		query		int		false	"Offset for pagination"	default(0)
//	@Param			character	query		string	false	"Filter by character ID"	Enums(narrator, misc_voices, keiichi, rena, mion, satoko, rika, shion, satoshi, tomitake, takano, irie, ooishi, hanyuu, akasaka, okonogi, kasai, kimiyoshi, oryou, teppei, rina, chie, tomoe, madoka, yamaoki, fujita, natsumi, chisato, tamako, akira, miyuki, otobe, towada, riku, ouka, kumagai, nagisa, akane, arakawa, maeno)
//	@Param			interactionA	query		string	false	"Interaction filter: first character ID (requires interactionB)"	Enums(narrator, misc_voices, keiichi, rena, mion, satoko, rika, shion, satoshi, tomitake, takano, irie, ooishi, hanyuu, akasaka, okonogi, kasai, kimiyoshi, oryou, teppei, rina, chie, tomoe, madoka, yamaoki, fujita, natsumi, chisato, tamako, akira, miyuki, otobe, towada, riku, ouka, kumagai, nagisa, akane, arakawa, maeno)
//	@Param			interactionB	query		string	false	"Interaction filter: second character ID (requires interactionA)"	Enums(narrator, misc_voices, keiichi, rena, mion, satoko, rika, shion, satoshi, tomitake, takano, irie, ooishi, hanyuu, akasaka, okonogi, kasai, kimiyoshi, oryou, teppei, rina, chie, tomoe, madoka, yamaoki, fujita, natsumi, chisato, tamako, akira, miyuki, otobe, towada, riku, ouka, kumagai, nagisa, akane, arakawa, maeno)
//	@Param			episode		query		int		false	"Filter by arc number"
//	@Param			arc			query		string	false	"Filter by arc name (e.g. onikakushi, watanagashi)"
//	@Param			exact		query		bool	false	"Match whole words only"	default(false)
//	@Success		200			{object}	dto.SearchAPIResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Router			/higurashi/search [get]
func (s *Service) higurashiSearch(ctx fiber.Ctx) error {
	return handleSearch[dto.HigurashiQuote](ctx, s.HigurashiService, strings.TrimSpace, fiber.Query[int](ctx, "episode"), params.QuoteScope{Segment: ctx.Query("arc")})
}

// higurashiRandom godoc
//
//	@Summary		Get random quote
//	@Description	Returns a random Higurashi quote with optional filters
//	@Tags			higurashi
//	@Produce		json
//	@Param			lang		query		string	false	"Language"	default(en)
//	@Param			character	query		string	false	"Filter by character ID"	Enums(narrator, misc_voices, keiichi, rena, mion, satoko, rika, shion, satoshi, tomitake, takano, irie, ooishi, hanyuu, akasaka, okonogi, kasai, kimiyoshi, oryou, teppei, rina, chie, tomoe, madoka, yamaoki, fujita, natsumi, chisato, tamako, akira, miyuki, otobe, towada, riku, ouka, kumagai, nagisa, akane, arakawa, maeno)
//	@Param			episode		query		int		false	"Filter by arc number"
//	@Param			arc			query		string	false	"Filter by arc name (e.g. onikakushi, watanagashi)"
//	@Success		200			{object}	dto.HigurashiQuote
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/higurashi/random [get]
func (s *Service) higurashiRandom(ctx fiber.Ctx) error {
	return handleRandom[dto.HigurashiQuote](ctx, s.HigurashiService, strings.TrimSpace, fiber.Query[int](ctx, "episode"), params.QuoteScope{Segment: ctx.Query("arc")})
}

// higurashiBrowse godoc
//
//	@Summary		Browse quotes
//	@Description	Browse all Higurashi quotes with optional filters and pagination
//	@Tags			higurashi
//	@Produce		json
//	@Param			lang		query		string	false	"Language"	default(en)
//	@Param			character	query		string	false	"Filter by character ID"	Enums(narrator, misc_voices, keiichi, rena, mion, satoko, rika, shion, satoshi, tomitake, takano, irie, ooishi, hanyuu, akasaka, okonogi, kasai, kimiyoshi, oryou, teppei, rina, chie, tomoe, madoka, yamaoki, fujita, natsumi, chisato, tamako, akira, miyuki, otobe, towada, riku, ouka, kumagai, nagisa, akane, arakawa, maeno)
//	@Param			interactionA	query		string	false	"Interaction filter: first character ID (requires interactionB)"	Enums(narrator, misc_voices, keiichi, rena, mion, satoko, rika, shion, satoshi, tomitake, takano, irie, ooishi, hanyuu, akasaka, okonogi, kasai, kimiyoshi, oryou, teppei, rina, chie, tomoe, madoka, yamaoki, fujita, natsumi, chisato, tamako, akira, miyuki, otobe, towada, riku, ouka, kumagai, nagisa, akane, arakawa, maeno)
//	@Param			interactionB	query		string	false	"Interaction filter: second character ID (requires interactionA)"	Enums(narrator, misc_voices, keiichi, rena, mion, satoko, rika, shion, satoshi, tomitake, takano, irie, ooishi, hanyuu, akasaka, okonogi, kasai, kimiyoshi, oryou, teppei, rina, chie, tomoe, madoka, yamaoki, fujita, natsumi, chisato, tamako, akira, miyuki, otobe, towada, riku, ouka, kumagai, nagisa, akane, arakawa, maeno)
//	@Param			limit		query		int		false	"Maximum results"	default(50)
//	@Param			offset		query		int		false	"Offset for pagination"	default(0)
//	@Param			episode		query		int		false	"Filter by arc number"
//	@Param			arc			query		string	false	"Filter by arc name (e.g. onikakushi, watanagashi)"
//	@Success		200			{object}	dto.CharacterResponse
//	@Router			/higurashi/browse [get]
func (s *Service) higurashiBrowse(ctx fiber.Ctx) error {
	return handleBrowse[dto.HigurashiQuote](ctx, s.HigurashiService, strings.TrimSpace, fiber.Query[int](ctx, "episode"), params.QuoteScope{Segment: ctx.Query("arc")})
}

// higurashiByAudioID godoc
//
//	@Summary		Get quote by audio ID
//	@Description	Returns a specific Higurashi quote identified by its audio ID (may contain slashes)
//	@Tags			higurashi
//	@Produce		json
//	@Param			audioId		path		string	true	"Audio ID of the quote (e.g. ps3/s01/01/hrs010010)"
//	@Param			lang		query		string	false	"Language"	default(en)
//	@Success		200			{object}	dto.HigurashiQuote
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/higurashi/quote/{audioId} [get]
func (s *Service) higurashiByAudioID(ctx fiber.Ctx) error {
	return handleByAudioID[dto.HigurashiQuote](ctx, s.HigurashiService, ctx.Params("+"))
}

// higurashiByIndex godoc
//
//	@Summary		Get quote by index
//	@Description	Returns a specific Higurashi quote identified by its position index in the parsed script
//	@Tags			higurashi
//	@Produce		json
//	@Param			index		path		int		true	"Index of the quote in the parsed script"
//	@Param			lang		query		string	false	"Language"	default(en)
//	@Success		200			{object}	dto.HigurashiQuote
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/higurashi/quote/index/{index} [get]
func (s *Service) higurashiByIndex(ctx fiber.Ctx) error {
	return handleByIndex[dto.HigurashiQuote](ctx, s.HigurashiService)
}

// higurashiContext godoc
//
//	@Summary		Get quote context
//	@Description	Returns surrounding dialogue lines for a specific Higurashi quote
//	@Tags			higurashi
//	@Produce		json
//	@Param			audioId		path		string	true	"Audio ID of the quote (e.g. ps3/s01/01/hrs010010)"
//	@Param			lang		query		string	false	"Language"	default(en)
//	@Param			lines		query		int		false	"Number of context lines before and after"	default(5)
//	@Success		200			{object}	dto.ContextResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/higurashi/context/{audioId} [get]
func (s *Service) higurashiContext(ctx fiber.Ctx) error {
	return handleContext[dto.HigurashiQuote](ctx, s.HigurashiService, ctx.Params("+"))
}

// higurashiNearestVoiced godoc
//
//	@Summary		Find nearest voiced quote
//	@Description	Returns the audio ID of the nearest Higurashi quote with voice audio in a given direction
//	@Tags			higurashi
//	@Produce		json
//	@Param			audioId		path		string	true	"Audio ID of the reference quote (e.g. ps3/s01/01/hrs010010)"
//	@Param			lang		query		string	false	"Language"	default(en)
//	@Param			direction	query		string	false	"Direction to search"	default(next)	Enums(next, prev)
//	@Success		200			{object}	map[string]string
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/higurashi/nearest-voiced/{audioId} [get]
func (s *Service) higurashiNearestVoiced(ctx fiber.Ctx) error {
	return handleNearestVoiced[dto.HigurashiQuote](ctx, s.HigurashiService, ctx.Params("+"))
}

// higurashiCharacters godoc
//
//	@Summary		List all characters
//	@Description	Returns a map of Higurashi character IDs to character names
//	@Tags			higurashi
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/higurashi/characters [get]
func (s *Service) higurashiCharacters(ctx fiber.Ctx) error {
	ctx.Set("Cache-Control", "public, max-age=86400")
	return ctx.JSON(s.HigurashiService.GetCharacters())
}

// higurashiStats godoc
//
//	@Summary		Get quote statistics
//	@Description	Returns statistics about Higurashi quotes including top speakers and character interactions
//	@Tags			higurashi
//	@Produce		json
//	@Param			episode		query		int		false	"Filter by arc number, 0 for all"
//	@Success		200			{object}	dto.HigurashiStatsResult
//	@Router			/higurashi/stats [get]
func (s *Service) higurashiStats(ctx fiber.Ctx) error {
	episode := fiber.Query[int](ctx, "episode")
	return ctx.JSON(s.HigurashiService.GetStats().Compute(episode))
}
