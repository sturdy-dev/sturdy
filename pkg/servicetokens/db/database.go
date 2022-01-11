package db

import (
	"context"
	"fmt"
	"mash/pkg/servicetokens"

	"github.com/jmoiron/sqlx"
)

var _ Repository = &database{}

type database struct {
	db *sqlx.DB
}

func NewDatabase(db *sqlx.DB) Repository {
	return &database{
		db: db,
	}
}

func (d *database) Create(ctx context.Context, token *servicetokens.Token) error {
	if _, err := d.db.NamedExecContext(ctx, `
		INSERT INTO servicetokens (
			id, codebase_id, hash, name, created_at, last_used_at
		) VALUES (
			:id, :codebase_id, :hash, :name, :created_at, :last_used_at
		)
	`, token); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (d *database) GetByID(ctx context.Context, id string) (*servicetokens.Token, error) {
	token := &servicetokens.Token{}
	if err := d.db.GetContext(ctx, token, `
		SELECT
			id, codebase_id, hash, name, created_at, last_used_at
		FROM servicetokens
		WHERE id = $1
	`, id); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}
	return token, nil
}
