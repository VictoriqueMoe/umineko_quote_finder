package main

import (
	"embed"
	"io/fs"
	"log"
	"umineko_quote/internal/audio"
	"umineko_quote/internal/config"
	"umineko_quote/internal/controllers"
	"umineko_quote/internal/og"
	"umineko_quote/internal/quote/ciconia"
	"umineko_quote/internal/quote/game"
	"umineko_quote/internal/quote/higurashi"
	"umineko_quote/internal/quote/umineko"
	"umineko_quote/internal/repository"
	"umineko_quote/internal/routes"
	"umineko_quote/internal/utils"

	_ "umineko_quote/docs"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/gofiber/contrib/v3/swaggo"
)

//go:embed static/*
var staticFiles embed.FS

// @title			When They Cry Quote API
// @version		1.0
// @description	API for searching and browsing quotes from Umineko no Naku Koro ni, Higurashi no Naku Koro ni, and Ciconia no Naku Koro ni
// @contact.name	Featherine Augustus Aurora
// @contact.url	https://x.com/FeatherineFAA
// @contact.email	FAA@auaurora.moe
// @BasePath		/api/v1
// @schemes		https http
func main() {
	app := fiber.New()

	app.Use("/api", cors.New())

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} ${path} ${queryParams}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	if err := config.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	repos, err := repository.New()
	if err != nil {
		log.Fatalf("failed to initialise repositories: %v", err)
	}
	defer func() {
		if err := repos.Close(); err != nil {
			log.Printf("failed to close repositories: %v", err)
		}
	}()

	audioResolver, err := audio.NewAudioResolver(game.UminekoAudioDir)
	if err != nil {
		log.Fatalf("failed to scan umineko audio directory: %v", err)
	}

	uminekoService, err := umineko.NewService(repos.Umineko, audioResolver)
	if err != nil {
		log.Fatalf("failed to initialise umineko service: %v", err)
	}

	higurashiService, err := higurashi.NewService(repos.Higurashi)
	if err != nil {
		log.Fatalf("failed to initialise higurashi service: %v", err)
	}

	ciconiaService, err := ciconia.NewService(repos.Ciconia)
	if err != nil {
		log.Fatalf("failed to initialise ciconia service: %v", err)
	}

	ogGen := og.NewImageGenerator()
	audioCombiner, err := audio.NewCombiner()
	if err != nil {
		log.Fatalf("failed to initialize audio combiner: %v", err)
	}
	htmlBytes, _ := staticFiles.ReadFile("static/index.html")
	service := controllers.NewService(uminekoService, higurashiService, ciconiaService, ogGen, audioCombiner, string(htmlBytes))
	routes.PublicRoutes(service, app)

	app.Get("/swagger/*", swaggo.HandlerDefault)

	staticFS, _ := fs.Sub(staticFiles, "static")
	app.Get("/*", static.New("", static.Config{
		FS: staticFS,
	}))

	utils.StartServerWithGracefulShutdown(app, ":3000")
}
