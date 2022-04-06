package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	service_auth "getsturdy.com/api/pkg/auth/service"
	service_changes "getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/codebases"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/statuses"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

type RootResolver struct {
	logger *zap.Logger

	svc              *service_statuses.Service
	changeService    *service_changes.Service
	workspaceService service_workspace.Service
	authService      *service_auth.Service
	snapshotsService snapshotter.Snapshotter

	changeRootResolver    resolvers.ChangeRootResolver
	gitHubPrResovler      resolvers.GitHubPullRequestRootResolver
	workspaceRootResolver *resolvers.WorkspaceRootResolver

	eventsSubscriber *eventsv2.Subscriber
}

func New(
	logger *zap.Logger,
	svc *service_statuses.Service,
	changeService *service_changes.Service,
	workspaceService service_workspace.Service,
	authService *service_auth.Service,
	changeRootResolver resolvers.ChangeRootResolver,
	gitHubPrResovler resolvers.GitHubPullRequestRootResolver,
	eventsReader *eventsv2.Subscriber,
	snapshotsService snapshotter.Snapshotter,
	workspaceRootResolver *resolvers.WorkspaceRootResolver,
) *RootResolver {
	return &RootResolver{
		logger:                logger,
		svc:                   svc,
		changeService:         changeService,
		workspaceService:      workspaceService,
		authService:           authService,
		changeRootResolver:    changeRootResolver,
		gitHubPrResovler:      gitHubPrResovler,
		eventsSubscriber:      eventsReader,
		snapshotsService:      snapshotsService,
		workspaceRootResolver: workspaceRootResolver,
	}
}

func (r *RootResolver) InternalStatus(status *statuses.Status) resolvers.StatusResolver {
	return &resolver{status: status, root: r}
}

func (r *RootResolver) InteralStatusesByCodebaseIDAndCommitID(ctx context.Context, codebaseID codebases.ID, commitID string) ([]resolvers.StatusResolver, error) {
	ss, err := r.svc.List(ctx, codebaseID, commitID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	rr := make([]resolvers.StatusResolver, 0, len(ss))
	for _, s := range ss {
		rr = append(rr, &resolver{status: s, root: r})
	}
	return rr, nil
}

type resolver struct {
	root   *RootResolver
	status *statuses.Status
}

func (r *resolver) ID() graphql.ID {
	return graphql.ID(r.status.ID)
}

func (r *resolver) DetailsUrl() *string {
	return r.status.DetailsURL
}

func (r *resolver) Type() (resolvers.StatusType, error) {
	switch r.status.Type {
	case statuses.TypePending:
		return resolvers.StatusTypePending, nil
	case statuses.TypeHealty:
		return resolvers.StatusTypeHealthy, nil
	case statuses.TypeFailing:
		return resolvers.StatusTypeFailing, nil
	default:
		return resolvers.StatusTypeUndefined, fmt.Errorf("undefined status: %s", r.status.Type)
	}
}

func (r *resolver) Title() string {
	return r.status.Title
}

func (r *resolver) Description() *string {
	return r.status.Description
}

func (r *resolver) Timestamp() int32 {
	return int32(r.status.Timestamp.Unix())
}

func (r *resolver) Change(ctx context.Context) (resolvers.ChangeResolver, error) {
	change, err := r.root.changeRootResolver.Change(ctx, resolvers.ChangeArgs{
		CodebaseID: (*graphql.ID)(&r.status.CodebaseID),
		CommitID:   (*graphql.ID)(&r.status.CommitSHA),
	})
	switch {
	case err == nil:
		return change, nil
	case errors.Is(err, gqlerrors.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (r *resolver) GitHubPullRequest(ctx context.Context) (resolvers.GitHubPullRequestResolver, error) {
	if pullRequest, err := r.root.gitHubPrResovler.InternalByCodebaseIDAndHeadSHA(ctx, r.status.CodebaseID, r.status.CommitSHA); errors.Is(err, gqlerrors.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return pullRequest, nil
	}
}

func (r *resolver) Workspace(ctx context.Context) (resolvers.WorkspaceResolver, error) {
	if snapshot, err := r.root.snapshotsService.GetByCommitSHA(ctx, r.status.CommitSHA); errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, gqlerrors.Error(err)
	} else if snapshot.WorkspaceID == nil {
		return nil, nil
	} else {
		t := true
		return (*r.root.workspaceRootResolver).Workspace(ctx, resolvers.WorkspaceArgs{
			ID:            graphql.ID(*snapshot.WorkspaceID),
			AllowArchived: &t,
		})
	}
}
