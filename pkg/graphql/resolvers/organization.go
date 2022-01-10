package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type OrganizationRootResolver interface {
	Organizations(context.Context) ([]OrganizationResolver, error)

	// Mutations
	CreateOrganization(context.Context, CreateOrganizationArgs) (OrganizationResolver, error)
	AddUserToOrganization(context.Context, AddUserToOrganizationArgs) (OrganizationResolver, error)
}

type OrganizationResolver interface {
	ID() graphql.ID
	Name() string
	Members(context.Context) ([]AuthorResolver, error)
	Codebases(context.Context) ([]CodebaseResolver, error)
}

type CreateOrganizationArgs struct {
	Input CreateOrganizationInput
}

type CreateOrganizationInput struct {
	Name string
}

type AddUserToOrganizationArgs struct {
	Input AddUserToOrganizationInput
}

type AddUserToOrganizationInput struct {
	OrganizationID graphql.ID
	Email          string
}
