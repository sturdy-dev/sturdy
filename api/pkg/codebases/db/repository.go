package db

import (
	"context"

	"getsturdy.com/api/pkg/codebases"
)

type CodebaseRepository interface {
	Create(codebases.Codebase) error
	Get(id string) (*codebases.Codebase, error)
	GetAllowArchived(id string) (*codebases.Codebase, error)
	GetByInviteCode(inviteCode string) (*codebases.Codebase, error)
	GetByShortID(shortID string) (*codebases.Codebase, error)
	Update(entity *codebases.Codebase) error
	ListByOrganization(ctx context.Context, organizationID string) ([]*codebases.Codebase, error)
	Count(context.Context) (uint64, error)
}
