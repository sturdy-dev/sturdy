package resolvers

import (
	"context"

	"getsturdy.com/api/pkg/snapshots"

	"github.com/graph-gophers/graphql-go"
)

type SnapshotsRootResolver interface {
	InternalByID(context.Context, snapshots.ID) (SnapshotResolver, error)
}

type SnapshotResolver interface {
	ID() graphql.ID
	Previous(context.Context) (SnapshotResolver, error)
	Next(context.Context) (SnapshotResolver, error)
	CreatedAt() int32
	Description(context.Context) (*string, error)
}
