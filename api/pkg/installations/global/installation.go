package global

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"getsturdy.com/api/pkg/installations"
	service_installations "getsturdy.com/api/pkg/installations/service"
)

func New(
	ctx context.Context,
	service *service_installations.Service,
) (*installations.Installation, error) {
	installation, err := service.Get(ctx)
	// The installations table does not exist yet
	// This happens during the first time setup, when an installation might be accessed before the migration to create the table has been created
	// In this scenario, we return a temporary fake installation
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code.Name() == "undefined_table" {
		return &installations.Installation{ID: "tmp-first-time-setup"}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get installation: %w %T", err, err)
	}
	return installation, nil
}
