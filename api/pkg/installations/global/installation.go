package global

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/installations"
	service_installations "getsturdy.com/api/pkg/installations/service"

	"go.uber.org/zap"
)

func New(
	ctx context.Context,
	logger *zap.Logger,
	service *service_installations.Service,
) (*installations.Installation, error) {
	installation, err := service.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get installation: %w", err)
	}

	logger.Info(
		"installation",
		zap.String("id", installation.ID),
		zap.Stringer("type", installation.Type),
		zap.String("version", installation.Version),
	)
	return installation, nil
}
