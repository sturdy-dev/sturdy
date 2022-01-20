package db

import (
	"context"

	"mash/pkg/onetime"
)

type Repository interface {
	Create(context.Context, *onetime.Token) error
	Update(context.Context, *onetime.Token) error
	Get(ctx context.Context, userID, key string) (*onetime.Token, error)
}
