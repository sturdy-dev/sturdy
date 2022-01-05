package db

import (
	"context"
	"database/sql"

	"mash/pkg/jwt/keys"
)

var _ Repository = &memory{}

type memory struct {
	byID map[string]*keys.Key
}

func NewInMemory() *memory {
	return &memory{
		byID: map[string]*keys.Key{},
	}
}

func (db *memory) Create(ctx context.Context, key *keys.Key) error {
	db.byID[key.ID] = key
	return nil
}

func (db *memory) Get(ctx context.Context, id string) (*keys.Key, error) {
	key, found := db.byID[id]
	if !found {
		return nil, sql.ErrNoRows
	}
	return key, nil
}
