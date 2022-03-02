package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type OrganizationRootResolver interface {
	Organizations(context.Context) ([]OrganizationResolver, error)
	Organization(context.Context, OrganizationArgs) (OrganizationResolver, error)

	// Mutations
	CreateOrganization(context.Context, CreateOrganizationArgs) (OrganizationResolver, error)
	UpdateOrganization(context.Context, UpdateOrganizationArgs) (OrganizationResolver, error)
	AddUserToOrganization(context.Context, AddUserToOrganizationArgs) (OrganizationResolver, error)
	RemoveUserFromOrganization(context.Context, RemoveUserFromOrganizationArgs) (OrganizationResolver, error)
}

type OrganizationResolver interface {
	ID() graphql.ID
	ShortID() graphql.ID
	Name() string
	Members(context.Context) ([]AuthorResolver, error)
	Codebases(context.Context) ([]CodebaseResolver, error)

	Licenses(context.Context) ([]LicenseResolver, error)

	Writeable(context.Context) bool
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

type UpdateOrganizationArgs struct {
	Input UpdateOrganizationInput
}

type UpdateOrganizationInput struct {
	ID   graphql.ID
	Name string
}

type OrganizationArgs struct {
	ID      *graphql.ID
	ShortID *graphql.ID
}

type RemoveUserFromOrganizationArgs struct {
	Input RemoveUserFromOrganizationInput
}

type RemoveUserFromOrganizationInput struct {
	OrganizationID graphql.ID
	UserID         graphql.ID
}
