package db

import (
	"context"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/integrations"
)

type IntegrationsRepository interface {
	Create(context.Context, *integrations.Integration) error
	Update(context.Context, *integrations.Integration) error
	ListByCodebaseID(context.Context, codebases.ID) ([]*integrations.Integration, error)
	Get(ctx context.Context, id string) (*integrations.Integration, error)
}
