package inmemory

import (
	"context"
	"database/sql"

	"getsturdy.com/api/pkg/codebase"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
)

type inMemoryCodebaseUserRepository struct {
	users []codebase.CodebaseUser
}

func NewInMemoryCodebaseUserRepo() db_codebase.CodebaseUserRepository {
	return &inMemoryCodebaseUserRepository{users: make([]codebase.CodebaseUser, 0)}
}

func (r *inMemoryCodebaseUserRepository) Create(entity codebase.CodebaseUser) error {
	r.users = append(r.users, entity)
	return nil
}

func (r *inMemoryCodebaseUserRepository) GetByUser(userID string) ([]*codebase.CodebaseUser, error) {
	var res []*codebase.CodebaseUser
	for _, u := range r.users {
		if u.UserID == userID {
			u2 := u
			res = append(res, &u2)
		}
	}
	return res, nil
}

func (r *inMemoryCodebaseUserRepository) GetByCodebase(codebaseID string) ([]*codebase.CodebaseUser, error) {
	var res []*codebase.CodebaseUser
	for _, u := range r.users {
		if u.CodebaseID == codebaseID {
			u2 := u
			res = append(res, &u2)
		}
	}
	return res, nil
}

func (r *inMemoryCodebaseUserRepository) GetByUserAndCodebase(userID, codebaseID string) (*codebase.CodebaseUser, error) {
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
			r.users[i] = r.users[len(r.users)-1]              // Copy last element to index i.
			r.users[len(r.users)-1] = codebase.CodebaseUser{} // Erase last element (write zero value).
			r.users = r.users[:len(r.users)-1]                // Truncate slice.
		}
	}
	return nil
}
