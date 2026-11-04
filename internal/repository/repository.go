package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"umineko_quote/internal/config"
	"umineko_quote/internal/dao"
	"umineko_quote/internal/dao/memory"
	"umineko_quote/internal/dao/sqlite"
	"umineko_quote/internal/db"
	"umineko_quote/internal/dto"
)

type (
	QuoteRepository[Q any] interface {
		dao.QuoteDAO[Q]
	}

	UminekoRepository interface {
		dao.UminekoQuoteDAO
	}

	Repositories struct {
		db        *sql.DB
		Umineko   UminekoRepository
		Higurashi QuoteRepository[dto.HigurashiQuote]
		Ciconia   QuoteRepository[dto.CiconiaQuote]
	}

	quoteRepository[Q any] struct {
		dao.QuoteDAO[Q]
	}

	uminekoRepository struct {
		dao.UminekoQuoteDAO
	}
)

func New() (*Repositories, error) {
	cfg := config.Cfg

	start := time.Now()

	if cfg.StoreMode == config.StoreModeMemory {
		daos, err := memory.New()
		if err != nil {
			return nil, fmt.Errorf("initialise memory daos: %w", err)
		}

		log.Printf("[repository] loaded in-memory quotes in %v", time.Since(start).Round(time.Millisecond))
		return NewFromDAOs(daos, nil), nil
	}

	conn, err := db.Open(cfg.SQLitePath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err := db.Migrate(conn); err != nil {
		return nil, errors.Join(fmt.Errorf("migrate sqlite database: %w", err), conn.Close())
	}

	daos, err := sqlite.New(conn)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("initialise sqlite daos: %w", err), conn.Close())
	}

	log.Printf("[repository] opened sqlite database %s in %v", cfg.SQLitePath, time.Since(start).Round(time.Millisecond))
	return NewFromDAOs(daos, conn), nil
}

func NewFromDAOs(daos dao.DAOs, conn *sql.DB) *Repositories {
	return &Repositories{
		db:        conn,
		Umineko:   &uminekoRepository{UminekoQuoteDAO: daos.Umineko},
		Higurashi: &quoteRepository[dto.HigurashiQuote]{QuoteDAO: daos.Higurashi},
		Ciconia:   &quoteRepository[dto.CiconiaQuote]{QuoteDAO: daos.Ciconia},
	}
}

func (r *Repositories) Close() error {
	if r.db == nil {
		return nil
	}
	return r.db.Close()
}
