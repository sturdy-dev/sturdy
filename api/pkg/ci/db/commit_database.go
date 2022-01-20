package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/ci"

	"github.com/jmoiron/sqlx"
)

func NewCommitRepository(db *sqlx.DB) CommitRepository {
	return &database{db: db}
}

type database struct {
	db *sqlx.DB
}

func (r *database) GetByCodebaseAndCiRepoCommitID(ctx context.Context, codebaseID, ciRepoCommitID string) (*ci.Commit, error) {
	var res ci.Commit
	err := r.db.GetContext(ctx, &res, `SELECT id, codebase_id, ci_repo_commit_id, trunk_commit_id, created_at FROM ci_commits WHERE codebase_id = $1 AND ci_repo_commit_id = $2`, codebaseID, ciRepoCommitID)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *database) Create(ctx context.Context, c *ci.Commit) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO ci_commits
		(id, codebase_id, ci_repo_commit_id, trunk_commit_id, created_at)
		VALUES(:id, :codebase_id, :ci_repo_commit_id, :trunk_commit_id, :created_at)`, &c)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}
