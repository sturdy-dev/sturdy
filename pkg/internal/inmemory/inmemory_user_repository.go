package inmemory

import (
	"context"
	"mash/pkg/user"
	db_user "mash/pkg/user/db"
)

// inMemoryUserRepo implements user.Repository
type inMemoryUserRepo struct {
	users []*user.User
}

func NewInMemoryUserRepo() db_user.Repository {
	return &inMemoryUserRepo{}
}

func (f *inMemoryUserRepo) Create(newUser *user.User) error {
	panic("not implemented")
}
func (f *inMemoryUserRepo) Get(id string) (*user.User, error) {
	return &user.User{
		ID:    id,
		Name:  "Test Testsson",
		Email: "email@pls.com",
	}, nil
}

func (f *inMemoryUserRepo) GetByIDs(_ context.Context, ids ...string) ([]*user.User, error) {
	return nil, nil
}

func (f *inMemoryUserRepo) GetByEmail(email string) (*user.User, error) {
	panic("not implemented")
}
func (f *inMemoryUserRepo) Update(u *user.User) error {
	panic("not implemented")
}
func (f *inMemoryUserRepo) UpdatePassword(u *user.User) error {
	panic("not implemented")
}
