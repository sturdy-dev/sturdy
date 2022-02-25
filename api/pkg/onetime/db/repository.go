package db

import (
	"context"

	"getsturdy.com/api/pkg/onetime"
	"getsturdy.com/api/pkg/users"
)

type Repository interface {
	Create(context.Context, *onetime.Token) error
	Update(context.Context, *onetime.Token) error
	Get(ctx context.Context, userID users.ID, key string) (*onetime.Token, error)
}
