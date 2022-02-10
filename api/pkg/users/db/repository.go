package db

import (
	"context"

	"getsturdy.com/api/pkg/users"
)

type Repository interface {
	Create(newUser *users.User) error
	Get(id string) (*users.User, error)
	GetByIDs(ctx context.Context, ids ...string) ([]*users.User, error)
	GetByEmail(email string) (*users.User, error)
	Update(*users.User) error
	UpdatePassword(u *users.User) error
	Count(context.Context) (uint64, error)
	List(ctx context.Context, limit uint64) ([]*users.User, error)
}
