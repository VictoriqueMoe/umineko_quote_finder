package controllers

import "github.com/gofiber/fiber/v2"

func (s *Service) getAllSystemRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupHealthRoute,
		s.setupConfigRoute,
	}
}

func (s *Service) setupHealthRoute(routeGroup fiber.Router) {
	routeGroup.Get("/health", s.healthCheck)
}

func (s *Service) healthCheck(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status":  "ok",
		"service": "umineko-quote-service",
	})
}

func (s *Service) setupConfigRoute(routeGroup fiber.Router) {
	routeGroup.Get("/config", s.config)
}

func (s *Service) config(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"hasAudio": s.QuoteService.HasAudio(),
	})
}
