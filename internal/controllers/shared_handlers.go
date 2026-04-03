package controllers

import (
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"

	"github.com/gofiber/fiber/v3"
)

type QuoteLookup interface {
	GetContext(lang language.Language, audioID string, lines int) *dto.ContextResponse
	NearestVoicedAudioID(lang language.Language, audioID string, direction string) string
}

func handleContext(ctx fiber.Ctx, lookup QuoteLookup) error {
	lang := language.English.Parse(ctx.Query("lang"))
	audioID := ctx.Params("audioId")
	if !audioIdPattern.MatchString(audioID) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid audio ID",
		})
	}

	lines := fiber.Query[int](ctx, "lines", 5)
	result := lookup.GetContext(lang, audioID, lines)
	if result == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(result)
}

func handleNearestVoiced(ctx fiber.Ctx, lookup QuoteLookup) error {
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

	result := lookup.NearestVoicedAudioID(lang, audioID, direction)
	if result == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no voiced quote found",
		})
	}
	return ctx.JSON(fiber.Map{"audioId": result})
}
