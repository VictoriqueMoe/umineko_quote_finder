package memory

import (
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"
)

type uminekoDAO struct {
	*quoteDAO[dto.UminekoQuote]
}

func (d *uminekoDAO) TruthPerEpisode(lang language.Language) (map[int]params.TruthCounts, error) {
	epTruth := make(map[int]params.TruthCounts)

	for _, q := range d.quotes[lang] {
		if !q.HasRedTruth && !q.HasBlueTruth && !q.HasGoldTruth && !q.HasPurpleTruth {
			continue
		}

		counts := epTruth[q.Episode]
		if q.HasRedTruth {
			counts.Red++
		}
		if q.HasBlueTruth {
			counts.Blue++
		}
		if q.HasGoldTruth {
			counts.Gold++
		}
		if q.HasPurpleTruth {
			counts.Purple++
		}
		epTruth[q.Episode] = counts
	}

	return epTruth, nil
}
