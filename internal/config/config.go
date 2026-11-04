package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type (
	StoreMode string

	Config struct {
		StoreMode  StoreMode
		SQLitePath string
	}
)

const (
	StoreModeMemory StoreMode = "memory"
	StoreModeSQLite StoreMode = "sqlite"
)

var Cfg Config

func Load() error {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return fmt.Errorf("load .env: %w", err)
		}
	}

	mode := StoreMode(envOr("STORE_MODE", string(StoreModeMemory)))
	if mode != StoreModeMemory && mode != StoreModeSQLite {
		return fmt.Errorf("STORE_MODE must be %q or %q, got %q", StoreModeMemory, StoreModeSQLite, mode)
	}

	Cfg = Config{
		StoreMode:  mode,
		SQLitePath: envOr("SQLITE_PATH", "./quotes.db"),
	}

	return nil
}

func envOr(key string, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
