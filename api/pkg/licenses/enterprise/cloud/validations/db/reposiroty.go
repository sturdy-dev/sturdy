package db

import (
	"context"

	"getsturdy.com/api/pkg/licenses/enterprise/cloud/validations"
)

type Repository interface {
	Create(context.Context, *validations.Validation) error
}
