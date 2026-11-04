package utils

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func InternalError(ctx fiber.Ctx, msg string, cause ...error) error {
	for _, err := range cause {
		if err == nil {
			continue
		}
		log.Printf("[%s %s] %s: %v", ctx.Method(), ctx.Path(), msg, err)
	}
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": msg})
}
