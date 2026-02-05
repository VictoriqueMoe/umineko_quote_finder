package og

import (
	"image/color"
	"strconv"
	"strings"
)

type textSegment struct {
	Text  string
	Color color.RGBA
}

var (
	redTruthColor  = color.RGBA{R: 255, G: 51, B: 51, A: 255}  // #ff3333
	blueTruthColor = color.RGBA{R: 51, G: 153, B: 255, A: 255} // #3399ff
)

func parseHTMLSegments(htmlStr string, defaultColor color.RGBA) []textSegment {
	var segments []textSegment
	var colorStack []color.RGBA
	currentColor := defaultColor
	var buf strings.Builder
	var skipDepth int

	i := 0
	for i < len(htmlStr) {
		if htmlStr[i] == '<' {
			if buf.Len() > 0 && skipDepth == 0 {
				segments = append(segments, textSegment{Text: buf.String(), Color: currentColor})
			}
			buf.Reset()

			end := strings.IndexByte(htmlStr[i:], '>')
			if end == -1 {
				break
			}
			tag := htmlStr[i : i+end+1]
			i += end + 1

			lowerTag := strings.ToLower(tag)

			switch {
			case lowerTag == "<br>" || lowerTag == "<br/>" || lowerTag == "<br />":
				if skipDepth == 0 {
					buf.WriteRune('\n')
				}
			case strings.HasPrefix(lowerTag, "<span"):
				colorStack = append(colorStack, currentColor)
				if strings.Contains(tag, `class="red-truth"`) {
					currentColor = redTruthColor
				} else if strings.Contains(tag, `class="blue-truth"`) {
					currentColor = blueTruthColor
				} else if idx := strings.Index(tag, "color:"); idx != -1 {
					if c, ok := parseHexColor(tag[idx+6:]); ok {
						currentColor = c
					}
				}
			case lowerTag == "</span>":
				if len(colorStack) > 0 {
					currentColor = colorStack[len(colorStack)-1]
					colorStack = colorStack[:len(colorStack)-1]
				}
			case strings.HasPrefix(lowerTag, "<rt") || strings.HasPrefix(lowerTag, "<rp"):
				skipDepth++
			case lowerTag == "</rt>" || lowerTag == "</rp>":
				if skipDepth > 0 {
					skipDepth--
				}
			}
		} else if htmlStr[i] == '&' {
			end := strings.IndexByte(htmlStr[i:], ';')
			if end != -1 {
				entity := htmlStr[i : i+end+1]
				i += end + 1
				if skipDepth == 0 {
					switch entity {
					case "&amp;":
						buf.WriteRune('&')
					case "&lt;":
						buf.WriteRune('<')
					case "&gt;":
						buf.WriteRune('>')
					case "&quot;":
						buf.WriteRune('"')
					case "&apos;":
						buf.WriteRune('\'')
					default:
						if r, ok := decodeNumericEntity(entity); ok {
							buf.WriteRune(r)
						} else {
							buf.WriteString(entity)
						}
					}
				}
			} else {
				if skipDepth == 0 {
					buf.WriteByte(htmlStr[i])
				}
				i++
			}
		} else {
			if skipDepth == 0 {
				buf.WriteByte(htmlStr[i])
			}
			i++
		}
	}

	if buf.Len() > 0 && skipDepth == 0 {
		segments = append(segments, textSegment{Text: buf.String(), Color: currentColor})
	}

	return segments
}

func truncateSegments(segments []textSegment, maxRunes int) []textSegment {
	total := 0
	for _, seg := range segments {
		total += len([]rune(seg.Text))
	}
	if total <= maxRunes {
		return segments
	}

	var result []textSegment
	remaining := maxRunes - 3
	if remaining < 0 {
		remaining = 0
	}

	for _, seg := range segments {
		runes := []rune(seg.Text)
		if remaining <= 0 {
			break
		}
		if len(runes) <= remaining {
			result = append(result, seg)
			remaining -= len(runes)
		} else {
			result = append(result, textSegment{Text: string(runes[:remaining]), Color: seg.Color})
			remaining = 0
		}
	}

	if len(result) > 0 {
		last := &result[len(result)-1]
		last.Text += "..."
	} else {
		result = append(result, textSegment{Text: "...", Color: creamColor})
	}

	return result
}

func decodeNumericEntity(entity string) (rune, bool) {
	if len(entity) < 4 || entity[0] != '&' || entity[len(entity)-1] != ';' {
		return 0, false
	}
	inner := entity[1 : len(entity)-1]
	if len(inner) < 2 || inner[0] != '#' {
		return 0, false
	}
	inner = inner[1:]
	var n uint64
	var err error
	if inner[0] == 'x' || inner[0] == 'X' {
		n, err = strconv.ParseUint(inner[1:], 16, 32)
	} else {
		n, err = strconv.ParseUint(inner, 10, 32)
	}
	if err != nil || n == 0 {
		return 0, false
	}
	return rune(n), true
}

func parseHexColor(s string) (color.RGBA, bool) {
	end := strings.IndexAny(s, "\";)")
	if end == -1 {
		return color.RGBA{}, false
	}
	hex := strings.TrimSpace(s[:end])
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return color.RGBA{}, false
	}
	r, err1 := strconv.ParseUint(hex[0:2], 16, 8)
	g, err2 := strconv.ParseUint(hex[2:4], 16, 8)
	b, err3 := strconv.ParseUint(hex[4:6], 16, 8)
	if err1 != nil || err2 != nil || err3 != nil {
		return color.RGBA{}, false
	}
	return color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}, true
}
