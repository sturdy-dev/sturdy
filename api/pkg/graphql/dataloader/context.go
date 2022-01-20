package dataloader

import (
	"context"

	"github.com/google/uuid"
	"github.com/graph-gophers/dataloader/v6"
)

type contextKeyType struct{}

var contextKey = contextKeyType{}

// returns a new context to use for the dataloader cache
func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKey, uuid.NewString())
}

func keyFromContext(ctx context.Context, key dataloader.Key) (dataloader.StringKey, bool) {
	contextID, ok := ctx.Value(contextKey).(string)
	if !ok {
		return "", false
	}
	return dataloader.StringKey(contextID + "/" + key.String()), true
}
