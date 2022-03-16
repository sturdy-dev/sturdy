package inmemory

import (
	"database/sql"

	"getsturdy.com/api/pkg/github"
	"getsturdy.com/api/pkg/users"
)

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
