package global

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/installations"
	"getsturdy.com/api/pkg/installations/db"
	"getsturdy.com/api/pkg/version"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func New(
	ctx context.Context,
	logger *zap.Logger,
	repo db.Repository,
) (*installations.Installation, error) {
	logger = logger.Named("installations")

	ii, err := repo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list installations: %w", err)
	}
	switch len(ii) {
	case 0:
		installation := &installations.Installation{
			ID:      uuid.New().String(),
			Type:    installationType,
			Version: version.Version,
		}
		if err := repo.Create(ctx, installation); err != nil {
			return nil, fmt.Errorf("failed to create installation: %w", err)
		}
		logger.Info(
			"installation",
			zap.String("id", installation.ID),
			zap.Stringer("type", installation.Type),
			zap.String("version", installation.Version),
		)
		return installation, nil
	case 1:
		installation := ii[0]
		installation.Type = installationType
		installation.Version = version.Version
		logger.Info(
			"installation",
			zap.String("id", installation.ID),
			zap.Stringer("type", installation.Type),
			zap.String("version", installation.Version),
		)
		return installation, nil
	default:
		return nil, fmt.Errorf("more than one installation found")
	}
}
