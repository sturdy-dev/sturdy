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
)

//go:embed migrations/*.sql
var migrations embed.FS

func Setup(dbSourceURL string) (*sqlx.DB, error) {
	if err := doMigrateUp(dbSourceURL); err != nil {
		return nil, fmt.Errorf("error applying migrations: %w", err)
	}

	db, err := sqlx.Open("postgres", dbSourceURL)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
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
