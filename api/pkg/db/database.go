package db

import (
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Setup(dbSourceURL string) (*sqlx.DB, error) {
	db, err := setup(dbSourceURL)
	if err != nil {
		return nil, err
	}
	if err := doMigrateUp(db.DB); err != nil {
		return nil, fmt.Errorf("error applying migrations: %w", err)
	}
	return db, nil
}

func setup(dbSourceURL string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dbSourceURL)
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func SetupWithTimeout(dbSourceURL string, timeout time.Duration) (*sqlx.DB, error) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			db, err := setup(dbSourceURL)
			if err != nil {
				continue
			}
			if err := doMigrateUp(db.DB); err != nil {
				return nil, fmt.Errorf("error applying migrations: %w", err)
			}

			return db, nil
		case <-time.After(timeout):
			return nil, fmt.Errorf("timeout waiting for db to start")
		}
	}
}

func doMigrateUp(db *sql.DB) error {
	migrations, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("error opening migrations: %w", err)
	}

	database, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("error creating database: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", migrations, "postgres", database)
	if err != nil {
		return fmt.Errorf("error connecting to db to migrate: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
