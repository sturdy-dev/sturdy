package graphql

import (
	"context"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
)

type rootResolver struct {
	snapshotService *service_snapshots.Service
}

func NewRoot(
	snapshotService *service_snapshots.Service,
) resolvers.SnapshotsRootResolver {
	return &rootResolver{
		snapshotService: snapshotService,
	}
}

func (r *rootResolver) InternalByID(ctx context.Context, id string) (resolvers.SnapshotResolver, error) {
	snap, err := r.snapshotService.GetByID(ctx, id)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &resolver{
		snapshot: snap,
		root:     r,
	}, nil
}
