package migrate

import (
	"context"
	"fmt"
	"sort"

	"getsturdy.com/api/pkg/db/migrate/data"
	"getsturdy.com/api/pkg/db/migrate/schema"
	"go.uber.org/zap"
)

type Service struct {
	datamigrations   *data.Service
	schemamigrations *schema.Service
	logger           *zap.Logger
}

func New(
	datamigrations *data.Service,
	schemamigrations *schema.Service,
	logger *zap.Logger,
) *Service {
	return &Service{
		logger:           logger.Named("migrate"),
		datamigrations:   datamigrations,
		schemamigrations: schemamigrations,
	}
}

func (s *Service) Migrate(ctx context.Context) error {
	migrateDataOnVersion := s.datamigrations.Versions()
	if noDataMigrations := len(migrateDataOnVersion) == 0; noDataMigrations {
		s.logger.Warn("applying db schema migrations")
		return s.schemamigrations.Up()
	}

	dataMigrationVersions := []uint{}
	for v := range migrateDataOnVersion {
		dataMigrationVersions = append(dataMigrationVersions, v)
	}
	sort.Slice(dataMigrationVersions, func(i, j int) bool {
		return dataMigrationVersions[i] < dataMigrationVersions[j]
	})

	for _, version := range dataMigrationVersions {
		s.logger.Warn("applying db schema migrations", zap.Uint("version", version-1))
		if err := s.schemamigrations.UpTo(version - 1); err != nil {
			return fmt.Errorf("failed to migrate schema to version %d: %w", version, err)
		}

		s.logger.Warn("applying db data migrations", zap.Uint("version", version))
		if err := s.datamigrations.Run(ctx, version); err != nil {
			return fmt.Errorf("failed to migrate data to version %d: %w", version, err)
		}

		s.logger.Warn("applying db schema migrations", zap.Uint("version", version))
		if err := s.schemamigrations.UpTo(version); err != nil {
			return fmt.Errorf("failed to migrate schema to version %d: %w", version, err)
		}
	}

	s.logger.Warn("applying db schema migrations")
	return s.schemamigrations.Up()
}
