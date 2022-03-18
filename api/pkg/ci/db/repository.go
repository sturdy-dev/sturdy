package db

import (
	"context"

	"getsturdy.com/api/pkg/ci"
	"getsturdy.com/api/pkg/codebases"
)

type CommitRepository interface {
	Create(context.Context, *ci.Commit) error
	GetByCodebaseAndCiRepoCommitID(ctx context.Context, codebaseID codebases.ID, ciRepoCommitID string) (*ci.Commit, error)
}
