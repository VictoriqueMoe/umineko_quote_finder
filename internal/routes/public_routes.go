package routes

import (
	"umineko_quote/internal/controllers"

	"github.com/gofiber/fiber/v3"
)

func PublicRoutes(service controllers.Service, app *fiber.App) {
	api := app.Group("/api/v1")

	umiGroup := api.Group("/umineko")
	umiRoutes := service.GetUminekoRoutes()
	for i := 0; i < len(umiRoutes); i++ {
		umiRoutes[i](umiGroup)
	}

	higuGroup := api.Group("/higurashi")
	higuRoutes := service.GetHigurashiRoutes()
	for i := 0; i < len(higuRoutes); i++ {
		higuRoutes[i](higuGroup)
	}

	systemRoutes := service.GetSystemRoutes()
	for i := 0; i < len(systemRoutes); i++ {
		systemRoutes[i](api)
	}

	pageRoutes := service.GetPageRoutes()
	for i := 0; i < len(pageRoutes); i++ {
		pageRoutes[i](app)
	}
}
