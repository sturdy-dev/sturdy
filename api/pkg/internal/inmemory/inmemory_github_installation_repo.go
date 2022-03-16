package inmemory

import (
	"database/sql"
	"getsturdy.com/api/pkg/github"
)

func NewInMemoryGitHubInstallationRepository() *inMemoryGithubInstallationRepo {
	return &inMemoryGithubInstallationRepo{
		installs: make([]github.Installation, 0),
	}
}

type inMemoryGithubInstallationRepo struct {
	installs []github.Installation
}

func (i *inMemoryGithubInstallationRepo) GetByID(ID string) (*github.Installation, error) {
	panic("implement me")
}

func (i *inMemoryGithubInstallationRepo) GetByOwner(owner string) (*github.Installation, error) {
	panic("implement me")
}

func (i *inMemoryGithubInstallationRepo) GetByInstallationID(i2 int64) (*github.Installation, error) {
	for _, in := range i.installs {
		if in.InstallationID == i2 {
			return &in, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (i *inMemoryGithubInstallationRepo) Create(installation github.Installation) error {
	i.installs = append(i.installs, installation)
	return nil
}

func (i *inMemoryGithubInstallationRepo) Update(installation *github.Installation) error {
	for k, v := range i.installs {
		if v.ID == installation.ID {
			i.installs[k] = *installation
			return nil
		}
	}
	return sql.ErrNoRows
}
