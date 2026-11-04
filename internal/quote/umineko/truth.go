package umineko

import "umineko_quote/internal/quote/params"

type Truth string

const (
	TruthAll    Truth = ""
	TruthRed    Truth = "red"
	TruthBlue   Truth = "blue"
	TruthGold   Truth = "gold"
	TruthPurple Truth = "purple"
)

func (Truth) Parse(s string) Truth {
	switch s {
	case "red":
		return TruthRed
	case "blue":
		return TruthBlue
	case "gold":
		return TruthGold
	case "purple":
		return TruthPurple
	default:
		return TruthAll
	}
}

func (t Truth) Scope() params.QuoteScope {
	return params.QuoteScope{
		RedTruth:    t == TruthRed,
		BlueTruth:   t == TruthBlue,
		GoldTruth:   t == TruthGold,
		PurpleTruth: t == TruthPurple,
	}
}
