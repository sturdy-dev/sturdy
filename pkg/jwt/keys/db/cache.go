package db

import (
	"context"
	"database/sql"
	"errors"
	"mash/pkg/jwt/keys"
	"sync"
)

var _ Repository = &cache{}

type cache struct {
	db Repository

	cache      map[string]*keys.Key
	cacheGuard *sync.RWMutex
}

func NewCache(db Repository) *cache {
	return &cache{
		db: db,

		cache:      make(map[string]*keys.Key),
		cacheGuard: &sync.RWMutex{},
	}
}

func (c *cache) Create(ctx context.Context, key *keys.Key) error {
	if err := c.db.Create(ctx, key); err != nil {
		return err
	}
	c.cacheGuard.Lock()
	c.cache[key.ID] = key
	c.cacheGuard.Unlock()
	return nil
}

func (c *cache) Get(ctx context.Context, id string) (*keys.Key, error) {
	c.cacheGuard.RLock()
	cached, foundInCache := c.cache[id]
	c.cacheGuard.RUnlock()

	if !foundInCache {
		key, err := c.db.Get(ctx, id)
		switch {
		case err == nil:
			c.cacheGuard.Lock()
			c.cache[id] = key
			c.cacheGuard.Unlock()
			return key, nil
		case errors.Is(err, sql.ErrNoRows):
			c.cacheGuard.Lock()
			c.cache[id] = nil
			c.cacheGuard.Unlock()
			return nil, sql.ErrNoRows
		default:
			return nil, err
		}
	} else {
		if cached == nil {
			return nil, sql.ErrNoRows
		} else {
			return cached, nil
		}
	}
}
