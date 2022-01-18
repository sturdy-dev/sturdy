package global

import (
	"context"
	"fmt"

	"mash/pkg/installations"
	"mash/pkg/installations/db"

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
			ID:   uuid.New().String(),
			Type: installationType,
		}
		if err := repo.Create(ctx, installation); err != nil {
			return nil, fmt.Errorf("failed to create installation: %w", err)
		}
		logger.Info("installation", zap.String("id", installation.ID), zap.Stringer("type", installation.Type))
		return installation, nil
	case 1:
		installation := ii[0]
		installation.Type = installationType
		logger.Info("installation", zap.String("id", installation.ID), zap.Stringer("type", installation.Type))
		return installation, nil
	default:
		return nil, fmt.Errorf("more than one installation found")
	}
}
