package inmemory

import (
	"context"

	"getsturdy.com/api/pkg/users"
	db_user "getsturdy.com/api/pkg/users/db"
)

// inMemoryUserRepo implements user.Repository
type inMemoryUserRepo struct {
	users []*users.User
}

func NewInMemoryUserRepo() db_user.Repository {
	return &inMemoryUserRepo{}
}

func (f *inMemoryUserRepo) Create(newUser *users.User) error {
	f.users = append(f.users, newUser)
	return nil
}

func (f *inMemoryUserRepo) Get(id string) (*users.User, error) {
	return &users.User{
		ID:    id,
		Name:  "Test Testsson",
		Email: "email@pls.com",
	}, nil
}

func (f *inMemoryUserRepo) GetByIDs(_ context.Context, ids ...string) ([]*users.User, error) {
	return nil, nil
}

func (f *inMemoryUserRepo) GetByEmail(email string) (*users.User, error) {
	panic("not implemented")
}

func (f *inMemoryUserRepo) Update(u *users.User) error {
	panic("not implemented")
}

func (f *inMemoryUserRepo) UpdatePassword(u *users.User) error {
	panic("not implemented")
}

func (f *inMemoryUserRepo) Count(_ context.Context) (uint64, error) {
	return uint64(len(f.users)), nil
}
