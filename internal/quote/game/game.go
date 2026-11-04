package game

import (
	"umineko_quote/internal/quote/language"

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
		TextJP         string
		TextJPHtml     string
	}
)

const NarratorID = "narrator"
