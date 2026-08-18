package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
)

type Options struct {
	Driver       string
	SQLitePath   string
	PostgresDSN  string
	MaxOpenConns int
	MaxIdleConns int
}

func Open(opts Options) (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)
	switch opts.Driver {
	case "sqlite":
		if opts.SQLitePath != "" {
			if err := os.MkdirAll(filepath.Dir(opts.SQLitePath), 0o755); err != nil {
				return nil, fmt.Errorf("create sqlite directory: %w", err)
			}
		}
		db, err = sql.Open("sqlite", opts.SQLitePath)
	case "postgres":
		db, err = sql.Open("postgres", opts.PostgresDSN)
	default:
		return nil, fmt.Errorf("unsupported database driver %q", opts.Driver)
	}
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	if opts.MaxOpenConns > 0 {
		db.SetMaxOpenConns(opts.MaxOpenConns)
	}
	if opts.MaxIdleConns > 0 {
		db.SetMaxIdleConns(opts.MaxIdleConns)
	}
	db.SetConnMaxLifetime(30 * time.Minute)
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

func EnsureSQLiteSchema(db *sql.DB) error {
	_, err := db.Exec(schema)
	if err != nil {
		return fmt.Errorf("ensure sqlite schema: %w", err)
	}
	return nil
}
