package db

import (
	"fmt"

	"getsturdy.com/api/pkg/pki"

	"github.com/jmoiron/sqlx"
)

type Repo interface {
	Create(upk pki.UserPublicKey) error
	GetByPublicKeyAndUserID(publicKey, userID string) (*pki.UserPublicKey, error)
	GetKeyByUserID(userID string) ([]pki.UserPublicKey, error)
}

type dbrepo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &dbrepo{db: db}
}

// The ID value is set inside this method
func (r *dbrepo) Create(upk pki.UserPublicKey) error {
	_, err := r.db.NamedExec(`INSERT INTO user_public_keys (public_key, user_id, added_at)
		VALUES (:public_key, :user_id, :added_at)`, &upk)
	if err != nil {
		return fmt.Errorf("failed to perform insert: %w", err)
	}
	return nil
}

func (r *dbrepo) GetByPublicKeyAndUserID(publicKey, userID string) (*pki.UserPublicKey, error) {
	var upk pki.UserPublicKey
	err := r.db.Get(&upk, "SELECT * FROM user_public_keys WHERE public_key=$1 AND user_id=$2", publicKey, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return &upk, nil
}

func (r *dbrepo) GetKeyByUserID(userID string) ([]pki.UserPublicKey, error) {
	var keys []pki.UserPublicKey
	err := r.db.Select(&keys, "SELECT * FROM user_public_keys WHERE user_id=$1", userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query table: %w", err)
	}
	return keys, nil
}
