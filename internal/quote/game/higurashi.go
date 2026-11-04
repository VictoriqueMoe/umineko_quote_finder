package game

import (
	"fmt"
	"log"
	"time"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote"
	"umineko_quote/internal/quote/language"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/higurashi"
)

type Higurashi struct{}

const higurashiDataFile = "data/higurashi/en.file"

func (Higurashi) Name() string {
	return "higurashi"
}

func (Higurashi) DataFiles() ([]string, error) {
	return []string{higurashiDataFile}, nil
}

func (Higurashi) Load() (map[language.Language][]dto.HigurashiQuote, error) {
	f, err := quote.DataFS.Open(higurashiDataFile)
	if err != nil {
		return nil, fmt.Errorf("open higurashi script %s: %w", higurashiDataFile, err)
	}
	defer f.Close()

	timeStart := time.Now()
	parsed, _, err := scriptparser.ParseReader(f, higurashi.NewParser())
	timeEnd := time.Now()

	if err != nil {
		return nil, fmt.Errorf("parse higurashi script %s: %w", higurashiDataFile, err)
	}
	log.Printf("[higurashi/en] parsed %d quotes in %v", len(parsed), timeEnd.Sub(timeStart))

	indexed := make([]dto.HigurashiQuote, len(parsed))
	for i, q := range parsed {
		indexed[i] = dto.HigurashiQuote{
			HigurashiQuote: q,
			Index:          i,
		}
	}

	return map[language.Language][]dto.HigurashiQuote{
		language.English: indexed,
	}, nil
}

func (Higurashi) Base(q *dto.HigurashiQuote) *scriptdto.BaseQuote {
	return &q.BaseQuote
}

func (Higurashi) Row(q *dto.HigurashiQuote) QuoteRow {
	return QuoteRow{
		Index:      q.Index,
		Segment:    q.Arc,
		TextJP:     q.TextJP,
		TextJPHtml: q.TextJPHtml,
	}
}

func (Higurashi) Build(base scriptdto.BaseQuote, r QuoteRow) dto.HigurashiQuote {
	return dto.HigurashiQuote{
		BaseQuote:  base,
		TextJP:     r.TextJP,
		TextJPHtml: r.TextJPHtml,
		Arc:        r.Segment,
		Index:      r.Index,
	}
}
