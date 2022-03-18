package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type RemoteRootResolver interface {
	InternalRemoteByCodebaseID(ctx context.Context, codebaseID string) (RemoteResolver, error)

	// Mutations
	CreateOrUpdateCodebaseRemote(ctx context.Context, args CreateOrUpdateCodebaseRemoteArgsArgs) (RemoteResolver, error)
	PushWorkspaceToRemote(ctx context.Context, args PushWorkspaceToRemoteArgs) (WorkspaceResolver, error)
	FetchRemoteToTrunk(ctx context.Context, args FetchRemoteToTrunkArgs) (CodebaseResolver, error)
}

type RemoteResolver interface {
	ID() graphql.ID
	Name() string
	URL() string
	TrackedBranch() string
	BasicAuthUsername() string
	BasicAuthPassword() string
}

type CreateOrUpdateCodebaseRemoteArgsArgs struct {
	Input CreateOrUpdateCodebaseRemoteArgsInput
}

type CreateOrUpdateCodebaseRemoteArgsInput struct {
	CodebaseID        string
	Name              string
	Url               string
	BasicAuthUsername string
	BasicAuthPassword string
	TrackedBranch     string
}

type PushWorkspaceToRemoteArgs struct {
	Input PushWorkspaceToRemoteInput
}

type PushWorkspaceToRemoteInput struct {
	WorkspaceID graphql.ID
}

type FetchRemoteToTrunkArgs struct {
	Input FetchRemoteToTrunkInput
}

type FetchRemoteToTrunkInput struct {
	CodebaseID graphql.ID
}
