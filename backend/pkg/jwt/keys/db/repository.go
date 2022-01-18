package db

import (
	"context"

	"mash/pkg/jwt/keys"
)

type Repository interface {
	Create(context.Context, *keys.Key) error
	Get(context.Context, string) (*keys.Key, error)
}
