package inmemory

import (
	"database/sql"
	"getsturdy.com/api/pkg/github"
)

func NewInMemoryGitHubRepositoryRepo() *inMemoryGithubRepositoryRepo {
	return &inMemoryGithubRepositoryRepo{
		repos: make([]github.GitHubRepository, 0),
	}
}

type inMemoryGithubRepositoryRepo struct {
	repos []github.GitHubRepository
}

func (i *inMemoryGithubRepositoryRepo) GetByInstallationAndGitHubRepoID(installationID, gitHubRepositoryID int64) (*github.GitHubRepository, error) {
	for _, r := range i.repos {
		if r.InstallationID == installationID && r.GitHubRepositoryID == gitHubRepositoryID {
			return &r, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (i *inMemoryGithubRepositoryRepo) GetByInstallationAndName(installationID int64, name string) (*github.GitHubRepository, error) {
	panic("implement me")
}

func (i *inMemoryGithubRepositoryRepo) GetByCodebaseID(codebaseID string) (*github.GitHubRepository, error) {
	for _, r := range i.repos {
		if r.CodebaseID == codebaseID {
			return &r, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (i *inMemoryGithubRepositoryRepo) GetByID(ID string) (*github.GitHubRepository, error) {
	for _, r := range i.repos {
		if r.ID == ID {
			return &r, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (i *inMemoryGithubRepositoryRepo) ListByInstallationID(installationID int64) ([]*github.GitHubRepository, error) {
	panic("implement me")
}

func (i *inMemoryGithubRepositoryRepo) ListByInstallationIDAndGitHubRepoIDs(installationID int64, gitHubRepositoryIDs []int64) ([]*github.GitHubRepository, error) {
	panic("implement me")
}

func (i *inMemoryGithubRepositoryRepo) Create(repository github.GitHubRepository) error {
	i.repos = append(i.repos, repository)
	return nil
}

func (i *inMemoryGithubRepositoryRepo) Update(repository *github.GitHubRepository) error {
	for idx, r := range i.repos {
		if r.ID == repository.ID {
			i.repos[idx] = *repository
			break
		}
	}
	return nil
}

func (i *inMemoryGithubRepositoryRepo) GetUnsyncedRepositories() ([]*github.GitHubRepository, error) {
	panic("implement me")
}
