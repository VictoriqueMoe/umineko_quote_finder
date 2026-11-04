package game

import (
	"fmt"
	"log"
	"time"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote"
	"umineko_quote/internal/quote/language"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
	"github.com/VictoriqueMoe/umineko_script_parser/ciconia"
	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
)

type Ciconia struct{}

const ciconiaDataFile = "data/Ciconia/en.file"

func (Ciconia) Name() string {
	return "ciconia"
}

func (Ciconia) DataFiles() ([]string, error) {
	return []string{ciconiaDataFile}, nil
}

func (Ciconia) Load() (map[language.Language][]dto.CiconiaQuote, error) {
	f, err := quote.DataFS.Open(ciconiaDataFile)
	if err != nil {
		return nil, fmt.Errorf("open ciconia script %s: %w", ciconiaDataFile, err)
	}
	defer f.Close()

	timeStart := time.Now()
	parsed, _, err := scriptparser.ParseReader(f, ciconia.NewParser())
	timeEnd := time.Now()

	if err != nil {
		return nil, fmt.Errorf("parse ciconia script %s: %w", ciconiaDataFile, err)
	}
	log.Printf("[ciconia/en] parsed %d quotes in %v", len(parsed), timeEnd.Sub(timeStart))

	indexed := make([]dto.CiconiaQuote, len(parsed))
	for i, q := range parsed {
		indexed[i] = dto.CiconiaQuote{
			CicroniaQuote: q,
			Index:         i,
		}
	}

	return map[language.Language][]dto.CiconiaQuote{
		language.English: indexed,
	}, nil
}

func (Ciconia) Base(q *dto.CiconiaQuote) *scriptdto.BaseQuote {
	return &q.BaseQuote
}

func (Ciconia) Row(q *dto.CiconiaQuote) QuoteRow {
	return QuoteRow{
		Index:      q.Index,
		Segment:    q.Chapter,
		TextJP:     q.TextJP,
		TextJPHtml: q.TextJPHtml,
	}
}

func (Ciconia) Build(base scriptdto.BaseQuote, r QuoteRow) dto.CiconiaQuote {
	return dto.CiconiaQuote{
		BaseQuote:  base,
		TextJP:     r.TextJP,
		TextJPHtml: r.TextJPHtml,
		Chapter:    r.Segment,
		Index:      r.Index,
	}
}
