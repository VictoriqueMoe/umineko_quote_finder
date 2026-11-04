package routes

import (
	"umineko_quote/internal/controllers"

	"github.com/gofiber/fiber/v3"
)

func PublicRoutes(service controllers.Service, app *fiber.App) {
	api := app.Group("/api/v1")

	umiGroup := api.Group("/umineko")
	for _, route := range service.GetUminekoRoutes() {
		route(umiGroup)
	}

	higuGroup := api.Group("/higurashi")
	for _, route := range service.GetHigurashiRoutes() {
		route(higuGroup)
	}

	ciconiaGroup := api.Group("/ciconia")
	for _, route := range service.GetCiconiaRoutes() {
		route(ciconiaGroup)
	}

	for _, route := range service.GetSystemRoutes() {
		route(api)
	}

	for _, route := range service.GetPageRoutes() {
		route(app)
	}
}
