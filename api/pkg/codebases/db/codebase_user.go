package db

import (
	"context"
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
	_, err := r.db.NamedExec(`INSERT INTO codebase_users (id, user_id, codebase_id, created_at)
		VALUES (:id, :user_id, :codebase_id, :created_at)`, &entity)
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
