package db

import (
	"context"

	"getsturdy.com/api/pkg/licenses"
	"getsturdy.com/api/pkg/licenses/enterprise/cloud/validations"
)

type Repository interface {
	Create(context.Context, *validations.Validation) error
	ListLatest(context.Context, licenses.ID) ([]*validations.Validation, error)
}
