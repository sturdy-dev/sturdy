package resolvers

import (
	"context"

	"github.com/graph-gophers/graphql-go"
)

type RemoteRootResolver interface {
	CreateCodebaseRemote(ctx context.Context, args CreateCodebaseRemoteArgs) (CodebaseResolver, error)
	PushWorkspaceToRemote(ctx context.Context, args PushWorkspaceToRemoteArgs) (WorkspaceResolver, error)
	FetchRemoteToTrunk(ctx context.Context, args FetchRemoteToTrunkArgs) (CodebaseResolver, error)
}

type RemoteResolver interface {
	ID() graphql.ID
	Name() string
	URL() string
}

type CreateCodebaseRemoteArgs struct {
	Input CreateCodebaseRemoteInput
}

type CreateCodebaseRemoteInput struct {
	CodebaseID        string
	Name              string
	Url               string
	BasicAuthUsername string
	BasicAuthPassword string
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
