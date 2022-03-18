package inmemory

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	"getsturdy.com/api/pkg/users"
)

type inMemoryCodebaseUserRepository struct {
	users []codebases.CodebaseUser
}

func NewInMemoryCodebaseUserRepo() db_codebases.CodebaseUserRepository {
	return &inMemoryCodebaseUserRepository{users: make([]codebases.CodebaseUser, 0)}
}

func (r *inMemoryCodebaseUserRepository) Create(entity codebases.CodebaseUser) error {
	r.users = append(r.users, entity)
	return nil
}

func (r *inMemoryCodebaseUserRepository) GetByUser(userID users.ID) ([]*codebases.CodebaseUser, error) {
	var res []*codebases.CodebaseUser
	for _, u := range r.users {
		if u.UserID == userID {
			u2 := u
			res = append(res, &u2)
		}
	}
	return res, nil
}

func (r *inMemoryCodebaseUserRepository) GetByCodebase(codebaseID codebases.ID) ([]*codebases.CodebaseUser, error) {
	var res []*codebases.CodebaseUser
	for _, u := range r.users {
		if u.CodebaseID == codebaseID {
			u2 := u
			res = append(res, &u2)
		}
	}
	return res, nil
}

func (r *inMemoryCodebaseUserRepository) GetByUserAndCodebase(userID users.ID, codebaseID codebases.ID) (*codebases.CodebaseUser, error) {
	for _, u := range r.users {
		if u.UserID == userID && u.CodebaseID == codebaseID {
			return &u, nil
		}
	}
	return nil, sql.ErrNoRows
}

func (r *inMemoryCodebaseUserRepository) DeleteByID(_ context.Context, id string) error {
	for i, u := range r.users {
		if u.ID == id {
			// Remove the element at index i from a.
			r.users[i] = r.users[len(r.users)-1]               // Copy last element to index i.
			r.users[len(r.users)-1] = codebases.CodebaseUser{} // Erase last element (write zero value).
			r.users = r.users[:len(r.users)-1]                 // Truncate slice.
		}
	}
	return nil
}
