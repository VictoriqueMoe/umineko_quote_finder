package sqlite

import (
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"
)

type uminekoDAO struct {
	*quoteDAO[dto.UminekoQuote]
}

func (d *uminekoDAO) TruthPerEpisode(lang language.Language) (map[int]params.TruthCounts, error) {
	rows, err := d.db.Query(`SELECT episode, SUM(has_red_truth), SUM(has_blue_truth), SUM(has_gold_truth), SUM(has_purple_truth)
		FROM quotes
		WHERE game = ? AND lang = ? AND (has_red_truth = 1 OR has_blue_truth = 1 OR has_gold_truth = 1 OR has_purple_truth = 1)
		GROUP BY episode`,
		d.game.Name(), string(lang),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	epTruth := make(map[int]params.TruthCounts)
	for rows.Next() {
		var ep int
		var counts params.TruthCounts
		if err := rows.Scan(&ep, &counts.Red, &counts.Blue, &counts.Gold, &counts.Purple); err != nil {
			return nil, err
		}
		epTruth[ep] = counts
	}

	return epTruth, rows.Err()
}
