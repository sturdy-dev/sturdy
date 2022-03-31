package db

import (
	"fmt"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/github"

	"github.com/jmoiron/sqlx"
)

type GitHubRepositoryRepository interface {
	GetByInstallationAndGitHubRepoID(installationID, gitHubRepositoryID int64) (*github.Repository, error)
	GetByInstallationAndName(installationID int64, name string) (*github.Repository, error)
	GetByCodebaseID(codebases.ID) (*github.Repository, error)
	GetByID(ID string) (*github.Repository, error)
	ListByInstallationID(installationID int64) ([]*github.Repository, error)
	ListByInstallationIDAndGitHubRepoIDs(installationID int64, gitHubRepositoryIDs []int64) ([]*github.Repository, error)
	Create(repository github.Repository) error
	Update(*github.Repository) error
}

type gitHubRepositoryRepo struct {
	db *sqlx.DB
}

func NewGitHubRepositoryRepository(db *sqlx.DB) GitHubRepositoryRepository {
	return &gitHubRepositoryRepo{db}
}

func (r *gitHubRepositoryRepo) GetByInstallationAndGitHubRepoID(installationID, gitHubRepositoryID int64) (*github.Repository, error) {
	var res github.Repository
	err := r.db.Get(&res, `SELECT *
		 FROM github_repositories
         WHERE installation_id = $1
           AND github_repository_id = $2
           AND deleted_at IS NULL`, installationID, gitHubRepositoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *gitHubRepositoryRepo) GetByInstallationAndName(installationID int64, name string) (*github.Repository, error) {
	var res github.Repository
	err := r.db.Get(&res, `SELECT *
		FROM github_repositories
		WHERE installation_id = $1
		  AND name = $2
		  AND deleted_at IS NULL`, installationID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *gitHubRepositoryRepo) GetByID(ID string) (*github.Repository, error) {
	var res github.Repository
	err := r.db.Get(&res, `SELECT * FROM github_repositories WHERE id = $1 AND deleted_at IS NULL`, ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *gitHubRepositoryRepo) GetByCodebaseID(codebaseID codebases.ID) (*github.Repository, error) {
	var res github.Repository
	err := r.db.Get(&res, `SELECT * FROM github_repositories WHERE codebase_id = $1 AND deleted_at IS NULL`, codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *gitHubRepositoryRepo) ListByInstallationID(installationID int64) ([]*github.Repository, error) {
	var entities []*github.Repository
	err := r.db.Select(&entities, `SELECT * FROM github_repositories WHERE installation_id = $1 AND deleted_at IS NULL`, installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entities, nil
}

func (r *gitHubRepositoryRepo) ListByInstallationIDAndGitHubRepoIDs(installationID int64, gitHubRepositoryIDs []int64) ([]*github.Repository, error) {
	query, args, err := sqlx.In(`SELECT * FROM github_repositories WHERE installation_id = ? AND github_repository_id IN(?) AND deleted_at IS NULL`,
		installationID,
		gitHubRepositoryIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}
	query = r.db.Rebind(query)
	var entities []*github.Repository
	err = r.db.Select(&entities, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entities, nil
}

func (r *gitHubRepositoryRepo) Create(i github.Repository) error {
	_, err := r.db.NamedExec(`INSERT INTO github_repositories (id, installation_id, name, created_at, github_repository_id, codebase_id, tracked_branch, synced_at, installation_access_token, installation_access_token_expires_at)
		VALUES (:id, :installation_id, :name, :created_at, :github_repository_id, :codebase_id, :tracked_branch, :synced_at, :installation_access_token, :installation_access_token_expires_at)`, &i)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *gitHubRepositoryRepo) Update(i *github.Repository) error {
	_, err := r.db.NamedExec(`UPDATE github_repositories
			SET uninstalled_at = :uninstalled_at,
				installation_access_token = :installation_access_token,
				installation_access_token_expires_at = :installation_access_token_expires_at,
				tracked_branch = :tracked_branch,
				synced_at = :synced_at,
			    integration_enabled = :integration_enabled,
			    github_source_of_truth = :github_source_of_truth,
			    last_push_at = :last_push_at,
			    last_push_error_message = :last_push_error_message,
			    deleted_at = :deleted_at
			WHERE id = :id`, i)
	if err != nil {
		return fmt.Errorf("failed to update repo: %w", err)
	}
	return nil
}
