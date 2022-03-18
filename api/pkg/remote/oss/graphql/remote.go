package graphql

import (
	"context"

	gqlerror "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
)

type remoteRootResolver struct{}

func New() resolvers.RemoteRootResolver {
	return &remoteRootResolver{}
}

func (r *remoteRootResolver) InternalRemoteByCodebaseID(ctx context.Context, codebaseID string) (resolvers.RemoteResolver, error) {
	return nil, gqlerror.ErrNotImplemented
}

func (r *remoteRootResolver) CreateOrUpdateCodebaseRemote(ctx context.Context, args resolvers.CreateOrUpdateCodebaseRemoteArgsArgs) (resolvers.RemoteResolver, error) {
	return nil, gqlerror.ErrNotImplemented
}

func (r *remoteRootResolver) PushWorkspaceToRemote(ctx context.Context, args resolvers.PushWorkspaceToRemoteArgs) (resolvers.WorkspaceResolver, error) {
	return nil, gqlerror.ErrNotImplemented
}

func (r *remoteRootResolver) FetchRemoteToTrunk(ctx context.Context, args resolvers.FetchRemoteToTrunkArgs) (resolvers.CodebaseResolver, error) {
	return nil, gqlerror.ErrNotImplemented
}
