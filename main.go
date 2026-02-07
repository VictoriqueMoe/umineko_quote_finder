package main

import (
	"embed"
	"io/fs"
	"log"
	"umineko_quote/internal/audio"
	"umineko_quote/internal/controllers"
	"umineko_quote/internal/og"
	"umineko_quote/internal/quote"
	"umineko_quote/internal/routes"
	"umineko_quote/internal/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	app := fiber.New()

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path} ${queryParams}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	quoteService := quote.NewService()
	ogGen := og.NewImageGenerator()
	audioCombiner, err := audio.NewCombiner()
	if err != nil {
		log.Fatalf("failed to initialize audio combiner: %v", err)
	}
	htmlBytes, _ := staticFiles.ReadFile("static/index.html")
	service := controllers.NewService(quoteService, ogGen, audioCombiner, string(htmlBytes))
	routes.PublicRoutes(service, app)

	staticFS, _ := fs.Sub(staticFiles, "static")
	app.Get("/*", static.New("", static.Config{
		FS: staticFS,
	}))

	utils.StartServerWithGracefulShutdown(app, ":3000")
}
