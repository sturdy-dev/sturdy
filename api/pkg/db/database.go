package db

import (
	"embed"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Setup(dbSourceURL string) (*sqlx.DB, error) {
	if err := doMigrateUp(dbSourceURL); err != nil {
		return nil, fmt.Errorf("error applying migrations: %w", err)
	}

	db, err := sqlx.Open("postgres", dbSourceURL)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func TrySetup(logger *zap.Logger, dbSourceURL string, timeout time.Duration) (*sqlx.DB, error) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			db, err := Setup(dbSourceURL)
			if err != nil {
				logger.Error("error connecting to db, will try again", zap.Error(err))
				continue
			}

			return db, nil
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for db to start")
		}
	}
}

func doMigrateUp(dbSourceURL string) error {
	migrations, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("error opening migrations: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", migrations, dbSourceURL)
	if err != nil {
		return fmt.Errorf("error connecting to db to migrate: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
