package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type InstallationsRootResolver interface {
	Installation(context.Context) (InstallationsResolver, error)

	UpdateInstallation(context.Context, UpdateInstallationArgs) (InstallationsResolver, error)
}

type InstallationsResolver interface {
	ID() graphql.ID
	UsersCount(context.Context) (int32, error)
	NeedsFirstTimeSetup(context.Context) (bool, error)
	Version() string
	DistributionType() string
	License(context.Context) (LicenseResolver, error)
}

type UpdateInstallationArgs struct {
	Input UpdateInstallationInput
}

type UpdateInstallationInput struct {
	LicenseKey *string
}
