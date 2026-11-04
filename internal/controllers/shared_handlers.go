package controllers

import (
	"umineko_quote/internal/controllers/utils"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"
	"umineko_quote/internal/quote/quotes"

	"github.com/gofiber/fiber/v3"
)

func handleSearch[Q any](ctx fiber.Ctx, svc quotes.Service[Q], characterID func(string) string, episode int, scope params.QuoteScope) error {
	query := ctx.Query("q")
	if query == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query parameter 'q' is required",
		})
	}

	characterParam := ctx.Query("character")
	interactionAParam := ctx.Query("interactionA")
	interactionBParam := ctx.Query("interactionB")

	if errMsg, invalid := validateInteractionQueryParams(characterParam, interactionAParam, interactionBParam); invalid {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errMsg,
		})
	}

	searchParams := params.NewSearchParams(
		query,
		language.English.Parse(ctx.Query("lang")),
		fiber.Query[int](ctx, "limit"),
		fiber.Query[int](ctx, "offset"),
		characterID(characterParam),
		episode,
		characterID(interactionAParam),
		characterID(interactionBParam),
		fiber.Query[bool](ctx, "exact"),
	)

	response, err := svc.Search(searchParams, scope)
	if err != nil {
		return utils.InternalError(ctx, "failed to search quotes", err)
	}

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

func handleBrowse[Q any](ctx fiber.Ctx, svc quotes.Service[Q], characterID func(string) string, episode int, scope params.QuoteScope) error {
	characterParam := ctx.Query("character")
	interactionAParam := ctx.Query("interactionA")
	interactionBParam := ctx.Query("interactionB")

	if errMsg, invalid := validateInteractionQueryParams(characterParam, interactionAParam, interactionBParam); invalid {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errMsg,
		})
	}

	browseParams := params.NewBrowseParams(
		language.English.Parse(ctx.Query("lang")),
		fiber.Query[int](ctx, "limit"),
		fiber.Query[int](ctx, "offset"),
		characterID(characterParam),
		episode,
		characterID(interactionAParam),
		characterID(interactionBParam),
	)

	response, err := svc.Browse(browseParams, scope)
	if err != nil {
		return utils.InternalError(ctx, "failed to browse quotes", err)
	}
	return ctx.JSON(response)
}

func handleRandom[Q any](ctx fiber.Ctx, svc quotes.Service[Q], characterID func(string) string, episode int, scope params.QuoteScope) error {
	lang := language.English.Parse(ctx.Query("lang"))

	q, err := svc.Random(lang, characterID(ctx.Query("character")), episode, scope)
	if err != nil {
		return utils.InternalError(ctx, "failed to load random quote", err)
	}
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no quotes available",
		})
	}
	return ctx.JSON(q)
}

func handleByAudioID[Q any](ctx fiber.Ctx, svc quotes.Service[Q], audioID string) error {
	lang := language.English.Parse(ctx.Query("lang"))

	q, err := svc.GetByAudioID(lang, audioID)
	if err != nil {
		return utils.InternalError(ctx, "failed to load quote", err)
	}
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(q)
}

func handleByIndex[Q any](ctx fiber.Ctx, svc quotes.Service[Q]) error {
	lang := language.English.Parse(ctx.Query("lang"))
	index := fiber.Params[int](ctx, "index")

	q, err := svc.GetByIndex(lang, index)
	if err != nil {
		return utils.InternalError(ctx, "failed to load quote", err)
	}
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(q)
}

func handleContext[Q any](ctx fiber.Ctx, svc quotes.Service[Q], audioID string) error {
	lang := language.English.Parse(ctx.Query("lang"))
	lines := fiber.Query[int](ctx, "lines")

	result, err := svc.GetContext(lang, audioID, lines)
	if err != nil {
		return utils.InternalError(ctx, "failed to load quote context", err)
	}
	if result == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(result)
}

func handleNearestVoiced[Q any](ctx fiber.Ctx, svc quotes.Service[Q], audioID string) error {
	lang := language.English.Parse(ctx.Query("lang"))

	direction := ctx.Query("direction", "next")
	if direction != "next" && direction != "prev" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "direction must be 'next' or 'prev'",
		})
	}

	result, err := svc.NearestVoicedAudioID(lang, audioID, direction)
	if err != nil {
		return utils.InternalError(ctx, "failed to find nearest voiced quote", err)
	}
	if result == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no voiced quote found",
		})
	}
	return ctx.JSON(fiber.Map{"audioId": result})
}
