package db

import (
	"context"
	"getsturdy.com/api/pkg/integrations"
)

type IntegrationsRepository interface {
	Create(context.Context, *integrations.Integration) error
	Update(context.Context, *integrations.Integration) error
	ListByCodebaseID(context.Context, string) ([]*integrations.Integration, error)
	Get(ctx context.Context, id string) (*integrations.Integration, error)
}
