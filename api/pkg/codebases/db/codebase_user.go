package db

import (
	"context"
	"database/sql"
	"fmt"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/users"

	"github.com/jmoiron/sqlx"
)

type CodebaseUserRepository interface {
	Create(codebases.CodebaseUser) error
	GetByID(context.Context, string) (*codebases.CodebaseUser, error)
	GetByUser(users.ID) ([]*codebases.CodebaseUser, error)
	GetByCodebase(codebases.ID) ([]*codebases.CodebaseUser, error)
	GetByUserAndCodebase(userID users.ID, codebaseID codebases.ID) (*codebases.CodebaseUser, error)
	DeleteByID(context.Context, string) error
}

type codebaseUserRepo struct {
	db *sqlx.DB
}

func NewCodebaseUserRepo(db *sqlx.DB) CodebaseUserRepository {
	return &codebaseUserRepo{db: db}
}

func (r *codebaseUserRepo) GetByID(ctx context.Context, id string) (*codebases.CodebaseUser, error) {
	var cb codebases.CodebaseUser
	err := r.db.GetContext(ctx, &cb, "SELECT * FROM codebase_users WHERE id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &cb, nil
}

func (r *codebaseUserRepo) Create(entity codebases.CodebaseUser) error {
	_, err := r.db.NamedExec(`INSERT INTO codebase_users (id, user_id, codebase_id, created_at, invited_by)
		VALUES (:id, :user_id, :codebase_id, :created_at, :invited_by)`, &entity)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *codebaseUserRepo) GetByUser(userID users.ID) ([]*codebases.CodebaseUser, error) {
	var entities []*codebases.CodebaseUser
	err := r.db.Select(&entities, "SELECT * FROM codebase_users WHERE user_id=$1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entities, nil
}

func (r *codebaseUserRepo) GetByCodebase(codebaseID codebases.ID) ([]*codebases.CodebaseUser, error) {
	var entities []*codebases.CodebaseUser
	err := r.db.Select(&entities, "SELECT * FROM codebase_users WHERE codebase_id=$1", codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entities, nil
}

func (r *codebaseUserRepo) GetByUserAndCodebase(userID users.ID, codebaseID codebases.ID) (*codebases.CodebaseUser, error) {
	var cb codebases.CodebaseUser
	err := r.db.Get(&cb, "SELECT * FROM codebase_users WHERE user_id = $1 AND codebase_id = $2 LIMIT 1", userID, codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &cb, nil
}

func (r *codebaseUserRepo) DeleteByID(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM codebase_users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete codebase_users by id: %w", err)
	}
	return nil
}

type inMemoryCodebaseUserRepository struct {
	users []codebases.CodebaseUser
}

func NewInMemoryCodebaseUserRepo() CodebaseUserRepository {
	return &inMemoryCodebaseUserRepository{}
}

func (r *inMemoryCodebaseUserRepository) GetByID(_ context.Context, id string) (*codebases.CodebaseUser, error) {
	for _, u := range r.users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, sql.ErrNoRows
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
