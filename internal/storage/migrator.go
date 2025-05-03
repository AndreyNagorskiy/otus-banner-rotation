package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"
	"github.com/pressly/goose"
)

var migrationsPath = "internal/storage/migrations"

func Migrate(dsn string, retryWithoutSSL bool) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open connection to database: %w", err)
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(db, migrationsPath); err != nil {
		if errors.Is(err, pq.ErrSSLNotSupported) && !retryWithoutSSL {
			err := Migrate(dsn+"?sslmode=disable", true)
			if err != nil {
				return fmt.Errorf("failed to run migrations without SSL: %w", err)
			}
		} else {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
	}
	return nil
}
