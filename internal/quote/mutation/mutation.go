package mutation

import (
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/mutation/engine"
)

type Engine interface {
	Apply(quotes []dto.ParsedQuote) []dto.ParsedQuote
}

type Pipeline struct {
	engines []Engine
}

func NewPipeline() *Pipeline {
	p := &Pipeline{
		engines: []Engine{
			&engine.KanonAttributionEngine{},
		},
	}
	return p
}

func (p *Pipeline) Apply(quotes []dto.ParsedQuote) []dto.ParsedQuote {
	for _, e := range p.engines {
		quotes = e.Apply(quotes)
	}
	return quotes
}
