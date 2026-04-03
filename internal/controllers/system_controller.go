package controllers

import (
	"umineko_quote/internal/quote/language"

	"github.com/gofiber/fiber/v3"
)

var expectedLanguages = []language.Language{
	language.English,
	language.WitchHunt,
	language.Japanese,
	language.Russian,
	language.Spanish,
	language.Portuguese,
}

func (s *Service) getAllSystemRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupHealthRoute,
		s.setupConfigRoute,
	}
}

func (s *Service) setupHealthRoute(routeGroup fiber.Router) {
	routeGroup.Get("/health", s.healthCheck)
}

func (s *Service) setupConfigRoute(routeGroup fiber.Router) {
	routeGroup.Get("/config", s.config)
}

func (s *Service) healthCheck(ctx fiber.Ctx) error {
	loaded := s.UminekoService.LoadedLanguages()

	languages := make(fiber.Map, len(expectedLanguages))
	healthy := true
	for _, lang := range expectedLanguages {
		count, ok := loaded[lang]
		if !ok || count == 0 {
			healthy = false
		}
		languages[string(lang)] = count
	}

	status := "ok"
	httpStatus := fiber.StatusOK
	if !healthy {
		status = "degraded"
		httpStatus = fiber.StatusServiceUnavailable
	}

	higurashiLoaded := s.HigurashiService.LoadedLanguages()

	return ctx.Status(httpStatus).JSON(fiber.Map{
		"status":    status,
		"service":   "umineko-quote-service",
		"languages": languages,
		"hasAudio":  s.UminekoService.HasAudio(),
		"higurashi": fiber.Map{
			"languages": higurashiLoaded,
		},
	})
}

func (s *Service) config(ctx fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"hasAudio": s.UminekoService.HasAudio(),
	})
}
