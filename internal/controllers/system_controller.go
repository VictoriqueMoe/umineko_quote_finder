package controllers

import (
	"umineko_quote/internal/quote/language"

	"github.com/gofiber/fiber/v3"
)

var expectedLanguages = []language.Language{
	language.English,
	language.Japanese,
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

// healthCheck godoc
//
//	@Summary		Health check
//	@Description	Returns the health status of the service including language loading status
//	@Tags			system
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		503	{object}	map[string]interface{}
//	@Router			/health [get]
func (s *Service) healthCheck(ctx fiber.Ctx) error {
	loaded := s.QuoteService.LoadedLanguages()

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

	return ctx.Status(httpStatus).JSON(fiber.Map{
		"status":    status,
		"service":   "umineko-quote-service",
		"languages": languages,
		"hasAudio":  s.QuoteService.HasAudio(),
	})
}

func (s *Service) setupConfigRoute(routeGroup fiber.Router) {
	routeGroup.Get("/config", s.config)
}

// config godoc
//
//	@Summary		Get configuration
//	@Description	Returns service configuration flags
//	@Tags			system
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/config [get]
func (s *Service) config(ctx fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"hasAudio": s.QuoteService.HasAudio(),
	})
}
