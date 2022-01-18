package db

import (
	"context"

	"mash/pkg/installations"
)

type Repository interface {
	Create(context.Context, *installations.Installation) error
	ListAll(context.Context) ([]*installations.Installation, error)
}
