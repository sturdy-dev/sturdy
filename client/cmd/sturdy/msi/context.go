package main

import (
	"context"
	"errors"
)

type contextKey int

const (
	ContextNotarizationUuidKey contextKey = iota
	ContextLauncherVersionKey
	ContextOsqueryVersionKey
)

// InitContext adds several pointers to the context to allow
// for data smuggling. Not the simplest way to move data around, but
// it allows us not to adjust all the returns.
func InitContext(ctx context.Context) context.Context {
	for _, key := range []contextKey{
		ContextNotarizationUuidKey,
		ContextLauncherVersionKey,
		ContextOsqueryVersionKey,
	} {
		var strPointer *string
		s := ""
		strPointer = &s

		ctx = context.WithValue(ctx, key, strPointer)
	}

	return ctx
}

func setInContext(ctx context.Context, key contextKey, val string) {
	// If there's no pointer, then there's no point in setting
	// this. It won't get back to the caller.
	ptr, ok := ctx.Value(key).(*string)
	if !ok || ptr == nil {
		return
	}
	*ptr = val
}

func GetFromContext(ctx context.Context, key contextKey) (string, error) {
	ptr, ok := ctx.Value(key).(*string)
	if !ok {
		return "", errors.New("Context didn't have a string pointer")
	}
	return *ptr, nil
}
