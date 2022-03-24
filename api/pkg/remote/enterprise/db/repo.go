package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/remote"
)

type Repository interface {
	GetByCodebaseID(ctx context.Context, codebaseID codebases.ID) (*remote.Remote, error)
	Create(ctx context.Context, r remote.Remote) error
	Update(ctx context.Context, r *remote.Remote) error
}

func New(db *sqlx.DB) Repository {
	return &repo{db: db}
}

type repo struct {
	db *sqlx.DB
}

func (r *repo) GetByCodebaseID(ctx context.Context, codebaseID codebases.ID) (*remote.Remote, error) {
	var res remote.Remote
	err := r.db.GetContext(ctx, &res, `SELECT * FROM remotes WHERE codebase_id = $1`, codebaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to GetByCodebaseID: %w", err)
	}
	return &res, nil
}

func (r *repo) Create(ctx context.Context, val remote.Remote) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO remotes (id, codebase_id, name, url, basic_username, basic_password, tracked_branch, browser_link_repo, browser_link_branch, keypair_id)
		VALUES(:id, :codebase_id, :name, :url, :basic_username, :basic_password, :tracked_branch, :browser_link_repo, :browser_link_branch, :keypair_id)`, val)
	if err != nil {
		return fmt.Errorf("failed to create remote: %w", err)
	}
	return nil
}

func (r *repo) Update(ctx context.Context, val *remote.Remote) error {
	_, err := r.db.NamedExecContext(ctx, `
		UPDATE remotes
	    SET name = :name,
	        url = :url,
	        basic_username = :basic_username,
			basic_password = :basic_password,
			tracked_branch = :tracked_branch,
			browser_link_repo = :browser_link_repo, 
			browser_link_branch = :browser_link_branch,
			keypair_id = :keypair_id
		WHERE id = :id`, val)
	if err != nil {
		return fmt.Errorf("failed to update remote: %w", err)
	}
	return nil
}
