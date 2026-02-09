package controllers

import "github.com/gofiber/fiber/v3"

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
//	@Description	Returns the health status of the service
//	@Tags			system
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (s *Service) healthCheck(ctx fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status":  "ok",
		"service": "umineko-quote-service",
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
