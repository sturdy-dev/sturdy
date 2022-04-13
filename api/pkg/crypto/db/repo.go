package db

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"getsturdy.com/api/pkg/crypto"
)

type KeyPairRepository interface {
	Get(ctx context.Context, id crypto.KeyPairID) (*crypto.KeyPair, error)
	Create(ctx context.Context, kp crypto.KeyPair) error
}

type repo struct {
	db *sqlx.DB
}

func New(d *sqlx.DB) KeyPairRepository {
	return &repo{db: d}
}

func (r *repo) Get(ctx context.Context, id crypto.KeyPairID) (*crypto.KeyPair, error) {
	var kp crypto.KeyPair
	err := r.db.GetContext(ctx, &kp, `SELECT id, public_key, private_key, created_at, created_by, last_used_at FROM keypairs WHERE id = $1`, id)
	if err != nil {
		return nil, fmt.Errorf("could not get keypair: %w", err)
	}
	return &kp, nil
}

func (r *repo) Create(ctx context.Context, kp crypto.KeyPair) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO keypairs (id, public_key, private_key, created_at, created_by, last_used_at)
		VALUES(:id, :public_key, :private_key, :created_at, :created_by, :last_used_at)`, kp)
	if err != nil {
		return fmt.Errorf("could not save keypair: %w", err)
	}
	return nil
}
