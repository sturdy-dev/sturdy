package db

import (
	"context"

	"getsturdy.com/api/pkg/codebases"
)

type CodebaseRepository interface {
	Create(codebases.Codebase) error
	Get(codebases.ID) (*codebases.Codebase, error)
	GetAllowArchived(codebases.ID) (*codebases.Codebase, error)
	GetByInviteCode(inviteCode string) (*codebases.Codebase, error)
	GetByShortID(codebases.ShortCodebaseID) (*codebases.Codebase, error)
	Update(entity *codebases.Codebase) error
	ListByOrganization(ctx context.Context, organizationID string) ([]*codebases.Codebase, error)
	Count(context.Context) (uint64, error)
}
