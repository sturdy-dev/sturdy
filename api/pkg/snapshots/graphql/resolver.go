package graphql

import (
	"context"
	"database/sql"
	"errors"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/snapshots"

	"github.com/graph-gophers/graphql-go"
)

var _ resolvers.SnapshotResolver = (*resolver)(nil)

type resolver struct {
	root     *rootResolver
	snapshot *snapshots.Snapshot
}

func (r *resolver) ID() graphql.ID {
	return graphql.ID(r.snapshot.ID)
}

func (r *resolver) Previous(ctx context.Context) (resolvers.SnapshotResolver, error) {
	if snap, err := r.root.snapshotService.Previous(ctx, r.snapshot); errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, gqlerrors.Error(err)
	} else {
		return &resolver{
			root:     r.root,
			snapshot: snap,
		}, nil
	}
}

func (r *resolver) Next(ctx context.Context) (resolvers.SnapshotResolver, error) {
	if snap, err := r.root.snapshotService.Next(ctx, r.snapshot); errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, gqlerrors.Error(err)
	} else {
		return &resolver{
			root:     r.root,
			snapshot: snap,
		}, nil
	}
}

func (r *resolver) CreatedAt() int32 {
	return int32(r.snapshot.CreatedAt.Unix())
}

func (r *resolver) Description(context.Context) (*string, error) {
	// TODO: implement
	return nil, nil
}
