package migrator

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"

	"getsturdy.com/api/pkg/datamigrations"
	"getsturdy.com/api/pkg/db"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Setup is only used by tests
func Setup(dbSourceURL string, datamigs datamigrations.Service) (*sqlx.DB, error) {
	db, err := db.SetupWithTimeout(dbSourceURL, time.Minute)
	if err != nil {
		return nil, err
	}
	if err := MigrateUP(db.DB, datamigs); err != nil {
		return nil, fmt.Errorf("error applying migrations: %w", err)
	}
	return db, nil
}

func MigrateUP(db *sql.DB, datamigs datamigrations.Service) error {
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

	// migrate one version at a time
	for {
		version, dirty, err := m.Version()
		log.Printf("Sturdy database schema version: %d, dirty: %t, err: %v", version, dirty, err)

		if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
			return fmt.Errorf("error getting db version: %w", err)
		}
		if dirty {
			return fmt.Errorf("database is dirty, aborting migrations")
		}

		// Trigger data migrations for this version of the database
		if err := datamigs.Run(context.Background(), version); err != nil {
			return fmt.Errorf("error running data migrations: %w", err)
		}

		err = m.Steps(1)
		if errors.Is(err, migrate.ErrNoChange) || errors.Is(err, os.ErrNotExist) {
			break
		}
		if err != nil {
			return err
		}
	}

	return nil
}
