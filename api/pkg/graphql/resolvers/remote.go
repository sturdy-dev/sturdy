package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"

	"getsturdy.com/api/pkg/codebases"
)

type RemoteRootResolver interface {
	InternalRemoteByCodebaseID(ctx context.Context, codebaseID codebases.ID) (RemoteResolver, error)

	// Mutations
	CreateOrUpdateCodebaseRemote(ctx context.Context, args CreateOrUpdateCodebaseRemoteArgsArgs) (RemoteResolver, error)
}

type RemoteResolver interface {
	ID() graphql.ID
	Name() string
	URL() string
	TrackedBranch() string

	BasicAuthUsername() *string
	BasicAuthPassword() *string

	KeyPair(context.Context) (KeyPairResolver, error)

	BrowserLinkRepo() string
	BrowserLinkBranch() string

	Enabled() bool
}

type CreateOrUpdateCodebaseRemoteArgsArgs struct {
	Input CreateOrUpdateCodebaseRemoteArgsInput
}

type CreateOrUpdateCodebaseRemoteArgsInput struct {
	CodebaseID    string
	Name          string
	Url           string
	TrackedBranch string

	BasicAuthUsername *string
	BasicAuthPassword *string
	KeyPairID         *graphql.ID

	BrowserLinkRepo   string
	BrowserLinkBranch string

	Enabled bool
}
