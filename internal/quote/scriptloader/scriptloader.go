package scriptloader

import (
	"bytes"
	"compress/zlib"
	"embed"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/lexar"
	"umineko_quote/internal/quote/character"
	"umineko_quote/internal/quote/mutation"
	"umineko_quote/internal/subtitle"
)

const (
	pass1XorA byte = 0x45
	pass1XorB byte = 0x71
	pass2XorA byte = 0x86
	pass2XorB byte = 0x23
)

type (
	ParseFunc func(lines []string) ([]dto.ParsedQuote, []lexar.SubtitleRef)

	Loader struct {
		fs        embed.FS
		parse     ParseFunc
		mutations mutation.Pipeline
	}

	xorReader struct {
		r    io.Reader
		xorA byte
		xorB byte
	}
)

var (
	inverseKeyTable = [256]byte{
		0x37, 0x6a, 0x09, 0x5e, 0x7a, 0xaf, 0xf5, 0xa4, 0xba, 0x78, 0x84, 0x58, 0x35, 0x1e, 0x6b, 0x0c,
		0x49, 0xc6, 0xc3, 0x44, 0x40, 0x9e, 0x6f, 0x65, 0xe4, 0xf6, 0xfe, 0x22, 0xe2, 0x95, 0xc7, 0x38,
		0xf0, 0x1a, 0x82, 0xe0, 0x5b, 0x2a, 0xd8, 0xe5, 0xce, 0x2f, 0x74, 0x25, 0xec, 0x59, 0xc0, 0x45,
		0x4b, 0x64, 0x43, 0xdc, 0xb0, 0xb9, 0x30, 0x6d, 0x28, 0xd1, 0x16, 0xbb, 0x66, 0x98, 0x92, 0x90,
		0x2c, 0xa7, 0xf1, 0x80, 0xc1, 0xd4, 0x8b, 0xd6, 0xdf, 0x24, 0x2d, 0xf7, 0xfb, 0x88, 0x4d, 0x3c,
		0x72, 0xf3, 0xdb, 0x2b, 0x93, 0x73, 0xef, 0x85, 0x83, 0xee, 0xc2, 0x8d, 0x5c, 0xb2, 0x0b, 0x94,
		0x3d, 0xa8, 0x3f, 0x1c, 0x4c, 0x6e, 0x03, 0x7b, 0x1d, 0x5a, 0x51, 0xa1, 0x70, 0x41, 0xd0, 0xaa,
		0xa0, 0x7e, 0xcd, 0xd5, 0x15, 0xa9, 0x18, 0x76, 0xc9, 0x7d, 0x7f, 0x0e, 0x3a, 0x99, 0xbf, 0xab,
		0x3b, 0x14, 0x3e, 0x9a, 0x04, 0xda, 0x02, 0xfd, 0x63, 0xd9, 0xfa, 0x9f, 0x4e, 0xe3, 0x61, 0xbe,
		0x07, 0x11, 0xa6, 0x1b, 0x19, 0x55, 0x8e, 0x77, 0x0a, 0x47, 0xe6, 0xf8, 0x0d, 0xcf, 0xd7, 0x33,
		0x23, 0x1f, 0xbc, 0x62, 0xde, 0x9b, 0x29, 0x53, 0x68, 0xe8, 0x21, 0xb6, 0x34, 0x52, 0x87, 0xcb,
		0x08, 0x79, 0xf4, 0x67, 0x69, 0x54, 0xe7, 0x86, 0xea, 0xb4, 0x20, 0x71, 0x01, 0xbd, 0x06, 0x31,
		0x00, 0x50, 0xc8, 0xb8, 0xac, 0x5d, 0x57, 0x7c, 0x89, 0xeb, 0xb7, 0x36, 0x8f, 0xf2, 0xe1, 0x56,
		0x81, 0x4a, 0xd2, 0x8c, 0xf9, 0xad, 0x60, 0xa5, 0x42, 0x10, 0x5f, 0x12, 0xb3, 0xff, 0x4f, 0xdd,
		0x46, 0x26, 0xa2, 0x17, 0xc5, 0x75, 0x91, 0x27, 0xb5, 0x8a, 0xd3, 0x13, 0x2e, 0xc4, 0xe9, 0x9d,
		0x97, 0x39, 0x32, 0x05, 0x0f, 0xca, 0xcc, 0x48, 0xfc, 0xae, 0x96, 0xed, 0x6c, 0x9c, 0xb1, 0xa3,
	}

	subtitleStyleCharacter = map[string]string{
		"Battler": "10",
	}
)

func New(efs embed.FS, parse ParseFunc) *Loader {
	return &Loader{
		fs:        efs,
		parse:     parse,
		mutations: *mutation.NewPipeline(),
	}
}

func (l *Loader) Load(lang string, path string) []dto.ParsedQuote {
	raw, err := l.fs.ReadFile(path)
	if err != nil {
		log.Printf("[%s] failed to read %s: %v", lang, path, err)
		return nil
	}

	decodeStart := time.Now()
	decoded, err := decode(raw)
	if err != nil {
		log.Printf("[%s] failed to decode %s: %v", lang, path, err)
		return nil
	}
	log.Printf("[%s] decoded %s (%d → %d bytes) in %v", lang, path, len(raw), len(decoded), time.Since(decodeStart).Round(time.Millisecond))

	lines := strings.Split(string(decoded), "\n")

	parseStart := time.Now()
	parsed, subtitleRefs := l.parse(lines)
	log.Printf("[%s] parsed %d lines → %d quotes in %v", lang, len(lines), len(parsed), time.Since(parseStart).Round(time.Millisecond))

	subQuotes := l.resolveSubtitleRefs(subtitleRefs)
	if len(subQuotes) > 0 {
		parsed = append(parsed, subQuotes...)
		log.Printf("[%s] added %d subtitle quotes", lang, len(subQuotes))
	}

	parsed = l.mutations.Apply(parsed)

	return parsed
}

func (x *xorReader) Read(p []byte) (int, error) {
	n, err := x.r.Read(p)
	for i := 0; i < n; i++ {
		p[i] = inverseKeyTable[p[i]^x.xorA] ^ x.xorB
	}
	return n, err
}

func decode(data []byte) ([]byte, error) {
	if len(data) < 16 {
		return nil, fmt.Errorf("data too short: %d bytes", len(data))
	}

	magic := string(data[:4])
	if magic != "ONS2" {
		return nil, fmt.Errorf("invalid magic: expected ONS2, got %s", magic)
	}

	compressedLen := binary.LittleEndian.Uint32(data[4:8])

	pass2 := &xorReader{
		r:    io.LimitReader(bytes.NewReader(data[16:]), int64(compressedLen)),
		xorA: pass2XorA,
		xorB: pass2XorB,
	}

	zlibReader, err := zlib.NewReader(pass2)
	if err != nil {
		return nil, fmt.Errorf("zlib init: %w", err)
	}
	defer zlibReader.Close()

	pass1 := &xorReader{
		r:    zlibReader,
		xorA: pass1XorA,
		xorB: pass1XorB,
	}

	decoded, err := io.ReadAll(pass1)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	return decoded, nil
}

func (l *Loader) resolveSubtitleRefs(refs []lexar.SubtitleRef) []dto.ParsedQuote {
	var quotes []dto.ParsedQuote

	for _, ref := range refs {
		filename := filepath.Base(strings.ReplaceAll(ref.SubPath, `\`, "/"))
		data, err := l.fs.ReadFile("data/sub/" + filename)
		if err != nil {
			log.Printf("[subtitle] could not read %s: %v", filename, err)
			continue
		}

		entries := subtitle.ParseASS(data)
		for i, entry := range entries {
			charID := ref.CharacterID
			if mapped, ok := subtitleStyleCharacter[entry.Style]; ok {
				charID = mapped
			}

			quotes = append(quotes, dto.ParsedQuote{
				Text:        entry.Text,
				TextHtml:    entry.Text,
				CharacterID: charID,
				Character:   character.CharacterNames.GetCharacterName(character.CharacterFromID(charID)),
				AudioID:     fmt.Sprintf("%s_s%d", ref.AudioID, i),
				Episode:     ref.Episode,
				ContentType: ref.ContentType,
			})
		}
	}

	return quotes
}
