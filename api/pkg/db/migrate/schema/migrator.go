package schema

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
)

//go:embed migrations/*.sql
var migrations embed.FS

type Service struct {
	migrator *migrate.Migrate
	db       *sqlx.DB
}

func New(db *sqlx.DB) (*Service, error) {
	m, err := getMigrator(db.DB)
	if err != nil {
		return nil, err
	}

	return &Service{
		migrator: m,
		db:       db,
	}, nil
}

// Up applies all of the migrations.
func (s *Service) Up() error {
	if err := s.migrator.Up(); errors.Is(err, migrate.ErrNoChange) {
		return nil
	} else {
		return err
	}
}

// UpTo applies all of the migrations up to the given version.
func (s *Service) UpTo(v uint) error {
	currentVersion, _, err := s.migrator.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		currentVersion = 0
	} else if err != nil {
		return fmt.Errorf("error getting current version: %w", err)
	}

	if currentVersion >= v {
		return nil
	}

	if err := s.migrator.Steps(int(v - currentVersion)); errors.Is(err, migrate.ErrNoChange) {
		return nil
	} else {
		return err
	}
}

func getMigrator(db *sql.DB) (*migrate.Migrate, error) {
	migrations, err := iofs.New(migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("error opening migrations: %w", err)
	}

	database, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("error creating database: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", migrations, "postgres", database)
	if err != nil {
		return nil, fmt.Errorf("error connecting to db to migrate: %w", err)
	}

	return m, nil
}
