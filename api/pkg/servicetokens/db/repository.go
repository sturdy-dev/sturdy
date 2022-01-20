package db

import (
	"context"

	"getsturdy.com/api/pkg/servicetokens"
)

type Repository interface {
	Create(context.Context, *servicetokens.Token) error
	GetByID(context.Context, string) (*servicetokens.Token, error)
}
