package db

import (
	"context"

	"getsturdy.com/api/pkg/installations/statistics"
)

type Repository interface {
	Create(context.Context, *statistics.Statistic) error
}
