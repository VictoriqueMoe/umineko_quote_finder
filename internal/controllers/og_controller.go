package controllers

import (
	"fmt"
	"strings"

	"umineko_quote/internal/og"
	"umineko_quote/internal/quote/language"

	"github.com/gofiber/fiber/v3"
)

var (
	uminekoOGDefaults = ogDefaults{
		title:       "Umineko Quote Search",
		description: "Search through the words of witches, humans, and furniture from Umineko no Naku Koro ni. When the seagulls cry, none shall remain.",
		twitterDesc: "Search through the words of witches, humans, and furniture from Umineko no Naku Koro ni.",
		image:       "https://waifuvault.moe/f/5e9cf90a-8a63-48b3-802d-1bc9be9062ea/clipboard-image-1769601762638.png",
		brand:       "Umineko Quote Search",
		quoteSuffix: "Umineko Quote",
	}

	higurashiOGDefaults = ogDefaults{
		title:       "Higurashi Quote Search",
		description: "Search through the words of Hinamizawa's residents from Higurashi no Naku Koro ni. When the cicadas cry, none shall escape.",
		twitterDesc: "Search through the words of Hinamizawa's residents from Higurashi no Naku Koro ni.",
		image:       "https://waifuvault.moe/f/5e9cf90a-8a63-48b3-802d-1bc9be9062ea/clipboard-image-1769601762638.png",
		brand:       "Higurashi Quote Search",
		quoteSuffix: "Higurashi Quote",
	}

	uminekoEpisodeNames = map[int]string{
		1: "Episode 1 \u2014 Legend",
		2: "Episode 2 \u2014 Turn",
		3: "Episode 3 \u2014 Banquet",
		4: "Episode 4 \u2014 Alliance",
		5: "Episode 5 \u2014 End",
		6: "Episode 6 \u2014 Dawn",
		7: "Episode 7 \u2014 Requiem",
		8: "Episode 8 \u2014 Twilight",
	}

	higurashiArcNames = map[string]string{
		"onikakushi":        "Onikakushi",
		"watanagashi":       "Watanagashi",
		"tatarigoroshi":     "Tatarigoroshi",
		"himatsubushi":      "Himatsubushi",
		"meakashi":          "Meakashi",
		"tsumihoroboshi":    "Tsumihoroboshi",
		"minagoroshi":       "Minagoroshi",
		"matsuribayashi":    "Matsuribayashi",
		"someutsushi":       "Someutsushi",
		"kageboshi":         "Kageboshi",
		"tsukiotoshi":       "Tsukiotoshi",
		"taraimawashi":      "Taraimawashi",
		"yoigoshi":          "Yoigoshi",
		"tokihogushi":       "Tokihogushi",
		"miotsukushi_omote": "Miotsukushi Omote",
		"kakera":            "Kakera",
		"miotsukushi_ura":   "Miotsukushi Ura",
		"kotohogushi":       "Kotohogushi",
		"hajisarashi":       "Hajisarashi",
	}
)

type ogDefaults struct {
	title       string
	description string
	twitterDesc string
	image       string
	brand       string
	quoteSuffix string
}

func (s *Service) getUminekoOGRoutes() []FSetupRoute {
	return []FSetupRoute{
		func(r fiber.Router) { r.Get("/og/builder.png", s.uminekoOGBuilderImage) },
		func(r fiber.Router) { r.Get("/og/:audioId.png", s.uminekoOGImage) },
	}
}

func (s *Service) getHigurashiOGRoutes() []FSetupRoute {
	return []FSetupRoute{
		func(r fiber.Router) { r.Get("/og/quote.png", s.higurashiOGImage) },
	}
}

func (s *Service) getAllOGPageRoutes() []FSetupRoute {
	return []FSetupRoute{
		func(r fiber.Router) { r.Get("/", s.ogPage) },
	}
}

func (s *Service) uminekoOGImage(ctx fiber.Ctx) error {
	audioId := ctx.Params("audioId")
	if !audioIdPattern.MatchString(audioId) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid audio ID"})
	}

	lang := language.English.Parse(ctx.Query("lang"))
	q := s.UminekoService.GetByAudioID(lang, audioId)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "quote not found"})
	}

	full := ctx.Query("full") == "true"
	data, err := s.OGImageGenerator.Generate(audioId, lang, q.Text, q.TextHtml, q.Character, full, uminekoOGDefaults.brand, uminekoEpisodeLabel(q.Episode, q.ContentType))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate image"})
	}

	ctx.Set("Content-Type", "image/png")
	ctx.Set("Cache-Control", "public, max-age=86400")
	return ctx.Send(data)
}

func (s *Service) higurashiOGImage(ctx fiber.Ctx) error {
	audioId := ctx.Query("audioId")
	if audioId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "audioId query parameter is required"})
	}

	lang := language.English.Parse(ctx.Query("lang"))
	q := s.HigurashiService.GetByAudioID(lang, audioId)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "quote not found"})
	}

	full := ctx.Query("full") == "true"
	data, err := s.OGImageGenerator.Generate(audioId, lang, q.Text, q.TextHtml, q.Character, full, higurashiOGDefaults.brand, higurashiArcLabel(q.Arc))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate image"})
	}

	ctx.Set("Content-Type", "image/png")
	ctx.Set("Cache-Control", "public, max-age=86400")
	return ctx.Send(data)
}

func (s *Service) uminekoOGBuilderImage(ctx fiber.Ctx) error {
	segmentsParam := ctx.Query("segments")
	if segmentsParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query parameter 'segments' is required"})
	}

	lang := language.English.Parse(ctx.Query("lang"))
	segments := s.parseUminekoBuilderSegments(segmentsParam, lang)
	if len(segments) == 0 {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "no valid segments found"})
	}

	var lines []og.DialogueLine
	for _, seg := range segments {
		if seg.Character != "" && seg.Text != "" {
			lines = append(lines, og.DialogueLine{Character: seg.Character, Text: seg.Text})
		}
	}

	data, err := s.OGImageGenerator.GenerateBuilder(segmentsParam, lang, lines, uminekoOGDefaults.brand)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to generate image"})
	}

	ctx.Set("Content-Type", "image/png")
	ctx.Set("Cache-Control", "public, max-age=86400")
	return ctx.Send(data)
}

func (s *Service) ogPage(ctx fiber.Ctx) error {
	game := ctx.Query("game", "umineko")
	audioId := ctx.Query("quote")
	builderParam := ctx.Query("builder")

	defaults := uminekoOGDefaults
	if game == "higurashi" {
		defaults = higurashiOGDefaults
	}

	if audioId == "" && builderParam == "" {
		html := s.replaceOGPlaceholders(defaults, defaults.title, defaults.description, defaults.twitterDesc, defaults.image)
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.SendString(html)
	}

	lang := language.English.Parse(ctx.Query("lang"))
	base := s.baseURL(ctx)

	if builderParam != "" && game != "higurashi" {
		segments := s.parseUminekoBuilderSegments(builderParam, lang)
		if len(segments) == 0 {
			html := s.replaceOGPlaceholders(defaults, defaults.title, defaults.description, defaults.twitterDesc, defaults.image)
			ctx.Set("Content-Type", "text/html; charset=utf-8")
			return ctx.SendString(html)
		}

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
			description = fmt.Sprintf("A voice build with %d clips from %s.", len(segments), defaults.title)
		}

		imageURL := fmt.Sprintf("%s/api/v1/umineko/og/builder.png?segments=%s&lang=%s", base, builderParam, lang)
		html := s.replaceOGPlaceholders(defaults, title, description, description, imageURL)
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.SendString(html)
	}

	if game == "higurashi" {
		q := s.HigurashiService.GetByAudioID(lang, audioId)
		if q == nil {
			html := s.replaceOGPlaceholders(defaults, defaults.title, defaults.description, defaults.twitterDesc, defaults.image)
			ctx.Set("Content-Type", "text/html; charset=utf-8")
			return ctx.SendString(html)
		}

		title := fmt.Sprintf("%s \u2014 %s", q.Character, defaults.quoteSuffix)
		description := q.Text
		if len(description) > 200 {
			description = description[:197] + "..."
		}
		imageURL := fmt.Sprintf("%s/api/v1/higurashi/og/quote.png?audioId=%s&lang=%s", base, audioId, lang)
		html := s.replaceOGPlaceholders(defaults, title, description, description, imageURL)
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.SendString(html)
	}

	q := s.UminekoService.GetByAudioID(lang, audioId)
	if q == nil {
		html := s.replaceOGPlaceholders(defaults, defaults.title, defaults.description, defaults.twitterDesc, defaults.image)
		ctx.Set("Content-Type", "text/html; charset=utf-8")
		return ctx.SendString(html)
	}

	title := fmt.Sprintf("%s \u2014 %s", q.Character, defaults.quoteSuffix)
	description := q.Text
	if len(description) > 200 {
		description = description[:197] + "..."
	}
	imageURL := fmt.Sprintf("%s/api/v1/umineko/og/%s.png?lang=%s", base, audioId, lang)
	html := s.replaceOGPlaceholders(defaults, title, description, description, imageURL)
	ctx.Set("Content-Type", "text/html; charset=utf-8")
	return ctx.SendString(html)
}

type builderSegmentMeta struct {
	CharID    string
	AudioID   string
	Character string
	Text      string
}

func (s *Service) parseUminekoBuilderSegments(param string, lang language.Language) []builderSegmentMeta {
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

		q := s.UminekoService.GetByAudioID(lang, audioId)
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

func (s *Service) replaceOGPlaceholders(defaults ogDefaults, title, description, twitterDesc, imageURL string) string {
	html := s.HTMLContent
	html = replaceMetaContent(html, "property", "og:title", defaults.title, escapeAttr(title))
	html = replaceMetaContent(html, "property", "og:description", defaults.description, escapeAttr(description))
	html = replaceMetaContent(html, "property", "og:image", defaults.image, imageURL)
	html = replaceMetaContent(html, "name", "twitter:title", defaults.title, escapeAttr(title))
	html = replaceMetaContent(html, "name", "twitter:description", defaults.twitterDesc, escapeAttr(twitterDesc))
	html = replaceMetaContent(html, "name", "twitter:image", defaults.image, imageURL)
	return html
}

func (s *Service) baseURL(ctx fiber.Ctx) string {
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

func replaceMetaContent(html, attrName, attrValue, oldContent, newContent string) string {
	old := attrName + `="` + attrValue + `" content="` + oldContent + `"`
	repl := attrName + `="` + attrValue + `" content="` + newContent + `"`
	return strings.Replace(html, old, repl, 1)
}

func uminekoEpisodeLabel(ep int, contentType string) string {
	name, ok := uminekoEpisodeNames[ep]
	if !ok {
		return ""
	}
	if contentType == "tea" {
		return name + " \u2014 Tea Party"
	}
	if contentType == "ura" {
		return name + " \u2014 Omake"
	}
	return name
}

func higurashiArcLabel(arc string) string {
	if name, ok := higurashiArcNames[arc]; ok {
		return name
	}
	return arc
}

func escapeAttr(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}
