package db

import (
	"context"

	"getsturdy.com/api/pkg/installations/statistics"
)

type memory struct{}

func NewMemory() Repository {
	return &memory{}
}

func (m *memory) Create(context.Context, *statistics.Statistic) error {
	return nil
}
