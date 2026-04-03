package controllers

import (
	"regexp"
	"strings"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/umineko"

	"github.com/VictoriqueMoe/umineko_script_parser/umineko/character"

	"umineko_quote/internal/dto"

	"github.com/gofiber/fiber/v3"
)

var audioIdPattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func (s *Service) getAllUminekoQuoteRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupUminekoSearchRoute,
		s.setupUminekoRandomRoute,
		s.setupUminekoBrowseRoute,
		s.setupUminekoByAudioIDRoute,
		s.setupUminekoByIndexRoute,
		s.setupUminekoContextRoute,
		s.setupUminekoNearestVoicedRoute,
		s.setupUminekoCharactersRoute,
		s.setupUminekoStatsRoute,
	}
}

func (s *Service) setupUminekoSearchRoute(routeGroup fiber.Router) {
	routeGroup.Get("/search", s.uminekoSearch)
}

func (s *Service) setupUminekoRandomRoute(routeGroup fiber.Router) {
	routeGroup.Get("/random", s.uminekoRandom)
}

func (s *Service) setupUminekoBrowseRoute(routeGroup fiber.Router) {
	routeGroup.Get("/browse", s.uminekoBrowse)
}

func (s *Service) setupUminekoByAudioIDRoute(routeGroup fiber.Router) {
	routeGroup.Get("/quote/:audioId", s.uminekoByAudioID)
}

func (s *Service) setupUminekoByIndexRoute(routeGroup fiber.Router) {
	routeGroup.Get("/quote/index/:index<int>", s.uminekoByIndex)
}

func (s *Service) setupUminekoContextRoute(routeGroup fiber.Router) {
	routeGroup.Get("/context/:audioId", s.uminekoContext)
}

func (s *Service) setupUminekoNearestVoicedRoute(routeGroup fiber.Router) {
	routeGroup.Get("/nearest-voiced/:audioId", s.uminekoNearestVoiced)
}

func (s *Service) setupUminekoCharactersRoute(routeGroup fiber.Router) {
	routeGroup.Get("/characters", s.uminekoCharacters)
}

func (s *Service) setupUminekoStatsRoute(routeGroup fiber.Router) {
	routeGroup.Get("/stats", s.uminekoStats)
}

// uminekoSearch godoc
//
//	@Summary		Search quotes
//	@Description	Search for quotes by text query with optional filters
//	@Tags			umineko
//	@Produce		json
//	@Param			q			query		string	true	"Search query"
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(auto, en, wh, ja, ru, es, pt)
//	@Param			limit		query		int		false	"Maximum results"	default(30)
//	@Param			offset		query		int		false	"Offset for pagination"	default(0)
//	@Param			character	query		string	false	"Filter by character ID"	Enums(group_voices, kinzo, krauss, natsuhi, jessica, eva, hideyoshi, george, rudolf, kyrie, battler, ange, rosa, maria, genji, shannon, kanon, gohda, kumasawa, nanjo, amakusa, okonogi, kasumi, professor, kawabata, nanjo_son, kumasawa_son, beatrice, bernkastel, lambdadelta, virgilia, ronove, gaap, sakutarou, eva_beatrice, chiester_45, chiester_410, chiester_00, lucifer, leviathan, satan, belphegor, mammon, beelzebub, asmodeus, goat, erika, dlanor, gertrude, cornelia, featherine, zepar, furfur, lion, will, clair, ikuko, tohya, kinzo_young, bice, beato_elder, misc_voices, narrator)
//	@Param			interactionA	query		string	false	"Interaction filter: first character ID (requires interactionB)"	Enums(group_voices, kinzo, krauss, natsuhi, jessica, eva, hideyoshi, george, rudolf, kyrie, battler, ange, rosa, maria, genji, shannon, kanon, gohda, kumasawa, nanjo, amakusa, okonogi, kasumi, professor, kawabata, nanjo_son, kumasawa_son, beatrice, bernkastel, lambdadelta, virgilia, ronove, gaap, sakutarou, eva_beatrice, chiester_45, chiester_410, chiester_00, lucifer, leviathan, satan, belphegor, mammon, beelzebub, asmodeus, goat, erika, dlanor, gertrude, cornelia, featherine, zepar, furfur, lion, will, clair, ikuko, tohya, kinzo_young, bice, beato_elder, misc_voices, narrator)
//	@Param			interactionB	query		string	false	"Interaction filter: second character ID (requires interactionA)"	Enums(group_voices, kinzo, krauss, natsuhi, jessica, eva, hideyoshi, george, rudolf, kyrie, battler, ange, rosa, maria, genji, shannon, kanon, gohda, kumasawa, nanjo, amakusa, okonogi, kasumi, professor, kawabata, nanjo_son, kumasawa_son, beatrice, bernkastel, lambdadelta, virgilia, ronove, gaap, sakutarou, eva_beatrice, chiester_45, chiester_410, chiester_00, lucifer, leviathan, satan, belphegor, mammon, beelzebub, asmodeus, goat, erika, dlanor, gertrude, cornelia, featherine, zepar, furfur, lion, will, clair, ikuko, tohya, kinzo_young, bice, beato_elder, misc_voices, narrator)
//	@Param			episode		query		int		false	"Filter by episode (1-8)"
//	@Param			truth		query		string	false	"Filter by truth type"	Enums(red, blue, gold, purple)
//	@Success		200			{object}	dto.SearchAPIResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Router			/umineko/search [get]
func (s *Service) uminekoSearch(ctx fiber.Ctx) error {
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

	charID := character.Character(strings.TrimSpace(characterParam)).ID()
	intA := character.Character(strings.TrimSpace(interactionAParam)).ID()
	intB := character.Character(strings.TrimSpace(interactionBParam)).ID()

	searchParams := params.NewSearchParams(
		query,
		lang,
		limit,
		offset,
		charID,
		episode,
		intA,
		intB,
	)

	response := s.UminekoService.Search(searchParams, ctx.Query("truth"))
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

// uminekoRandom godoc
//
//	@Summary		Get random quote
//	@Description	Returns a random quote with optional filters
//	@Tags			umineko
//	@Produce		json
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(en, wh, ja, ru, es, pt)
//	@Param			character	query		string	false	"Filter by character ID"	Enums(group_voices, kinzo, krauss, natsuhi, jessica, eva, hideyoshi, george, rudolf, kyrie, battler, ange, rosa, maria, genji, shannon, kanon, gohda, kumasawa, nanjo, amakusa, okonogi, kasumi, professor, kawabata, nanjo_son, kumasawa_son, beatrice, bernkastel, lambdadelta, virgilia, ronove, gaap, sakutarou, eva_beatrice, chiester_45, chiester_410, chiester_00, lucifer, leviathan, satan, belphegor, mammon, beelzebub, asmodeus, goat, erika, dlanor, gertrude, cornelia, featherine, zepar, furfur, lion, will, clair, ikuko, tohya, kinzo_young, bice, beato_elder, misc_voices, narrator)
//	@Param			episode		query		int		false	"Filter by episode (1-8)"
//	@Param			truth		query		string	false	"Filter by truth type"	Enums(red, blue, gold, purple)
//	@Success		200			{object}	dto.UminekoQuote
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/umineko/random [get]
func (s *Service) uminekoRandom(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	characterID := character.Character(ctx.Query("character")).ID()
	episode := fiber.Query[int](ctx, "episode", 0)
	truth := umineko.TruthAll.Parse(ctx.Query("truth"))
	q := s.UminekoService.Random(lang, characterID, episode, truth)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no quotes available",
		})
	}
	return ctx.JSON(q)
}

// uminekoBrowse godoc
//
//	@Summary		Browse quotes
//	@Description	Browse all quotes with optional filters and pagination
//	@Tags			umineko
//	@Produce		json
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(en, wh, ja, ru, es, pt)
//	@Param			character	query		string	false	"Filter by character ID"	Enums(group_voices, kinzo, krauss, natsuhi, jessica, eva, hideyoshi, george, rudolf, kyrie, battler, ange, rosa, maria, genji, shannon, kanon, gohda, kumasawa, nanjo, amakusa, okonogi, kasumi, professor, kawabata, nanjo_son, kumasawa_son, beatrice, bernkastel, lambdadelta, virgilia, ronove, gaap, sakutarou, eva_beatrice, chiester_45, chiester_410, chiester_00, lucifer, leviathan, satan, belphegor, mammon, beelzebub, asmodeus, goat, erika, dlanor, gertrude, cornelia, featherine, zepar, furfur, lion, will, clair, ikuko, tohya, kinzo_young, bice, beato_elder, misc_voices, narrator)
//	@Param			interactionA	query		string	false	"Interaction filter: first character ID (requires interactionB)"	Enums(group_voices, kinzo, krauss, natsuhi, jessica, eva, hideyoshi, george, rudolf, kyrie, battler, ange, rosa, maria, genji, shannon, kanon, gohda, kumasawa, nanjo, amakusa, okonogi, kasumi, professor, kawabata, nanjo_son, kumasawa_son, beatrice, bernkastel, lambdadelta, virgilia, ronove, gaap, sakutarou, eva_beatrice, chiester_45, chiester_410, chiester_00, lucifer, leviathan, satan, belphegor, mammon, beelzebub, asmodeus, goat, erika, dlanor, gertrude, cornelia, featherine, zepar, furfur, lion, will, clair, ikuko, tohya, kinzo_young, bice, beato_elder, misc_voices, narrator)
//	@Param			interactionB	query		string	false	"Interaction filter: second character ID (requires interactionA)"	Enums(group_voices, kinzo, krauss, natsuhi, jessica, eva, hideyoshi, george, rudolf, kyrie, battler, ange, rosa, maria, genji, shannon, kanon, gohda, kumasawa, nanjo, amakusa, okonogi, kasumi, professor, kawabata, nanjo_son, kumasawa_son, beatrice, bernkastel, lambdadelta, virgilia, ronove, gaap, sakutarou, eva_beatrice, chiester_45, chiester_410, chiester_00, lucifer, leviathan, satan, belphegor, mammon, beelzebub, asmodeus, goat, erika, dlanor, gertrude, cornelia, featherine, zepar, furfur, lion, will, clair, ikuko, tohya, kinzo_young, bice, beato_elder, misc_voices, narrator)
//	@Param			limit		query		int		false	"Maximum results"	default(50)
//	@Param			offset		query		int		false	"Offset for pagination"	default(0)
//	@Param			episode		query		int		false	"Filter by episode (1-8)"
//	@Param			truth		query		string	false	"Filter by truth type"	Enums(red, blue, gold, purple)
//	@Success		200			{object}	dto.CharacterResponse
//	@Router			/umineko/browse [get]
func (s *Service) uminekoBrowse(ctx fiber.Ctx) error {
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

	charID := character.Character(strings.TrimSpace(characterParam)).ID()
	intA := character.Character(strings.TrimSpace(interactionAParam)).ID()
	intB := character.Character(strings.TrimSpace(interactionBParam)).ID()

	browseParams := params.NewBrowseParams(
		lang,
		limit,
		offset,
		charID,
		episode,
		intA,
		intB,
	)

	response := s.UminekoService.Browse(browseParams, ctx.Query("truth"))
	return ctx.JSON(response)
}

// uminekoByAudioID godoc
//
//	@Summary		Get quote by audio ID
//	@Description	Returns a specific quote identified by its audio ID
//	@Tags			umineko
//	@Produce		json
//	@Param			audioId		path		string	true	"Audio ID of the quote"
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(en, wh, ja, ru, es, pt)
//	@Success		200			{object}	dto.UminekoQuote
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/umineko/quote/{audioId} [get]
func (s *Service) uminekoByAudioID(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	audioID := ctx.Params("audioId")

	q := s.UminekoService.GetByAudioID(lang, audioID)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(q)
}

// uminekoByIndex godoc
//
//	@Summary		Get quote by index
//	@Description	Returns a specific quote identified by its position index in the parsed script
//	@Tags			umineko
//	@Produce		json
//	@Param			index		path		int		true	"Index of the quote in the parsed script"
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(en, wh, ja, ru, es, pt)
//	@Success		200			{object}	dto.UminekoQuote
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/umineko/quote/index/{index} [get]
func (s *Service) uminekoByIndex(ctx fiber.Ctx) error {
	lang := language.English.Parse(ctx.Query("lang"))
	index := fiber.Params[int](ctx, "index")

	q := s.UminekoService.GetByIndex(lang, index)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(q)
}

// uminekoContext godoc
//
//	@Summary		Get quote context
//	@Description	Returns surrounding dialogue lines for a specific quote
//	@Tags			umineko
//	@Produce		json
//	@Param			audioId		path		string	true	"Audio ID of the quote"
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(en, wh, ja, ru, es, pt)
//	@Param			lines		query		int		false	"Number of context lines before and after"	default(5)
//	@Success		200			{object}	dto.ContextResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/umineko/context/{audioId} [get]
func (s *Service) uminekoContext(ctx fiber.Ctx) error {
	return handleContext(ctx, s.UminekoService)
}

// uminekoNearestVoiced godoc
//
//	@Summary		Find nearest voiced quote
//	@Description	Returns the audio ID of the nearest quote with voice audio in a given direction
//	@Tags			umineko
//	@Produce		json
//	@Param			audioId		path		string	true	"Audio ID of the reference quote"
//	@Param			lang		query		string	false	"Language"	default(en)	Enums(en, wh, ja, ru, es, pt)
//	@Param			direction	query		string	false	"Direction to search"	default(next)	Enums(next, prev)
//	@Success		200			{object}	map[string]string
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		404			{object}	dto.ErrorResponse
//	@Router			/umineko/nearest-voiced/{audioId} [get]
func (s *Service) uminekoNearestVoiced(ctx fiber.Ctx) error {
	return handleNearestVoiced(ctx, s.UminekoService)
}

// uminekoCharacters godoc
//
//	@Summary		List all characters
//	@Description	Returns a map of character IDs to character names
//	@Tags			umineko
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/umineko/characters [get]
func (s *Service) uminekoCharacters(ctx fiber.Ctx) error {
	ctx.Set("Cache-Control", "public, max-age=86400")
	chars := s.UminekoService.GetCharacters()
	characters := make(map[string]string, len(chars))
	for k, v := range chars {
		characters[string(k)] = v
	}
	return ctx.JSON(dto.CharactersResult{
		Characters: characters,
		Additional: map[string]string{},
	})
}

// uminekoStats godoc
//
//	@Summary		Get quote statistics
//	@Description	Returns statistics about quotes including top speakers, truth per episode, and character interactions
//	@Tags			umineko
//	@Produce		json
//	@Param			episode		query		int		false	"Filter by episode (1-8), 0 for all"
//	@Success		200			{object}	dto.StatsResult
//	@Router			/umineko/stats [get]
func (s *Service) uminekoStats(ctx fiber.Ctx) error {
	episode := fiber.Query[int](ctx, "episode", 0)
	return ctx.JSON(s.UminekoService.GetStats().Compute(episode))
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
