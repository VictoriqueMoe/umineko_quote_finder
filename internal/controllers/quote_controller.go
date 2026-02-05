package controllers

import (
	"os"
	"regexp"
	"strings"

	"umineko_quote/internal/audio"
	"umineko_quote/internal/quote"
	"umineko_quote/internal/utils"

	"github.com/gofiber/fiber/v2"
)

var audioIdPattern = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

func (s *Service) getAllQuoteRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupSearchRoute,
		s.setupRandomRoute,
		s.setupBrowseRoute,
		s.setupByCharacterRoute,
		s.setupByAudioIDRoute,
		s.setupContextRoute,
		s.setupCharactersRoute,
		s.setupCombinedAudioRoute,
		s.setupAudioRoute,
		s.setupStatsRoute,
	}
}

func (s *Service) setupSearchRoute(routeGroup fiber.Router) {
	routeGroup.Get("/search", s.search)
}

func (s *Service) setupRandomRoute(routeGroup fiber.Router) {
	routeGroup.Get("/random", s.random)
}

func (s *Service) setupBrowseRoute(routeGroup fiber.Router) {
	routeGroup.Get("/browse", s.browse)
}

func (s *Service) setupByCharacterRoute(routeGroup fiber.Router) {
	routeGroup.Get("/character/:id", s.byCharacter)
}

func (s *Service) setupByAudioIDRoute(routeGroup fiber.Router) {
	routeGroup.Get("/quote/:audioId", s.byAudioID)
}

func (s *Service) setupCharactersRoute(routeGroup fiber.Router) {
	routeGroup.Get("/characters", s.characters)
}

func (s *Service) search(ctx *fiber.Ctx) error {
	query := ctx.Query("q")
	if query == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query parameter 'q' is required",
		})
	}

	lang := ctx.Query("lang", "en")
	limit := ctx.QueryInt("limit", 30)
	offset := ctx.QueryInt("offset", 0)
	characterID := ctx.Query("character")
	episode := ctx.QueryInt("episode", 0)
	truth := quote.TruthAll.Parse(ctx.Query("truth"))

	response := s.QuoteService.Search(query, lang, limit, offset, characterID, episode, truth)
	return ctx.JSON(fiber.Map{
		"query":   query,
		"results": response.Results,
		"total":   response.Total,
		"limit":   response.Limit,
		"offset":  response.Offset,
	})
}

func (s *Service) random(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "en")
	characterID := ctx.Query("character")
	episode := ctx.QueryInt("episode", 0)
	truth := quote.TruthAll.Parse(ctx.Query("truth"))
	q := s.QuoteService.Random(lang, characterID, episode, truth)
	if q == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "no quotes available",
		})
	}
	return ctx.JSON(q)
}

func (s *Service) browse(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "en")
	characterID := ctx.Query("character")
	limit := ctx.QueryInt("limit", 50)
	offset := ctx.QueryInt("offset", 0)
	episode := ctx.QueryInt("episode", 0)
	truth := quote.TruthAll.Parse(ctx.Query("truth"))

	response := s.QuoteService.Browse(lang, characterID, limit, offset, episode, truth)
	return ctx.JSON(response)
}

func (s *Service) byCharacter(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "en")
	characterID := ctx.Params("id")
	limit := ctx.QueryInt("limit", 50)
	offset := ctx.QueryInt("offset", 0)
	episode := ctx.QueryInt("episode", 0)
	truth := quote.TruthAll.Parse(ctx.Query("truth"))

	response := s.QuoteService.GetByCharacter(lang, characterID, limit, offset, episode, truth)
	return ctx.JSON(response)
}

func (s *Service) byAudioID(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "en")
	audioID := ctx.Params("audioId")

	quote := s.QuoteService.GetByAudioID(lang, audioID)
	if quote == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(quote)
}

func (s *Service) setupContextRoute(routeGroup fiber.Router) {
	routeGroup.Get("/context/:audioId", s.context)
}

func (s *Service) context(ctx *fiber.Ctx) error {
	lang := ctx.Query("lang", "en")
	audioID := ctx.Params("audioId")
	if !audioIdPattern.MatchString(audioID) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid audio ID",
		})
	}

	lines := ctx.QueryInt("lines", 5)
	result := s.QuoteService.GetContext(lang, audioID, lines)
	if result == nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "quote not found",
		})
	}
	return ctx.JSON(result)
}

func (s *Service) characters(ctx *fiber.Ctx) error {
	return ctx.JSON(s.QuoteService.GetCharacters())
}

func (s *Service) setupStatsRoute(routeGroup fiber.Router) {
	routeGroup.Get("/stats", s.stats)
}

func (s *Service) stats(ctx *fiber.Ctx) error {
	episode := ctx.QueryInt("episode", 0)
	return ctx.JSON(s.QuoteService.GetStats().Compute(episode))
}

func (s *Service) setupCombinedAudioRoute(routeGroup fiber.Router) {
	routeGroup.Get("/audio/combined", s.combinedAudioSegments)
	routeGroup.Get("/audio/:charId/combined", s.combinedAudioLegacy)
}

func (s *Service) setupAudioRoute(routeGroup fiber.Router) {
	routeGroup.Get("/audio/:charId/:audioId", s.audio)
}

func (s *Service) audio(ctx *fiber.Ctx) error {
	charId := ctx.Params("charId")
	audioId := ctx.Params("audioId")
	if !audioIdPattern.MatchString(charId) || !audioIdPattern.MatchString(audioId) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid audio ID",
		})
	}

	filePath := s.QuoteService.AudioFilePath(charId, audioId)
	if filePath == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "audio file not found",
		})
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read audio file",
		})
	}

	return utils.ServeAudio(ctx, data)
}

func (s *Service) combinedAudioSegments(ctx *fiber.Ctx) error {
	segmentsParam := ctx.Query("segments")
	if segmentsParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query parameter 'segments' is required",
		})
	}

	parts := strings.Split(segmentsParam, ",")
	if len(parts) > 20 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "maximum 20 audio segments allowed",
		})
	}

	segments := make([]audio.AudioSegment, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		colonIdx := strings.IndexByte(part, ':')
		if colonIdx < 1 || colonIdx >= len(part)-1 {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid segment format: " + part + " (expected charId:audioId)",
			})
		}
		charId := part[:colonIdx]
		audioId := part[colonIdx+1:]
		if !audioIdPattern.MatchString(charId) || !audioIdPattern.MatchString(audioId) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid segment: " + part,
			})
		}
		segments = append(segments, audio.AudioSegment{CharID: charId, AudioID: audioId})
	}

	data, err := s.AudioCombiner.CombineOgg(segments, s.QuoteService.AudioFilePath)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.ServeAudio(ctx, data)
}

func (s *Service) combinedAudioLegacy(ctx *fiber.Ctx) error {
	charId := ctx.Params("charId")
	if !audioIdPattern.MatchString(charId) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid character ID",
		})
	}

	idsParam := ctx.Query("ids")
	if idsParam == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "query parameter 'ids' is required",
		})
	}

	ids := strings.Split(idsParam, ",")
	if len(ids) > 20 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "maximum 20 audio IDs allowed",
		})
	}

	segments := make([]audio.AudioSegment, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if !audioIdPattern.MatchString(id) {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid audio ID: " + id,
			})
		}
		segments = append(segments, audio.AudioSegment{CharID: charId, AudioID: id})
	}

	data, err := s.AudioCombiner.CombineOgg(segments, s.QuoteService.AudioFilePath)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return utils.ServeAudio(ctx, data)
}
