package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type GitHubAccountRootResolver interface {
	// internal
	InteralByID(context.Context, string) (GitHubAccountResolver, error)
}

type GitHubAccountResolver interface {
	ID() graphql.ID
	Login() string
}
