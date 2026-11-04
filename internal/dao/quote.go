package dao

import (
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"
)

type (
	QuoteDAO[Q any] interface {
		FindQuotes(f params.QuoteFilter, limit int, offset int) ([]Q, int, error)
		RandomQuote(f params.QuoteFilter) (*Q, error)
		QuoteIndexByAudioID(lang language.Language, audioID string) (int, bool, error)
		QuoteByIndex(lang language.Language, index int) (*Q, error)
		QuotesBetween(lang language.Language, from int, to int) ([]Q, error)
		NearestVoiced(lang language.Language, index int, previous bool) (string, error)
		LanguageCounts() (map[language.Language]int, error)
		Speakers(lang language.Language) (map[string]string, error)
		SpeakerTallies(lang language.Language) (params.SpeakerTallies, error)
	}

	UminekoQuoteDAO interface {
		QuoteDAO[dto.UminekoQuote]
		TruthPerEpisode(lang language.Language) (map[int]params.TruthCounts, error)
	}

	DAOs struct {
		Umineko   UminekoQuoteDAO
		Higurashi QuoteDAO[dto.HigurashiQuote]
		Ciconia   QuoteDAO[dto.CiconiaQuote]
	}
)
