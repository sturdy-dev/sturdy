package db

import (
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Setup(dbSourceURL string, migrateUp bool, migrationsPath string) (*sqlx.DB, error) {
	if migrateUp {
		err := doMigrateUp(dbSourceURL, migrationsPath)
		if err != nil {
			return nil, err
		}
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

func doMigrateUp(dbSourceURL, migrationsPath string) error {
	m, err := migrate.New(
		migrationsPath,
		dbSourceURL)
	if err != nil {
		return fmt.Errorf("error connecting to db to migrate: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("error running db migrations: %w", err)
	}
	return nil
}
