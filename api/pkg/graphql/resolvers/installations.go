package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type InstallationsRootResolver interface {
	Installation(context.Context) (InstallationsResolver, error)
}

type InstallationsResolver interface {
	ID() graphql.ID
	NeedsFirstTimeSetup(context.Context) (bool, error)
	Version() string
	License(context.Context) (LicenseResolver, error)
}
