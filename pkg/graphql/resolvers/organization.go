package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type OrganizationRootResolver interface {
	Organizations(context.Context) ([]OrganizationResolver, error)
}

type OrganizationResolver interface {
	ID() graphql.ID
	Name() string
	Members(context.Context) ([]AuthorResolver, error)
	Codebases(context.Context) ([]CodebaseResolver, error)
}
