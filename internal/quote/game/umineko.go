package game

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"maps"
	"slices"
	"strconv"
	"sync"
	"time"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/subtitle"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/umineko"
)

type (
	Umineko struct{}

	uminekoLangParseResult struct {
		lang   language.Language
		parsed []dto.UminekoQuote
		err    error
	}
)

const (
	UminekoAudioDir     = "internal/quote/data/audio"
	uminekoSubtitleGlob = "data/sub/*.ass"
)

var uminekoLangFiles = map[language.Language]string{
	language.English:    "data/en.file",
	language.WitchHunt:  "data/wh.file",
	language.Japanese:   "data/ja.file",
	language.Russian:    "data/ru.file",
	language.Spanish:    "data/es.file",
	language.Portuguese: "data/pt.file",
}

func (Umineko) Name() string {
	return "umineko"
}

func (Umineko) DataFiles() ([]string, error) {
	paths := slices.Sorted(maps.Values(uminekoLangFiles))

	subFiles, err := fs.Glob(quote.DataFS, uminekoSubtitleGlob)
	if err != nil {
		return nil, err
	}
	slices.Sort(subFiles)

	return append(paths, subFiles...), nil
}

func (Umineko) Load() (map[language.Language][]dto.UminekoQuote, error) {
	results := make(chan uminekoLangParseResult, len(uminekoLangFiles))
	var wg sync.WaitGroup

	for lang, path := range uminekoLangFiles {
		wg.Go(func() {
			parsed, err := loadUminekoLanguage(lang, path)
			results <- uminekoLangParseResult{lang: lang, parsed: parsed, err: err}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	quotes := make(map[language.Language][]dto.UminekoQuote)
	var errs []error

	for r := range results {
		if r.err != nil {
			errs = append(errs, r.err)
			continue
		}
		quotes[r.lang] = r.parsed
	}

	if err := errors.Join(errs...); err != nil {
		return nil, err
	}

	return quotes, nil
}

func loadUminekoLanguage(lang language.Language, path string) ([]dto.UminekoQuote, error) {
	f, err := quote.DataFS.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open umineko %s script %s: %w", lang, path, err)
	}
	defer f.Close()

	timeStart := time.Now()
	parser := umineko.NewParser()
	parsed, validationErrors, err := scriptparser.ParseReader(f, parser)
	subtitleRefs := parser.SubtitleRefs()
	timeEnd := time.Now()

	if err != nil {
		return nil, fmt.Errorf("parse umineko %s script %s: %w", lang, path, err)
	}
	log.Printf("[umineko/%s] parsed %s: %d quotes in %v", lang, path, len(parsed), timeEnd.Sub(timeStart))

	if len(validationErrors) > 0 {
		errorCount := 0
		warningCount := 0
		for _, ve := range validationErrors {
			if ve.Severity == scriptparser.SeverityError {
				errorCount++
			} else {
				warningCount++
			}
		}
		log.Printf("[umineko/%s] validation: %d errors, %d warnings", lang, errorCount, warningCount)
		limit := min(len(validationErrors), 10)
		for _, ve := range validationErrors[:limit] {
			log.Printf("[umineko/%s]   %s", lang, ve)
		}
		if len(validationErrors) > 10 {
			log.Printf("[umineko/%s]   ... and %d more", lang, len(validationErrors)-10)
		}
	}

	indexed := make([]dto.UminekoQuote, len(parsed))
	for i, q := range parsed {
		indexed[i] = dto.UminekoQuote{
			UminekoQuote: q,
			Index:        i,
		}
	}

	subQuotes := subtitle.ResolveRefs(quote.DataFS, subtitleRefs)
	if len(subQuotes) > 0 {
		startIdx := len(indexed)
		for i, sq := range subQuotes {
			indexed = append(indexed, dto.UminekoQuote{
				UminekoQuote: sq.UminekoQuote,
				Index:        startIdx + i,
			})
		}
		log.Printf("[umineko/%s] added %d subtitle quotes", lang, len(subQuotes))
	}

	return indexed, nil
}

func (Umineko) Base(q *dto.UminekoQuote) *scriptdto.BaseQuote {
	return &q.BaseQuote
}

func (Umineko) Row(q *dto.UminekoQuote) QuoteRow {
	return QuoteRow{
		Index:          q.Index,
		Segment:        strconv.Itoa(q.Episode),
		HasRedTruth:    q.HasRedTruth,
		HasBlueTruth:   q.HasBlueTruth,
		HasGoldTruth:   q.HasGoldTruth,
		HasPurpleTruth: q.HasPurpleTruth,
	}
}

func (Umineko) Build(base scriptdto.BaseQuote, r QuoteRow) dto.UminekoQuote {
	return dto.UminekoQuote{
		BaseQuote:      base,
		HasRedTruth:    r.HasRedTruth,
		HasBlueTruth:   r.HasBlueTruth,
		HasGoldTruth:   r.HasGoldTruth,
		HasPurpleTruth: r.HasPurpleTruth,
		Index:          r.Index,
	}
}
