package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/licenses"
)

type Repository interface {
	Get(context.Context, licenses.ID) (*licenses.License, error)
	Create(context.Context, *licenses.License) error
	ListByOrganizationID(context.Context, string) ([]*licenses.License, error)
	GetByKey(context.Context, string) (*licenses.License, error)
}

var (
	ErrNotFound = fmt.Errorf("not found")
)
