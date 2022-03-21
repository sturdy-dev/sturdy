package service

import (
	"context"

	"getsturdy.com/api/pkg/codebases"
)

type Service interface {
	Pull(ctx context.Context, codebaseID codebases.ID) error
}
