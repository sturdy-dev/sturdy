package inmemory

import (
	"database/sql"
	"getsturdy.com/api/pkg/github"
)

func NewInMemoryGitHubInstallationRepository() *inMemoryGithubInstallationRepo {
	return &inMemoryGithubInstallationRepo{
		installs: make([]github.GitHubInstallation, 0),
	}
}

type inMemoryGithubInstallationRepo struct {
	installs []github.GitHubInstallation
}

func (i *inMemoryGithubInstallationRepo) GetByID(ID string) (*github.GitHubInstallation, error) {
	panic("implement me")
}

func (i *inMemoryGithubInstallationRepo) GetByOwner(owner string) (*github.GitHubInstallation, error) {
	panic("implement me")
}

func (i *inMemoryGithubInstallationRepo) GetByInstallationID(i2 int64) (*github.GitHubInstallation, error) {
	for _, in := range i.installs {
		if in.InstallationID == i2 {
			return &in, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (i *inMemoryGithubInstallationRepo) Create(installation github.GitHubInstallation) error {
	i.installs = append(i.installs, installation)
	return nil
}

func (i *inMemoryGithubInstallationRepo) Update(installation *github.GitHubInstallation) error {
	for k, v := range i.installs {
		if v.ID == installation.ID {
			i.installs[k] = *installation
			return nil
		}
	}
	return sql.ErrNoRows
}
