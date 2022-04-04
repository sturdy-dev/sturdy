package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type AuthorRootResolver interface {
	// Query
	Author(context.Context, graphql.ID) (AuthorResolver, error)

	// Internal
	InternalAuthorFromNameAndEmail(ctx context.Context, name, email string) AuthorResolver
	// InternalAuthors(context.Context, ...string) ([]AuthorResolver, error)
}

type AuthorResolver interface {
	ID() graphql.ID
	Name() string
	AvatarUrl() *string
	Email() string
	Status() (UserStatus, error)
}
