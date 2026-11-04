package game

import (
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"

	"github.com/VictoriqueMoe/umineko_script_parser/ciconia"
	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
)

type Ciconia struct{}

const ciconiaDataFile = "data/Ciconia/en.file"

var (
	ciconiaSources = map[language.Language]scriptSource[ciconia.ParsedQuote]{
		language.English:  {path: ciconiaDataFile, parser: ciconia.NewParser()},
		language.Japanese: {path: ciconiaDataFile, parser: ciconia.NewJapaneseParser()},
	}

	ciconiaLoader = scriptLoader[ciconia.ParsedQuote, dto.CiconiaQuote]{
		name: "ciconia",
		index: func(q ciconia.ParsedQuote, i int) dto.CiconiaQuote {
			return dto.CiconiaQuote{CicroniaQuote: q, Index: i}
		},
	}
)

func (Ciconia) Name() string {
	return "ciconia"
}

func (Ciconia) DataFiles() ([]string, error) {
	return sourceFiles(ciconiaSources), nil
}

func (Ciconia) Load() (map[language.Language][]dto.CiconiaQuote, error) {
	return loadLanguages(ciconiaSources, ciconiaLoader.load)
}

func (Ciconia) Base(q *dto.CiconiaQuote) *scriptdto.BaseQuote {
	return &q.BaseQuote
}

func (Ciconia) Row(q *dto.CiconiaQuote) QuoteRow {
	return QuoteRow{
		Index:   q.Index,
		Segment: q.Chapter,
	}
}

func (Ciconia) Build(base scriptdto.BaseQuote, r QuoteRow) dto.CiconiaQuote {
	return dto.CiconiaQuote{
		BaseQuote: base,
		Chapter:   r.Segment,
		Index:     r.Index,
	}
}
