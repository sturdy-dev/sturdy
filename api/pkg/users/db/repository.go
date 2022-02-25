package db

import (
	"context"

	"getsturdy.com/api/pkg/users"
)

type Repository interface {
	Create(*users.User) error
	Get(users.ID) (*users.User, error)
	GetByIDs(context.Context, ...users.ID) ([]*users.User, error)
	GetByEmail(string) (*users.User, error)
	Update(*users.User) error
	UpdatePassword(*users.User) error
	Count(context.Context) (uint64, error)
	List(ctx context.Context, limit uint64) ([]*users.User, error)
}
