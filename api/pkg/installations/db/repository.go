package db

import (
	"context"

	"getsturdy.com/api/pkg/installations"
)

type Repository interface {
	Create(context.Context, *installations.Installation) error
	ListAll(context.Context) ([]*installations.Installation, error)
	Update(ctx context.Context, installation *installations.Installation) error
}
