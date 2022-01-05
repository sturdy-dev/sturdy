package inmemory

import (
	"database/sql"
	"mash/pkg/github"
	db_github "mash/pkg/github/db"
)

func NewInMemoryGitHubUserRepo() db_github.GitHubUserRepo {
	return &inMemoryGitHubUserRepo{
		users: make([]github.GitHubUser, 0),
	}
}

type inMemoryGitHubUserRepo struct {
	users []github.GitHubUser
}

func (i *inMemoryGitHubUserRepo) Create(user github.GitHubUser) error {
	i.users = append(i.users, user)
	return nil
}

func (i *inMemoryGitHubUserRepo) GetByUsername(username string) (*github.GitHubUser, error) {
	panic("implement me")
}

func (i *inMemoryGitHubUserRepo) GetByUserID(userID string) (*github.GitHubUser, error) {
	for _, u := range i.users {
		if u.UserID == userID {
			return &u, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (i *inMemoryGitHubUserRepo) Update(ouser *github.GitHubUser) error {
	panic("implement me")
}
