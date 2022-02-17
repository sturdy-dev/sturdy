package global

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/installations"
	service_installations "getsturdy.com/api/pkg/installations/service"
)

func New(
	ctx context.Context,
	service *service_installations.Service,
) (*installations.Installation, error) {
	installation, err := service.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get installation: %w", err)
	}
	return installation, nil
}
