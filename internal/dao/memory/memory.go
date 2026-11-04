package memory

import (
	"umineko_quote/internal/dao"
	"umineko_quote/internal/quote/game"
)

func New() (dao.DAOs, error) {
	umineko, err := game.Umineko{}.Load()
	if err != nil {
		return dao.DAOs{}, err
	}

	higurashi, err := game.Higurashi{}.Load()
	if err != nil {
		return dao.DAOs{}, err
	}

	ciconia, err := game.Ciconia{}.Load()
	if err != nil {
		return dao.DAOs{}, err
	}

	return dao.DAOs{
		Umineko:   &uminekoDAO{QuoteDAO: NewQuotes(game.Umineko{}, umineko)},
		Higurashi: NewQuotes(game.Higurashi{}, higurashi),
		Ciconia:   NewQuotes(game.Ciconia{}, ciconia),
	}, nil
}
