package db

import (
	"context"

	"getsturdy.com/api/pkg/ci"
)

type CommitRepository interface {
	Create(context.Context, *ci.Commit) error
	GetByCodebaseAndCiRepoCommitID(ctx context.Context, codebaseID, ciRepoCommitID string) (*ci.Commit, error)
}
