package db

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/github"

	"github.com/jmoiron/sqlx"
)

type GitHubPRRepo interface {
	Create(pr github.PullRequest) error
	Get(id string) (*github.PullRequest, error)
	GetByGitHubIDAndCodebaseID(gitHubID int64, codebaseID codebases.ID) (*github.PullRequest, error)
	GetByCodebaseIDaAndHeadSHA(ctx context.Context, codebaseID codebases.ID, headSHA string) (*github.PullRequest, error)
	ListByHeadAndRepositoryID(head string, repositoryID int64) ([]*github.PullRequest, error)
	GetMostRecentlyClosedByWorkspace(workspaceID string) (*github.PullRequest, error)
	ListOpenedByWorkspace(workspaceID string) ([]*github.PullRequest, error)
	Update(context.Context, *github.PullRequest) error
}

type gitHubPRRepo struct {
	db *sqlx.DB
}

func NewGitHubPRRepo(db *sqlx.DB) GitHubPRRepo {
	return &gitHubPRRepo{db: db}
}

func (r *gitHubPRRepo) Create(pr github.PullRequest) error {
	_, err := r.db.NamedExec(`INSERT INTO github_pull_requests (
		id,
		workspace_id,
		github_id,
		github_repository_id,
		created_by,
		github_pr_number,
		head,
		head_sha,
		codebase_id,
		base,
		created_at,
		updated_at,
		closed_at,
		merged_at,
		state,
		open,
		merged
		) VALUES (
		:id,
		:workspace_id,
		:github_id,
		:github_repository_id,
		:created_by,
		:github_pr_number,
		:head,
		:head_sha,
		:codebase_id,
		:base,
		:created_at,
		:updated_at,
		:closed_at,
		:merged_at, 
		:state,
		:state = 'open',
		:state = 'merged'
	)`, &pr)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r gitHubPRRepo) GetByCodebaseIDaAndHeadSHA(ctx context.Context, codebaseID codebases.ID, headSHA string) (*github.PullRequest, error) {
	var pr github.PullRequest
	if err := r.db.GetContext(ctx, &pr, `
		SELECT
			*
		FROM
			github_pull_requests
		WHERE
			codebase_id = $1
			AND head_sha = $2
	`, codebaseID, headSHA); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return &pr, nil
}

func (r *gitHubPRRepo) Get(id string) (*github.PullRequest, error) {
	var pr github.PullRequest
	err := r.db.Get(&pr, `SELECT
		id, workspace_id, github_id, github_repository_id, created_by, github_pr_number, head, head_sha, codebase_id, base, created_at, updated_at, closed_at, merged_at, state 
		FROM github_pull_requests WHERE id=$1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &pr, nil
}

func (r *gitHubPRRepo) GetByGitHubIDAndCodebaseID(gitHubID int64, codebaseID codebases.ID) (*github.PullRequest, error) {
	var pr github.PullRequest
	err := r.db.Get(&pr, `SELECT
		id, workspace_id, github_id, github_repository_id, created_by, github_pr_number, head, head_sha, codebase_id, base, created_at, updated_at, closed_at, merged_at, state
		FROM github_pull_requests WHERE github_id=$1 AND codebase_id = $2 `, gitHubID, codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &pr, nil
}

func (r *gitHubPRRepo) ListOpenedByWorkspace(workspaceID string) ([]*github.PullRequest, error) {
	var entities []*github.PullRequest
	err := r.db.Select(&entities, `SELECT
		id, workspace_id, github_id, github_repository_id, created_by, github_pr_number, head, head_sha, codebase_id, base, created_at, updated_at, closed_at, merged_at, state
		FROM github_pull_requests WHERE workspace_id = $1 and open = true`, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entities, nil
}

func (r *gitHubPRRepo) GetMostRecentlyClosedByWorkspace(workspaceID string) (*github.PullRequest, error) {
	var pr github.PullRequest
	err := r.db.Get(&pr, `SELECT 
		id, workspace_id, github_id, github_repository_id, created_by, github_pr_number, head, head_sha, codebase_id, base, created_at, updated_at, closed_at, merged_at, state
		FROM github_pull_requests WHERE workspace_id=$1 ORDER BY closed_at DESC LIMIT 1`, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &pr, nil
}

func (r *gitHubPRRepo) ListByHeadAndRepositoryID(head string, repositoryID int64) ([]*github.PullRequest, error) {
	var entities []*github.PullRequest
	err := r.db.Select(&entities, `SELECT 
		id, workspace_id, github_id, github_repository_id, created_by, github_pr_number, head, head_sha, codebase_id, base, created_at, updated_at, closed_at, merged_at, state
		FROM github_pull_requests WHERE head = $1 AND github_repository_id = $2`, head, repositoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entities, nil
}

func (r *gitHubPRRepo) Update(ctx context.Context, pr *github.PullRequest) error {
	_, err := r.db.NamedExecContext(ctx, `UPDATE github_pull_requests
		SET updated_at = :updated_at,
		    closed_at = :closed_at,
		    merged_at = :merged_at,
			head_sha = :head_sha,
			state = :state,
			open = (:state = 'open'),
			merged = (:state = 'merged')
		WHERE id = :id
	`, pr)
	if err != nil {
		return fmt.Errorf("failed to update %w", err)
	}
	return nil
}
