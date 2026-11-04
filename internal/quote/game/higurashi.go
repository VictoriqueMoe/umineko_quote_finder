package game

import (
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/higurashi"
)

type Higurashi struct{}

const higurashiDataFile = "data/higurashi/en.file"

var (
	higurashiSources = map[language.Language]scriptSource[higurashi.ParsedQuote]{
		language.English:  {path: higurashiDataFile, parser: higurashi.NewParser()},
		language.Japanese: {path: higurashiDataFile, parser: higurashi.NewJapaneseParser()},
	}

	higurashiLoader = scriptLoader[higurashi.ParsedQuote, dto.HigurashiQuote]{
		name: "higurashi",
		index: func(q higurashi.ParsedQuote, i int) dto.HigurashiQuote {
			return dto.HigurashiQuote{HigurashiQuote: q, Index: i}
		},
	}
)

func (Higurashi) Name() string {
	return "higurashi"
}

func (Higurashi) DataFiles() ([]string, error) {
	return sourceFiles(higurashiSources), nil
}

func (Higurashi) Load() (map[language.Language][]dto.HigurashiQuote, error) {
	return loadLanguages(higurashiSources, higurashiLoader.load)
}

func (Higurashi) Base(q *dto.HigurashiQuote) *scriptdto.BaseQuote {
	return &q.BaseQuote
}

func (Higurashi) Row(q *dto.HigurashiQuote) QuoteRow {
	return QuoteRow{
		Index:   q.Index,
		Segment: q.Arc,
	}
}

func (Higurashi) Build(base scriptdto.BaseQuote, r QuoteRow) dto.HigurashiQuote {
	return dto.HigurashiQuote{
		BaseQuote: base,
		Arc:       r.Segment,
		Index:     r.Index,
	}
}
