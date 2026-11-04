package game

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"slices"
	"sync"
	"time"

	"umineko_quote/internal/quote"
	"umineko_quote/internal/quote/language"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
)

type (
	Game[Q any] interface {
		Name() string
		Load() (map[language.Language][]Q, error)
		DataFiles() ([]string, error)
		Base(q *Q) *scriptdto.BaseQuote
		Row(q *Q) QuoteRow
		Build(base scriptdto.BaseQuote, r QuoteRow) Q
	}

	QuoteRow struct {
		Index          int
		Segment        string
		HasRedTruth    bool
		HasBlueTruth   bool
		HasGoldTruth   bool
		HasPurpleTruth bool
	}

	languageResult[Q any] struct {
		lang   language.Language
		parsed []Q
		err    error
	}

	scriptSource[P any] struct {
		path   string
		parser scriptparser.Parser[P]
	}

	scriptLoader[P any, Q any] struct {
		name  string
		index func(P, int) Q
	}
)

const NarratorID = "narrator"

func loadLanguages[S any, Q any](sources map[language.Language]S, load func(language.Language, S) ([]Q, error)) (map[language.Language][]Q, error) {
	results := make(chan languageResult[Q], len(sources))
	var wg sync.WaitGroup

	for lang, src := range sources {
		wg.Go(func() {
			parsed, err := load(lang, src)
			results <- languageResult[Q]{lang: lang, parsed: parsed, err: err}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	quotes := make(map[language.Language][]Q, len(sources))
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

func (l scriptLoader[P, Q]) load(lang language.Language, src scriptSource[P]) ([]Q, error) {
	data, err := fs.ReadFile(quote.DataFS, src.path)
	if err != nil {
		return nil, fmt.Errorf("read %s %s script %s: %w", l.name, lang, src.path, err)
	}

	start := time.Now()
	parsed, _, err := scriptparser.ParseReader(bytes.NewReader(data), src.parser)
	if err != nil {
		return nil, fmt.Errorf("parse %s %s script %s: %w", l.name, lang, src.path, err)
	}
	log.Printf("[%s/%s] parsed %d quotes in %v", l.name, lang, len(parsed), time.Since(start))

	indexed := make([]Q, len(parsed))
	for i, q := range parsed {
		indexed[i] = l.index(q, i)
	}

	return indexed, nil
}

func sourceFiles[P any](sources map[language.Language]scriptSource[P]) []string {
	paths := make([]string, 0, len(sources))
	for _, src := range sources {
		paths = append(paths, src.path)
	}
	slices.Sort(paths)

	return slices.Compact(paths)
}
