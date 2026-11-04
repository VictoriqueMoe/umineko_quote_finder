package controllers

import (
	"umineko_quote/internal/controllers/utils"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/quotes"

	"github.com/gofiber/fiber/v3"
)

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
