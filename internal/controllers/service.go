package controllers

import (
	"umineko_quote/internal/audio"
	"umineko_quote/internal/og"
	"umineko_quote/internal/quote/ciconia"
	"umineko_quote/internal/quote/higurashi"
	"umineko_quote/internal/quote/umineko"
)

type Service struct {
	UminekoService   umineko.Service
	HigurashiService higurashi.Service
	CiconiaService   ciconia.Service
	OGImageGenerator *og.ImageGenerator
	AudioCombiner    audio.Combiner
	HTMLContent      string
}

func NewService(
	uminekoService umineko.Service,
	higurashiService higurashi.Service,
	ciconiaService ciconia.Service,
	ogGen *og.ImageGenerator,
	audioCombiner audio.Combiner,
	htmlContent string,
) Service {
	return Service{
		UminekoService:   uminekoService,
		HigurashiService: higurashiService,
		CiconiaService:   ciconiaService,
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

func (s *Service) GetCiconiaRoutes() []FSetupRoute {
	var all []FSetupRoute
	all = append(all, s.getAllCiconiaQuoteRoutes()...)
	all = append(all, s.getCiconiaOGRoutes()...)
	return all
}

func (s *Service) GetSystemRoutes() []FSetupRoute {
	return s.getAllSystemRoutes()
}

func (s *Service) GetPageRoutes() []FSetupRoute {
	return s.getAllOGPageRoutes()
}
