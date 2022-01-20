package dataloader

import (
	"context"

	dataloader "github.com/graph-gophers/dataloader/v6"
	"go.uber.org/zap"
)

var _ dataloader.Cache = new(ContextCache)

// ContextCache is a dataloader.Cache implementation that removes cache entries
// after the context is done.
type ContextCache struct {
	s      dataloader.Cache
	logger *zap.Logger
}

func NewContextCache(logger *zap.Logger) *ContextCache {
	return &ContextCache{
		s:      dataloader.NewCache(),
		logger: logger.Named("ContextCache"),
	}
}

func (c *ContextCache) Get(ctx context.Context, key dataloader.Key) (dataloader.Thunk, bool) {
	k, ok := keyFromContext(ctx, key)
	if !ok {
		c.logger.Error("the context does not contain a cache key, skipping")
		return nil, false
	}

	return c.s.Get(ctx, k)
}

func (c *ContextCache) Set(ctx context.Context, key dataloader.Key, thunk dataloader.Thunk) {
	k, ok := keyFromContext(ctx, key)
	if !ok {
		c.logger.Error("the context does not contain a cache key, skipping")
		return
	}

	go func() {
		<-ctx.Done()
		c.Delete(ctx, key)
	}()

	c.s.Set(ctx, k, thunk)
}

func (c *ContextCache) Delete(ctx context.Context, key dataloader.Key) bool {
	k, ok := keyFromContext(ctx, key)
	if !ok {
		c.logger.Error("the context does not contain a cache key, skipping")
		return false
	}

	return c.s.Delete(ctx, k)
}

func (c *ContextCache) Clear() {
	c.s.Clear()
}
