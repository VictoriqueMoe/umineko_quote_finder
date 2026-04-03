package controllers

import (
	"umineko_quote/internal/audio"
	"umineko_quote/internal/og"
	"umineko_quote/internal/quote/higurashi"
	"umineko_quote/internal/quote/umineko"
)

type Service struct {
	UminekoService   umineko.Service
	HigurashiService higurashi.Service
	OGImageGenerator *og.ImageGenerator
	AudioCombiner    audio.Combiner
	HTMLContent      string
}

func NewService(
	uminekoService umineko.Service,
	higurashiService higurashi.Service,
	ogGen *og.ImageGenerator,
	audioCombiner audio.Combiner,
	htmlContent string,
) Service {
	return Service{
		UminekoService:   uminekoService,
		HigurashiService: higurashiService,
		OGImageGenerator: ogGen,
		AudioCombiner:    audioCombiner,
		HTMLContent:      htmlContent,
	}
}

func (s *Service) GetUminekoRoutes() []FSetupRoute {
	var all []FSetupRoute
	all = append(all, s.getAllUminekoQuoteRoutes()...)
	all = append(all, s.getAllAudioRoutes()...)
	all = append(all, s.getUminekoOGRoutes()...)
	return all
}

func (s *Service) GetHigurashiRoutes() []FSetupRoute {
	var all []FSetupRoute
	all = append(all, s.getAllHigurashiQuoteRoutes()...)
	all = append(all, s.getHigurashiOGRoutes()...)
	return all
}

func (s *Service) GetSystemRoutes() []FSetupRoute {
	return s.getAllSystemRoutes()
}

func (s *Service) GetPageRoutes() []FSetupRoute {
	return s.getAllOGPageRoutes()
}
