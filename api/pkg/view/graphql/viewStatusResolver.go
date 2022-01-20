package graphql

import (
	"context"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/mutagen"
	db_mutagen "getsturdy.com/api/pkg/mutagen/db"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

type ViewStatusRootResolver struct {
	viewStatusRepo db_mutagen.ViewStatusRepository
	logger         *zap.Logger
}

func NewViewStatusRootResolver(
	viewStatusRepo db_mutagen.ViewStatusRepository,
	logger *zap.Logger,
) resolvers.ViewStatusRootResolver {
	return &ViewStatusRootResolver{
		viewStatusRepo: viewStatusRepo,
		logger:         logger,
	}
}

func (r *ViewStatusRootResolver) InternalViewStatus(ctx context.Context, viewID string) (resolvers.ViewStatusResolver, error) {
	data, err := r.viewStatusRepo.GetByViewID(viewID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return &viewStatusResolver{
		root: r,
		data: data,
		id:   viewID,
	}, nil
}

type viewStatusResolver struct {
	root *ViewStatusRootResolver
	data *mutagen.ViewStatus
	id   string
}

func (r *viewStatusResolver) ID() graphql.ID {
	return graphql.ID(r.id)
}

func (r *viewStatusResolver) State() string {
	switch r.data.State {
	case mutagen.ViewStatusStateDisconnected, mutagen.ViewStatusStateHaltedOnRootEmptied, mutagen.ViewStatusStateHaltedOnRootDeletion, mutagen.ViewStatusStateHaltedOnRootTypeChange:
		return "Disconnected"
	case mutagen.ViewStatusStateConnectingAlpha, mutagen.ViewStatusStateConnectingBeta:
		return "Connecting"
	case mutagen.ViewStatusStateReconciling:
		return "Reconciling"
	case mutagen.ViewStatusStateScanning:
		return "Scanning"
	case mutagen.ViewStatusStateWatching, mutagen.ViewStatusStateWaitingForRescan:
		return "Ready"
	case mutagen.ViewStatusStateStagingAlpha, mutagen.ViewStatusStateStagingBeta:
		return "Transferring"
	case mutagen.ViewStatusStateTransitioning, mutagen.ViewStatusStateSaving:
		return "Finishing"
	default:
		return "Ready"
	}
}

func (r *viewStatusResolver) ProgressPath() *string {
	return r.data.StagingStatusPath
}

func (r *viewStatusResolver) ProgressReceived() *int32 {
	if r.data.StagingStatusReceived == nil {
		return nil
	}
	i := int32(*r.data.StagingStatusReceived)
	return &i
}

func (r *viewStatusResolver) ProgressTotal() *int32 {
	if r.data.StagingStatusTotal == nil {
		return nil
	}
	i := int32(*r.data.StagingStatusTotal)
	return &i
}

func (r *viewStatusResolver) LastError() *string {
	return r.data.LastError
}

func (r *viewStatusResolver) SturdyVersion() string {
	return r.data.SturdyVersion
}

func (r *viewStatusResolver) UpdatedAt() int32 {
	return int32(r.data.UpdatedAt.Unix())
}
