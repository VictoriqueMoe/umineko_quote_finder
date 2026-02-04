package controllers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (s *Service) getAllOGAPIRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupOGImageRoute,
	}
}

func (s *Service) getAllOGPageRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupOGPageRoute,
	}
}

func (s *Service) setupOGImageRoute(routeGroup fiber.Router) {
	routeGroup.Get("/og/:audioId.png", s.ogImage)
}

func (s *Service) ogImage(ctx *fiber.Ctx) error {
	audioId := ctx.Params("audioId")
	if !audioIdPattern.MatchString(audioId) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid audio ID",
		})
	}

	lang := ctx.Query("lang", "en")

	q := s.QuoteService.GetByAudioID(lang, audioId)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}

	data, err := s.OGImageGenerator.Generate(audioId, lang, q.Text, q.Character, q.Episode, q.ContentType)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate image",
		})
	}

	ctx.Set("Content-Type", "image/png")
	ctx.Set("Cache-Control", "public, max-age=86400")
	return ctx.Send(data)
}

func (s *Service) setupOGPageRoute(routeGroup fiber.Router) {
	routeGroup.Get("/", s.ogPage)
}

const (
	defaultOGTitle       = "Umineko Quote Search"
	defaultOGDescription = "Search through the words of witches, humans, and furniture from Umineko no Naku Koro ni. When the seagulls cry, none shall remain."
	defaultOGImage       = "https://waifuvault.moe/f/5e9cf90a-8a63-48b3-802d-1bc9be9062ea/clipboard-image-1769601762638.png"
	defaultTwitterDesc   = "Search through the words of witches, humans, and furniture from Umineko no Naku Koro ni."
)

func replaceMetaContent(html, attrName, attrValue, oldContent, newContent string) string {
	old := attrName + `="` + attrValue + `" content="` + oldContent + `"`
	repl := attrName + `="` + attrValue + `" content="` + newContent + `"`
	return strings.Replace(html, old, repl, 1)
}

func (s *Service) replaceOGPlaceholders(title, description, twitterDesc, imageURL string) string {
	html := s.HTMLContent
	html = replaceMetaContent(html, "property", "og:title", defaultOGTitle, escapeAttr(title))
	html = replaceMetaContent(html, "property", "og:description", defaultOGDescription, escapeAttr(description))
	html = replaceMetaContent(html, "property", "og:image", defaultOGImage, imageURL)
	html = replaceMetaContent(html, "name", "twitter:title", defaultOGTitle, escapeAttr(title))
	html = replaceMetaContent(html, "name", "twitter:description", defaultTwitterDesc, escapeAttr(twitterDesc))
	html = replaceMetaContent(html, "name", "twitter:image", defaultOGImage, imageURL)
	return html
}

func (s *Service) ogPage(ctx *fiber.Ctx) error {
	audioId := ctx.Query("quote")
	if audioId == "" {
		html := s.replaceOGPlaceholders(defaultOGTitle, defaultOGDescription, defaultTwitterDesc, defaultOGImage)
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.SendString(html)
	}

	lang := ctx.Query("lang", "en")

	q := s.QuoteService.GetByAudioID(lang, audioId)
	if q == nil {
		html := s.replaceOGPlaceholders(defaultOGTitle, defaultOGDescription, defaultTwitterDesc, defaultOGImage)
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.SendString(html)
	}

	scheme := "https"
	if strings.HasPrefix(ctx.Hostname(), "localhost") || strings.HasPrefix(ctx.Hostname(), "127.0.0.1") {
		scheme = "http"
	}
	proto := ctx.Get("X-Forwarded-Proto")
	if proto != "" {
		scheme = proto
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, ctx.Hostname())

	title := fmt.Sprintf("%s \u2014 Umineko Quote", q.Character)
	description := q.Text
	if len(description) > 200 {
		description = description[:197] + "..."
	}
	imageURL := fmt.Sprintf("%s/api/v1/og/%s.png?lang=%s", baseURL, audioId, lang)

	html := s.replaceOGPlaceholders(title, description, description, imageURL)
	ctx.Set("Content-Type", "text/html; charset=utf-8")
	return ctx.SendString(html)
}

func escapeAttr(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
