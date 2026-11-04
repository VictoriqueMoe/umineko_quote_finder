package sqlite

import (
	"database/sql"

	"umineko_quote/internal/dao"
	"umineko_quote/internal/quote/game"
)

func New(db *sql.DB) (dao.DAOs, error) {
	umineko, err := newQuotes(db, game.Umineko{})
	if err != nil {
		return dao.DAOs{}, err
	}

	higurashi, err := newQuotes(db, game.Higurashi{})
	if err != nil {
		return dao.DAOs{}, err
	}

	ciconia, err := newQuotes(db, game.Ciconia{})
	if err != nil {
		return dao.DAOs{}, err
	}

	return dao.DAOs{
		Umineko:   &uminekoDAO{quoteDAO: umineko},
		Higurashi: higurashi,
		Ciconia:   ciconia,
	}, nil
}
