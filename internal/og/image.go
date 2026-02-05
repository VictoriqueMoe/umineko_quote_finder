package og

import (
	"bytes"
	"embed"
	"fmt"
	"image/color"
	"image/png"
	"strings"
	"sync"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

//go:embed fonts/NotoSansJP-Regular.ttf
var notoSansJPData embed.FS

const (
	imgWidth  = 1200
	imgHeight = 630
)

var (
	bgColor        = color.RGBA{R: 10, G: 6, B: 18, A: 255}
	goldColor      = color.RGBA{R: 212, G: 168, B: 75, A: 255}
	creamColor     = color.RGBA{R: 232, G: 224, B: 240, A: 255}
	mutedColor     = color.RGBA{R: 168, G: 155, B: 184, A: 255}
	accentBarColor = color.RGBA{R: 212, G: 168, B: 75, A: 255}
)

type ImageGenerator struct {
	regularFont *sfnt.Font
	boldFont    *sfnt.Font
	jpFont      *sfnt.Font
	cache       sync.Map
}

func NewImageGenerator() *ImageGenerator {
	regular, _ := opentype.Parse(goregular.TTF)
	bold, _ := opentype.Parse(gobold.TTF)

	jpData, _ := notoSansJPData.ReadFile("fonts/NotoSansJP-Regular.ttf")
	jp, _ := opentype.Parse(jpData)

	return &ImageGenerator{
		regularFont: regular,
		boldFont:    bold,
		jpFont:      jp,
	}
}

func (g *ImageGenerator) textFont(lang string) *sfnt.Font {
	if lang == "ja" {
		return g.jpFont
	}
	return g.regularFont
}

func (g *ImageGenerator) boldOrFallback(lang string) *sfnt.Font {
	if lang == "ja" {
		return g.jpFont
	}
	return g.boldFont
}

func (g *ImageGenerator) finalise(dc *gg.Context, cacheKey string) ([]byte, error) {
	brandFace, err := opentype.NewFace(g.regularFont, &opentype.FaceOptions{Size: 16, DPI: 72})
	if err != nil {
		return nil, err
	}
	dc.SetFontFace(brandFace)
	dc.SetColor(mutedColor)
	dc.DrawStringAnchored("Umineko Quote Search", float64(imgWidth)-40, float64(imgHeight)-30, 1, 0)

	var buf bytes.Buffer
	if err := png.Encode(&buf, dc.Image()); err != nil {
		return nil, err
	}

	data := buf.Bytes()
	g.cache.Store(cacheKey, data)
	return data, nil
}

type DialogueLine struct {
	Character string
	Text      string
}

func (g *ImageGenerator) Generate(audioId, lang, text, textHtml, character string, episode int, contentType string) ([]byte, error) {
	cacheKey := audioId + ":" + lang
	if cached, ok := g.cache.Load(cacheKey); ok {
		return cached.([]byte), nil
	}

	dc := gg.NewContext(imgWidth, imgHeight)

	dc.SetColor(bgColor)
	dc.Clear()

	dc.SetColor(accentBarColor)
	dc.DrawRectangle(0, 0, 8, float64(imgHeight))
	dc.Fill()

	quoteFace, err := opentype.NewFace(g.boldFont, &opentype.FaceOptions{Size: 72, DPI: 72})
	if err != nil {
		return nil, err
	}
	dc.SetFontFace(quoteFace)
	dc.SetColor(color.RGBA{R: 212, G: 168, B: 75, A: 60})
	dc.DrawString("\u201C", 40, 100)

	textFace, err := opentype.NewFace(g.textFont(lang), &opentype.FaceOptions{Size: 28, DPI: 72})
	if err != nil {
		return nil, err
	}
	dc.SetFontFace(textFace)

	maxWidth := float64(imgWidth) - 120
	if textHtml != "" {
		segments := parseHTMLSegments(textHtml, creamColor)
		segments = truncateSegments(segments, 300)
		g.drawColouredText(dc, segments, 60, 120, maxWidth, 1.5)
	} else {
		displayText := truncateText(text, 300)
		dc.SetColor(creamColor)
		dc.DrawStringWrapped(displayText, 60, 120, 0, 0, maxWidth, 1.5, gg.AlignLeft)
	}

	charFace, err := opentype.NewFace(g.boldOrFallback(lang), &opentype.FaceOptions{Size: 24, DPI: 72})
	if err != nil {
		return nil, err
	}
	dc.SetFontFace(charFace)
	dc.SetColor(goldColor)
	dc.DrawString("\u2014 "+character, 60, float64(imgHeight)-120)

	if episode > 0 {
		epFace, err := opentype.NewFace(g.regularFont, &opentype.FaceOptions{Size: 18, DPI: 72})
		if err != nil {
			return nil, err
		}
		dc.SetFontFace(epFace)
		dc.SetColor(mutedColor)
		dc.DrawString(g.episodeName(episode, contentType), 60, float64(imgHeight)-88)
	}

	return g.finalise(dc, cacheKey)
}

func (g *ImageGenerator) GenerateBuilder(segmentsParam, lang string, lines []DialogueLine) ([]byte, error) {
	cacheKey := "builder:" + segmentsParam + ":" + lang
	if cached, ok := g.cache.Load(cacheKey); ok {
		return cached.([]byte), nil
	}

	dc := gg.NewContext(imgWidth, imgHeight)

	dc.SetColor(bgColor)
	dc.Clear()

	dc.SetColor(accentBarColor)
	dc.DrawRectangle(0, 0, 8, float64(imgHeight))
	dc.Fill()

	titleFace, err := opentype.NewFace(g.boldFont, &opentype.FaceOptions{Size: 36, DPI: 72})
	if err != nil {
		return nil, err
	}
	dc.SetFontFace(titleFace)
	dc.SetColor(goldColor)
	dc.DrawString("\u266B Voice Build", 60, 60)

	dc.SetColor(color.RGBA{R: 212, G: 168, B: 75, A: 80})
	dc.DrawRectangle(60, 75, float64(imgWidth)-120, 1)
	dc.Fill()

	charFace, err := opentype.NewFace(g.boldOrFallback(lang), &opentype.FaceOptions{Size: 20, DPI: 72})
	if err != nil {
		return nil, err
	}
	defer func() { _ = charFace.Close() }()
	textFace, err := opentype.NewFace(g.textFont(lang), &opentype.FaceOptions{Size: 20, DPI: 72})
	if err != nil {
		return nil, err
	}
	defer func() { _ = textFace.Close() }()

	maxWidth := float64(imgWidth) - 140
	curY := 110.0
	maxY := float64(imgHeight) - 80

	for i, line := range lines {
		if curY > maxY {
			dc.SetFontFace(textFace)
			dc.SetColor(mutedColor)
			remaining := len(lines) - i
			if remaining == 1 {
				dc.DrawString("+ 1 more clip\u2026", 80, curY)
			} else {
				dc.DrawString(fmt.Sprintf("+ %d more clips\u2026", remaining), 80, curY)
			}
			break
		}

		dc.SetFontFace(charFace)
		dc.SetColor(goldColor)
		nameStr := line.Character + ":"
		nameW, _ := dc.MeasureString(nameStr)
		dc.DrawString(nameStr, 80, curY)

		dc.SetFontFace(textFace)
		dc.SetColor(creamColor)
		quoteText := truncateText(line.Text, 120)

		textX := 80 + nameW + 10
		remainingWidth := maxWidth - nameW - 10
		if remainingWidth < 200 {
			curY += 30
			dc.DrawStringWrapped(quoteText, 80, curY, 0, 0, maxWidth, 1.4, gg.AlignLeft)
			_, textH := dc.MeasureString("Mg")
			wrappedLines := dc.WordWrap(quoteText, maxWidth)
			curY += textH * 1.4 * float64(len(wrappedLines))
		} else {
			wrappedLines := dc.WordWrap(quoteText, remainingWidth)
			if len(wrappedLines) > 0 {
				dc.DrawString(wrappedLines[0], textX, curY)
			}
			_, textH := dc.MeasureString("Mg")
			for j := 1; j < len(wrappedLines); j++ {
				curY += textH * 1.4
				if curY > maxY {
					break
				}
				dc.DrawString(wrappedLines[j], 80, curY)
			}
			curY += textH * 1.4
		}

		curY += 8
	}

	return g.finalise(dc, cacheKey)
}

func (*ImageGenerator) drawColouredText(dc *gg.Context, segments []textSegment, x, y, maxWidth, lineSpacing float64) {
	_, fh := dc.MeasureString("Mg")
	lineH := fh * lineSpacing
	curX := x
	curY := y

	for _, seg := range segments {
		tokens := splitTokens(seg.Text)
		for _, tok := range tokens {
			if tok == "\n" {
				curX = x
				curY += lineH
				continue
			}

			w, _ := dc.MeasureString(tok)

			if w > maxWidth && tok != " " {
				for _, r := range tok {
					rs := string(r)
					rw, _ := dc.MeasureString(rs)
					if curX > x && (curX-x)+rw > maxWidth {
						curX = x
						curY += lineH
					}
					dc.SetColor(seg.Color)
					dc.DrawString(rs, curX, curY)
					curX += rw
				}
				continue
			}

			if curX > x && (curX-x)+w > maxWidth {
				curX = x
				curY += lineH
			}

			if curX == x && tok == " " {
				continue
			}

			dc.SetColor(seg.Color)
			dc.DrawString(tok, curX, curY)
			curX += w
		}
	}
}

func splitTokens(s string) []string {
	var tokens []string
	var buf strings.Builder

	for _, r := range s {
		switch r {
		case '\n':
			if buf.Len() > 0 {
				tokens = append(tokens, buf.String())
				buf.Reset()
			}
			tokens = append(tokens, "\n")
		case ' ':
			if buf.Len() > 0 {
				tokens = append(tokens, buf.String())
				buf.Reset()
			}
			tokens = append(tokens, " ")
		default:
			buf.WriteRune(r)
		}
	}

	if buf.Len() > 0 {
		tokens = append(tokens, buf.String())
	}

	return tokens
}

func truncateText(s string, maxRunes int) string {
	runes := []rune(s)
	if len(runes) > maxRunes {
		return string(runes[:maxRunes-3]) + "..."
	}
	return s
}

func (*ImageGenerator) episodeName(ep int, contentType string) string {
	names := map[int]string{
		1: "Episode 1 \u2014 Legend",
		2: "Episode 2 \u2014 Turn",
		3: "Episode 3 \u2014 Banquet",
		4: "Episode 4 \u2014 Alliance",
		5: "Episode 5 \u2014 End",
		6: "Episode 6 \u2014 Dawn",
		7: "Episode 7 \u2014 Requiem",
		8: "Episode 8 \u2014 Twilight",
	}
	name, ok := names[ep]
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
