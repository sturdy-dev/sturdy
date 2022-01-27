package db

import (
	"context"

	"getsturdy.com/api/pkg/licenses/enterprise/cloud/validations"
)

type memory struct{}

func NewMemory() Repository {
	return &memory{}
}

func (m *memory) Create(context.Context, *validations.Validation) error {
	return nil
}
