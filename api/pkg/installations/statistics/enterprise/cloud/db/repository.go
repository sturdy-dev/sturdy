package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/installations/statistics"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

type Repository interface {
	Create(context.Context, *statistics.Statistic) error
	GetByLicenseKey(context.Context, string) (*statistics.Statistic, error)
}
