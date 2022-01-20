package db

import (
	"fmt"
	"mash/pkg/github"

	"github.com/jmoiron/sqlx"
)

type GitHubRepositoryRepo interface {
	GetByInstallationAndGitHubRepoID(installationID, gitHubRepositoryID int64) (*github.GitHubRepository, error)
	GetByInstallationAndName(installationID int64, name string) (*github.GitHubRepository, error)
	GetByCodebaseID(repositoryID string) (*github.GitHubRepository, error)
	GetByID(ID string) (*github.GitHubRepository, error)
	ListByInstallationID(installationID int64) ([]*github.GitHubRepository, error)
	ListByInstallationIDAndGitHubRepoIDs(installationID int64, gitHubRepositoryIDs []int64) ([]*github.GitHubRepository, error)
	Create(repository github.GitHubRepository) error
	Update(*github.GitHubRepository) error
}

type gitHubRepositoryRepo struct {
	db *sqlx.DB
}

func NewGitHubRepositoryRepo(db *sqlx.DB) GitHubRepositoryRepo {
	return &gitHubRepositoryRepo{db}
}

func (r *gitHubRepositoryRepo) GetByInstallationAndGitHubRepoID(installationID, gitHubRepositoryID int64) (*github.GitHubRepository, error) {
	var res github.GitHubRepository
	err := r.db.Get(&res, "SELECT * FROM github_repositories WHERE installation_id = $1 AND github_repository_id = $2", installationID, gitHubRepositoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *gitHubRepositoryRepo) GetByInstallationAndName(installationID int64, name string) (*github.GitHubRepository, error) {
	var res github.GitHubRepository
	err := r.db.Get(&res, "SELECT * FROM github_repositories WHERE installation_id = $1 AND name = $2", installationID, name)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *gitHubRepositoryRepo) GetByID(ID string) (*github.GitHubRepository, error) {
	var res github.GitHubRepository
	err := r.db.Get(&res, "SELECT * FROM github_repositories WHERE id = $1", ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *gitHubRepositoryRepo) GetByCodebaseID(repositoryID string) (*github.GitHubRepository, error) {
	var res github.GitHubRepository
	err := r.db.Get(&res, "SELECT * FROM github_repositories WHERE codebase_id = $1", repositoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &res, nil
}

func (r *gitHubRepositoryRepo) ListByInstallationID(installationID int64) ([]*github.GitHubRepository, error) {
	var entities []*github.GitHubRepository
	err := r.db.Select(&entities, "SELECT * FROM github_repositories WHERE installation_id = $1", installationID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entities, nil
}

func (r *gitHubRepositoryRepo) ListByInstallationIDAndGitHubRepoIDs(installationID int64, gitHubRepositoryIDs []int64) ([]*github.GitHubRepository, error) {
	query, args, err := sqlx.In("SELECT * FROM github_repositories WHERE installation_id = ? AND github_repository_id IN(?)",
		installationID,
		gitHubRepositoryIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to create query: %w", err)
	}
	query = r.db.Rebind(query)
	var entities []*github.GitHubRepository
	err = r.db.Select(&entities, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entities, nil
}

func (r *gitHubRepositoryRepo) Create(i github.GitHubRepository) error {
	_, err := r.db.NamedExec(`INSERT INTO github_repositories (id, installation_id, name, created_at, github_repository_id, codebase_id, tracked_branch, synced_at, installation_access_token, installation_access_token_expires_at)
		VALUES (:id, :installation_id, :name, :created_at, :github_repository_id, :codebase_id, :tracked_branch, :synced_at, :installation_access_token, :installation_access_token_expires_at)`, &i)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *gitHubRepositoryRepo) Update(i *github.GitHubRepository) error {
	_, err := r.db.NamedExec(`UPDATE github_repositories
			SET uninstalled_at = :uninstalled_at,
				installation_access_token = :installation_access_token,
				installation_access_token_expires_at = :installation_access_token_expires_at,
				tracked_branch = :tracked_branch,
				synced_at = :synced_at,
			    integration_enabled = :integration_enabled,
			    github_source_of_truth = :github_source_of_truth,
			    last_push_at = :last_push_at,
			    last_push_error_message = :last_push_error_message
			WHERE id = :id`, i)
	if err != nil {
		return fmt.Errorf("failed to update repo: %w", err)
	}
	return nil
}
