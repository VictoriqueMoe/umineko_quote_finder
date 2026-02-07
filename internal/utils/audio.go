package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func ServeAudio(ctx fiber.Ctx, data []byte) error {
	ctx.Set("Content-Type", "audio/ogg")
	ctx.Set("Accept-Ranges", "bytes")

	total := len(data)
	rangeHeader := ctx.Get("Range")
	if rangeHeader == "" {
		return ctx.Send(data)
	}

	if !strings.HasPrefix(rangeHeader, "bytes=") {
		ctx.Set("Content-Range", fmt.Sprintf("bytes */%d", total))
		return ctx.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
	}

	rangeParts := strings.SplitN(rangeHeader[6:], "-", 2)
	if len(rangeParts) != 2 {
		ctx.Set("Content-Range", fmt.Sprintf("bytes */%d", total))
		return ctx.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
	}

	var start, end int
	var err error

	if rangeParts[0] == "" {
		suffixLen, parseErr := strconv.Atoi(rangeParts[1])
		if parseErr != nil || suffixLen <= 0 {
			ctx.Set("Content-Range", fmt.Sprintf("bytes */%d", total))
			return ctx.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
		}
		start = total - suffixLen
		if start < 0 {
			start = 0
		}
		end = total - 1
	} else {
		start, err = strconv.Atoi(rangeParts[0])
		if err != nil || start < 0 || start >= total {
			ctx.Set("Content-Range", fmt.Sprintf("bytes */%d", total))
			return ctx.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
		}
		if rangeParts[1] == "" {
			end = total - 1
		} else {
			end, err = strconv.Atoi(rangeParts[1])
			if err != nil || end < start || end >= total {
				ctx.Set("Content-Range", fmt.Sprintf("bytes */%d", total))
				return ctx.SendStatus(fiber.StatusRequestedRangeNotSatisfiable)
			}
		}
	}

	ctx.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, total))
	ctx.Set("Content-Length", strconv.Itoa(end-start+1))
	ctx.Status(fiber.StatusPartialContent)
	return ctx.Send(data[start : end+1])
}
