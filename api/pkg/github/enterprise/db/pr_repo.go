package db

import (
	"context"
	"fmt"
	"getsturdy.com/api/pkg/github"

	"github.com/jmoiron/sqlx"
)

type gitHubPRRepo struct {
	db *sqlx.DB
}

type GitHubPRRepo interface {
	Create(pr github.PullRequest) error
	Get(ID string) (*github.PullRequest, error)
	GetByGitHubID(gitHubID int64) (*github.PullRequest, error)
	GetByCodebaseIDaAndHeadSHA(ctx context.Context, codebaseID, headSHA string) (*github.PullRequest, error)
	ListByHeadAndRepositoryID(head string, repositoryID int64) ([]*github.PullRequest, error)
	GetMostRecentlyClosedByWorkspace(workspaceID string) (*github.PullRequest, error)
	ListOpenedByWorkspace(workspaceID string) ([]*github.PullRequest, error)
	Update(pr *github.PullRequest) error
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
		open,
		merged,
		created_at,
		updated_at,
		closed_at,
		merged_at)
			VALUES (
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
		:open,
		:merged,
		:created_at,
		:updated_at,
		:closed_at,
		:merged_at)`, &pr)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r gitHubPRRepo) GetByCodebaseIDaAndHeadSHA(ctx context.Context, codebaseID, headSHA string) (*github.PullRequest, error) {
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

func (r *gitHubPRRepo) Get(ID string) (*github.PullRequest, error) {
	var pr github.PullRequest
	err := r.db.Get(&pr, "SELECT * FROM github_pull_requests WHERE id=$1", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &pr, nil
}

func (r *gitHubPRRepo) GetByGitHubID(gitHubID int64) (*github.PullRequest, error) {
	var pr github.PullRequest
	err := r.db.Get(&pr, "SELECT * FROM github_pull_requests WHERE github_id=$1", gitHubID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &pr, nil
}

func (r *gitHubPRRepo) ListOpenedByWorkspace(workspaceID string) ([]*github.PullRequest, error) {
	var entities []*github.PullRequest
	err := r.db.Select(&entities, "SELECT * FROM github_pull_requests WHERE workspace_id = $1 and open = true", workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entities, nil
}

func (r *gitHubPRRepo) GetMostRecentlyClosedByWorkspace(workspaceID string) (*github.PullRequest, error) {
	var pr github.PullRequest
	err := r.db.Get(&pr, "SELECT * FROM github_pull_requests WHERE workspace_id=$1 ORDER BY closed_at DESC LIMIT 1", workspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &pr, nil
}

func (r *gitHubPRRepo) ListByHeadAndRepositoryID(head string, repositoryID int64) ([]*github.PullRequest, error) {
	var entities []*github.PullRequest
	err := r.db.Select(&entities, "SELECT * FROM github_pull_requests WHERE head = $1 AND github_repository_id = $2", head, repositoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entities, nil
}

func (r *gitHubPRRepo) Update(pr *github.PullRequest) error {
	_, err := r.db.NamedExec(`UPDATE github_pull_requests
		SET open = :open,
		    merged = :merged,
		    updated_at = :updated_at,
		    closed_at = :closed_at,
		    merged_at = :merged_at,
			head_sha = :head_sha
		WHERE id=:id`, pr)
	if err != nil {
		return fmt.Errorf("failed to update %w", err)
	}
	return nil
}
