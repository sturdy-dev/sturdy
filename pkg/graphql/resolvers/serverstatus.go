package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type ServerStatusRootResolver interface {
	ServerStatus(context.Context) (ServerStatusResolver, error)
}

type ServerStatusResolver interface {
	ID() graphql.ID
	NeedsFirstTimeSetup(context.Context) (bool, error)
	Version() string
}
