package resolvers

import (
	"context"

	"getsturdy.com/api/pkg/users"
	"github.com/graph-gophers/graphql-go"
)

type GitHubAccountRootResolver interface {
	// internal
	InteralByID(context.Context, users.ID) (GitHubAccountResolver, error)
}

type GitHubAccountResolver interface {
	ID() graphql.ID
	Login() string
}
