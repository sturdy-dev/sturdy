package db

import (
	"fmt"
	"mash/pkg/codebase"

	"github.com/jmoiron/sqlx"
)

type CodebaseUserRepository interface {
	Create(entity codebase.CodebaseUser) error
	GetByUser(userID string) ([]*codebase.CodebaseUser, error)
	GetByCodebase(codebaseID string) ([]*codebase.CodebaseUser, error)
	GetByUserAndCodebase(userID, codebaseID string) (*codebase.CodebaseUser, error)
}

type codebaseUserRepo struct {
	db *sqlx.DB
}

func NewCodebaseUserRepo(db *sqlx.DB) CodebaseUserRepository {
	return &codebaseUserRepo{db: db}
}

func (r *codebaseUserRepo) Create(entity codebase.CodebaseUser) error {
	_, err := r.db.NamedExec(`INSERT INTO codebase_users (id, user_id, codebase_id, created_at)
		VALUES (:id, :user_id, :codebase_id, :created_at)`, &entity)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *codebaseUserRepo) GetByUser(userID string) ([]*codebase.CodebaseUser, error) {
	var entities []*codebase.CodebaseUser
	err := r.db.Select(&entities, "SELECT * FROM codebase_users WHERE user_id=$1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entities, nil
}

func (r *codebaseUserRepo) GetByCodebase(codebaseID string) ([]*codebase.CodebaseUser, error) {
	var entities []*codebase.CodebaseUser
	err := r.db.Select(&entities, "SELECT * FROM codebase_users WHERE codebase_id=$1", codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return entities, nil
}

func (r *codebaseUserRepo) GetByUserAndCodebase(userID, codebaseID string) (*codebase.CodebaseUser, error) {
	var cb codebase.CodebaseUser
	err := r.db.Get(&cb, "SELECT * FROM codebase_users WHERE user_id = $1 AND codebase_id = $2 LIMIT 1", userID, codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &cb, nil
}
