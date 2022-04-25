package db

import (
	"database/sql"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/users"
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

func NewInMemoryGitHubUserRepo() *inMemoryGitHubUserRepo {
	return &inMemoryGitHubUserRepo{
		users: make([]github.User, 0),
	}
}

type inMemoryGitHubUserRepo struct {
	users []github.User
}

func (i *inMemoryGitHubUserRepo) Create(user github.User) error {
	i.users = append(i.users, user)
	return nil
}

func (i *inMemoryGitHubUserRepo) GetByUsername(username string) (*github.User, error) {
	for _, v := range i.users {
		if v.Username == username {
			return &v, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (i *inMemoryGitHubUserRepo) GetByUserID(userID users.ID) (*github.User, error) {
	for _, u := range i.users {
		if u.UserID == userID {
			return &u, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (i *inMemoryGitHubUserRepo) Update(ouser *github.User) error {
	panic("implement me")
}
