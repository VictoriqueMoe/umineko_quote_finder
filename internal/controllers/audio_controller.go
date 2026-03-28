package controllers

import (
	"os"
	"path/filepath"
	"strings"

	"umineko_quote/internal/audio"
	"umineko_quote/internal/utils"

	"github.com/gofiber/fiber/v3"
)

const seAudioDir = "internal/quote/data/se"

func (s *Service) getAllAudioRoutes() []FSetupRoute {
	return []FSetupRoute{
		s.setupCombinedAudioRoute,
		s.setupSeAudioRoute,
		s.setupAudioRoute,
	}
}

func (s *Service) setupCombinedAudioRoute(routeGroup fiber.Router) {
	routeGroup.Get("/audio/voice/combined", s.combinedAudioSegments)
}

func (s *Service) setupSeAudioRoute(routeGroup fiber.Router) {
	routeGroup.Get("/audio/se/:filename", s.seAudio)
}

func (s *Service) setupAudioRoute(routeGroup fiber.Router) {
	routeGroup.Get("/audio/voice/:charId/:audioId", s.audio)
}

func (s *Service) audio(ctx fiber.Ctx) error {
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

func (s *Service) seAudio(ctx fiber.Ctx) error {
	filename := ctx.Params("filename")
	if !audioIdPattern.MatchString(filename) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid SE filename",
		})
	}

	filePath := filepath.Join(seAudioDir, filename+".ogg")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "SE file not found",
		})
	}

	return utils.ServeAudio(ctx, data)
}

func (s *Service) combinedAudioSegments(ctx fiber.Ctx) error {
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
