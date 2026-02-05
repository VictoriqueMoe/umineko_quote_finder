package controllers

import (
	"fmt"
	"strings"

	"umineko_quote/internal/og"

	"github.com/gofiber/fiber/v2"
)

func (s *Service) getAllOGAPIRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupOGBuilderImageRoute,
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

func (s *Service) setupOGBuilderImageRoute(routeGroup fiber.Router) {
	routeGroup.Get("/og/builder.png", s.ogBuilderImage)
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

	data, err := s.OGImageGenerator.Generate(audioId, lang, q.Text, q.TextHtml, q.Character, q.Episode, q.ContentType)
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

func (s *Service) baseURL(ctx *fiber.Ctx) string {
	scheme := "https"
	if strings.HasPrefix(ctx.Hostname(), "localhost") || strings.HasPrefix(ctx.Hostname(), "127.0.0.1") {
		scheme = "http"
	}
	proto := ctx.Get("X-Forwarded-Proto")
	if proto != "" {
		scheme = proto
	}
	return fmt.Sprintf("%s://%s", scheme, ctx.Hostname())
}

type builderSegmentMeta struct {
	CharID    string
	AudioID   string
	Character string
	Text      string
}

func (s *Service) parseBuilderSegments(param, lang string) []builderSegmentMeta {
	parts := strings.Split(param, ",")
	if len(parts) > 20 {
		parts = parts[:20]
	}

	var segments []builderSegmentMeta
	for _, part := range parts {
		part = strings.TrimSpace(part)
		colonIdx := strings.IndexByte(part, ':')
		if colonIdx < 1 || colonIdx >= len(part)-1 {
			continue
		}
		charId := part[:colonIdx]
		audioId := part[colonIdx+1:]
		if !audioIdPattern.MatchString(charId) || !audioIdPattern.MatchString(audioId) {
			continue
		}

		q := s.QuoteService.GetByAudioID(lang, audioId)
		if q != nil {
			clipText := q.Text
			if q.AudioTextMap != nil {
				if mapped, ok := q.AudioTextMap[audioId]; ok {
					clipText = mapped
				}
			}
			segments = append(segments, builderSegmentMeta{
				CharID:    charId,
				AudioID:   audioId,
				Character: q.Character,
				Text:      clipText,
			})
		} else {
			segments = append(segments, builderSegmentMeta{
				CharID:  charId,
				AudioID: audioId,
			})
		}
	}
	return segments
}

func (s *Service) ogPage(ctx *fiber.Ctx) error {
	audioId := ctx.Query("quote")
	builderParam := ctx.Query("builder")

	if audioId == "" && builderParam == "" {
		html := s.replaceOGPlaceholders(defaultOGTitle, defaultOGDescription, defaultTwitterDesc, defaultOGImage)
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.SendString(html)
	}

	lang := ctx.Query("lang", "en")
	base := s.baseURL(ctx)

	// Handle builder links
	if builderParam != "" {
		segments := s.parseBuilderSegments(builderParam, lang)
		if len(segments) == 0 {
			html := s.replaceOGPlaceholders(defaultOGTitle, defaultOGDescription, defaultTwitterDesc, defaultOGImage)
			ctx.Set("Content-Type", "text/html; charset=utf-8")
			return ctx.SendString(html)
		}

		// Build title from unique character names
		seen := map[string]bool{}
		var names []string
		for _, seg := range segments {
			if seg.Character != "" && !seen[seg.Character] {
				seen[seg.Character] = true
				names = append(names, seg.Character)
			}
		}
		title := "Voice Build"
		if len(names) > 0 {
			title = fmt.Sprintf("Voice Build \u2014 %s", strings.Join(names, ", "))
		}

		// Build description from quote texts
		var descParts []string
		for _, seg := range segments {
			if seg.Character != "" && seg.Text != "" {
				descParts = append(descParts, fmt.Sprintf("%s: \u201C%s\u201D", seg.Character, seg.Text))
			}
		}
		description := strings.Join(descParts, " \u2192 ")
		if len(description) > 200 {
			description = description[:197] + "..."
		}
		if description == "" {
			description = fmt.Sprintf("A voice build with %d clips from Umineko no Naku Koro ni.", len(segments))
		}

		imageURL := fmt.Sprintf("%s/api/v1/og/builder.png?segments=%s&lang=%s", base, builderParam, lang)

		html := s.replaceOGPlaceholders(title, description, description, imageURL)
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.SendString(html)
	}

	// Handle single quote links
	q := s.QuoteService.GetByAudioID(lang, audioId)
	if q == nil {
		html := s.replaceOGPlaceholders(defaultOGTitle, defaultOGDescription, defaultTwitterDesc, defaultOGImage)
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.SendString(html)
	}

	title := fmt.Sprintf("%s \u2014 Umineko Quote", q.Character)
	description := q.Text
	if len(description) > 200 {
		description = description[:197] + "..."
	}
	imageURL := fmt.Sprintf("%s/api/v1/og/%s.png?lang=%s", base, audioId, lang)

	html := s.replaceOGPlaceholders(title, description, description, imageURL)
	ctx.Set("Content-Type", "text/html; charset=utf-8")
	return ctx.SendString(html)
}

func (s *Service) ogBuilderImage(ctx *fiber.Ctx) error {
	segmentsParam := ctx.Query("segments")
	if segmentsParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query parameter 'segments' is required",
		})
	}

	lang := ctx.Query("lang", "en")
	segments := s.parseBuilderSegments(segmentsParam, lang)
	if len(segments) == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no valid segments found",
		})
	}

	var lines []og.DialogueLine
	for _, seg := range segments {
		if seg.Character != "" && seg.Text != "" {
			lines = append(lines, og.DialogueLine{Character: seg.Character, Text: seg.Text})
		}
	}

	data, err := s.OGImageGenerator.GenerateBuilder(segmentsParam, lang, lines)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate image",
		})
	}

	ctx.Set("Content-Type", "image/png")
	ctx.Set("Cache-Control", "public, max-age=86400")
	return ctx.Send(data)
}

func escapeAttr(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
