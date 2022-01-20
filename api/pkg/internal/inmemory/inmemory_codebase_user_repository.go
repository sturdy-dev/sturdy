package inmemory

import (
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
	panic("not implemented")
}
func (r *inMemoryCodebaseUserRepository) GetByCodebase(codebaseID string) ([]*codebase.CodebaseUser, error) {
	var res []*codebase.CodebaseUser
	for _, u := range r.users {
		if u.CodebaseID == codebaseID {
			res = append(res, &u)
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
