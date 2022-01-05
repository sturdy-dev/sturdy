package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type ACLRootResolver interface {
	// Internal APIs
	InternalACLByCodebaseID(ctx context.Context, codebaseID graphql.ID) (ACLResolver, error)

	// Queries
	CanI(ctx context.Context, args CanIArgs) (bool, error)

	// Mutations
	UpdateACL(ctx context.Context, args UpdateACLArgs) (ACLResolver, error)
}

type CanIArgs struct {
	CodebaseID string
	Action     string
	Resource   string
}

type UpdateACLArgs struct {
	Input UpdateACLInput
}

type UpdateACLInput struct {
	CodebaseID graphql.ID
	Policy     *string
}

type ACLResolver interface {
	ID() graphql.ID
	Policy() (string, error)
}
