package inmemory

import (
	"database/sql"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/github"
)

func NewInMemoryGitHubRepositoryRepo() *inMemoryGithubRepositoryRepo {
	return &inMemoryGithubRepositoryRepo{
		repos: make([]github.Repository, 0),
	}
}

type inMemoryGithubRepositoryRepo struct {
	repos []github.Repository
}

func (i *inMemoryGithubRepositoryRepo) GetByInstallationAndGitHubRepoID(installationID, gitHubRepositoryID int64) (*github.Repository, error) {
	for _, r := range i.repos {
		if r.InstallationID == installationID && r.GitHubRepositoryID == gitHubRepositoryID && r.DeletedAt == nil {
			return &r, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (i *inMemoryGithubRepositoryRepo) GetByInstallationAndName(installationID int64, name string) (*github.Repository, error) {
	panic("implement me")
}

func (i *inMemoryGithubRepositoryRepo) GetByCodebaseID(codebaseID codebases.ID) (*github.Repository, error) {
	for _, r := range i.repos {
		if r.CodebaseID == codebaseID && r.DeletedAt == nil {
			return &r, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (i *inMemoryGithubRepositoryRepo) GetByID(ID string) (*github.Repository, error) {
	for _, r := range i.repos {
		if r.ID == ID && r.DeletedAt == nil {
			return &r, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (i *inMemoryGithubRepositoryRepo) ListByInstallationID(installationID int64) ([]*github.Repository, error) {
	panic("implement me")
}

func (i *inMemoryGithubRepositoryRepo) ListByInstallationIDAndGitHubRepoIDs(installationID int64, gitHubRepositoryIDs []int64) ([]*github.Repository, error) {
	panic("implement me")
}

func (i *inMemoryGithubRepositoryRepo) Create(repository github.Repository) error {
	i.repos = append(i.repos, repository)
	return nil
}

func (i *inMemoryGithubRepositoryRepo) Update(repository *github.Repository) error {
	for idx, r := range i.repos {
		if r.ID == repository.ID {
			i.repos[idx] = *repository
			break
		}
	}
	return nil
}

func (i *inMemoryGithubRepositoryRepo) GetUnsyncedRepositories() ([]*github.Repository, error) {
	panic("implement me")
}
