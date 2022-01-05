package db

import (
	"context"
	"fmt"

	"mash/pkg/onetime"

	"github.com/jmoiron/sqlx"
)

type database struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *database {
	return &database{db: db}
}

func (d *database) Create(ctx context.Context, token *onetime.Token) error {
	if _, err := d.db.NamedExecContext(ctx, `
		INSERT INTO onetime_tokens (
			key,
			user_id,
			created_at,
			clicks
		) VALUES (
			:key,
			:user_id,
			:created_at,
			:clicks
		)
	`, token); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (d *database) Update(ctx context.Context, token *onetime.Token) error {
	if _, err := d.db.NamedExecContext(ctx, `
		UPDATE onetime_tokens SET
			clicks = :clicks
		WHERE key = :key
	`, token); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	return nil
}

func (m *database) Get(ctx context.Context, userID, key string) (*onetime.Token, error) {
	var token onetime.Token
	if err := m.db.GetContext(ctx, &token, `
		SELECT
			key,
			user_id,
			created_at,
			clicks
		FROM onetime_tokens
		WHERE user_id = $1 AND key = $2
	`, userID, key); err != nil {
		return nil, fmt.Errorf("failed to get: %w", err)
	}
	return &token, nil
}
