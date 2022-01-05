package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type ViewStatusRootResolver interface {
	InternalViewStatus(ctx context.Context, viewID string) (ViewStatusResolver, error)
}

type ViewStatusResolver interface {
	ID() graphql.ID
	State() string
	ProgressPath() *string
	ProgressReceived() *int32
	ProgressTotal() *int32
	LastError() *string
	SturdyVersion() string
	UpdatedAt() int32
}
