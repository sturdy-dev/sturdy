package db

import (
	"context"
	"fmt"

	"mash/pkg/jwt/keys"

	"github.com/jmoiron/sqlx"
)

var _ Repository = &database{}

type database struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *database {
	return &database{
		db: db,
	}
}

func (db *database) Create(ctx context.Context, key *keys.Key) error {
	if _, err := db.db.NamedExecContext(ctx, `
	INSERT INTO jwt_keys (
		id,
		public_der
	) VALUES (
		:id, :public_der
	)
	`, key); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}
	return nil
}

func (db *database) Get(ctx context.Context, id string) (*keys.Key, error) {
	key := &keys.Key{}
	if err := db.db.GetContext(ctx, key, `
	SELECT
		id, public_der
	FROM
		jwt_keys
	WHERE
		id = $1
	`, id); err != nil {
		return nil, err
	}
	return key, nil
}
