package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/changes"
	service_changes "getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/codebases"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/github"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/snapshots"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	"getsturdy.com/api/pkg/statuses"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	service_workspace_statuses "getsturdy.com/api/pkg/workspaces/statuses/service"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

type RootResolver struct {
	logger *zap.Logger

	svc                      *service_statuses.Service
	changeService            *service_changes.Service
	workspaceService         *service_workspace.Service
	authService              *service_auth.Service
	snapshotsService         *service_snapshots.Service
	workspaceStatusesService *service_workspace_statuses.Service

	changeRootResolver    resolvers.ChangeRootResolver
	gitHubPrResolver      resolvers.GitHubPullRequestRootResolver
	workspaceRootResolver *resolvers.WorkspaceRootResolver

	eventsSubscriber *eventsv2.Subscriber
}

var _ resolvers.StatusesRootResolver = &RootResolver{}

func New(
	logger *zap.Logger,

	svc *service_statuses.Service,
	changeService *service_changes.Service,
	workspaceService *service_workspace.Service,
	authService *service_auth.Service,
	snapshotsService *service_snapshots.Service,
	workspaceStatusesService *service_workspace_statuses.Service,

	changeRootResolver resolvers.ChangeRootResolver,
	gitHubPrResolver resolvers.GitHubPullRequestRootResolver,
	workspaceRootResolver *resolvers.WorkspaceRootResolver,

	eventsReader *eventsv2.Subscriber,
) *RootResolver {
	return &RootResolver{
		logger: logger,

		svc:                      svc,
		changeService:            changeService,
		workspaceService:         workspaceService,
		authService:              authService,
		snapshotsService:         snapshotsService,
		workspaceStatusesService: workspaceStatusesService,

		changeRootResolver:    changeRootResolver,
		gitHubPrResolver:      gitHubPrResolver,
		workspaceRootResolver: workspaceRootResolver,

		eventsSubscriber: eventsReader,
	}
}

func (r *RootResolver) InternalGitHubPullRequestStatuses(context.Context, *github.PullRequest) ([]resolvers.GitHubPullRequestStatusResolver, error) {
	return nil, gqlerrors.ErrNotImplemented
}

func (r *RootResolver) InternalWorkspaceStatuses(ctx context.Context, workspaceID string) ([]resolvers.WorkspaceStatusResolver, error) {
	ss, err := r.svc.ListByWorkspaceID(ctx, workspaceID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	rr := make([]resolvers.WorkspaceStatusResolver, 0, len(ss))
	for _, s := range ss {
		rr = append(rr, &workspaceResolver{
			resolver: &resolver{status: s, root: r},
		})
	}
	return rr, nil
}

func (r *RootResolver) InternalChangeStatuses(ctx context.Context, change *changes.Change) ([]resolvers.ChangeStatusResolver, error) {
	if change.CommitID == nil {
		return nil, gqlerrors.ErrNotFound
	}
	ss, err := r.svc.List(ctx, change.CodebaseID, *change.CommitID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	rr := make([]resolvers.ChangeStatusResolver, 0, len(ss))
	for _, s := range ss {
		rr = append(rr, &changeResolver{
			resolver: &resolver{status: s, root: r},
		})
	}
	return rr, nil
}

func (r *RootResolver) InternalStatus(status *statuses.Status) resolvers.StatusResolver {
	return &resolver{status: status, root: r}
}

func (r *RootResolver) statusesByCodebaseIDAndCommitID(ctx context.Context, codebaseID codebases.ID, commitID string) ([]resolvers.StatusResolver, error) {
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
	case statuses.TypeHealthy:
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

func (r *resolver) ToGitHubPullRequestStatus() (resolvers.GitHubPullRequestStatusResolver, bool) {
	return nil, false
}

func (r *resolver) ToWorkspaceStatus() (resolvers.WorkspaceStatusResolver, bool) {
	return &workspaceResolver{
		resolver: &resolver{status: r.status, root: r.root},
	}, true
}

func (r *resolver) ToChangeStatus() (resolvers.ChangeStatusResolver, bool) {
	return &changeResolver{
		resolver: &resolver{status: r.status, root: r.root},
	}, true
}

type changeResolver struct {
	*resolver
}

func (r *changeResolver) Change(ctx context.Context) (resolvers.ChangeResolver, error) {
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

type workspaceResolver struct {
	*resolver

	snapshot     *snapshots.Snapshot
	snapshotErr  error
	snapshotOnce sync.Once
}

func (r *workspaceResolver) getSnapshot(ctx context.Context) (*snapshots.Snapshot, error) {
	r.snapshotOnce.Do(func() {
		r.snapshot, r.snapshotErr = r.root.snapshotsService.GetByCommitSHA(ctx, r.status.CommitSHA)
	})
	return r.snapshot, r.snapshotErr
}

func (r *workspaceResolver) Stale(ctx context.Context) (bool, error) {
	if snapshot, err := r.getSnapshot(ctx); err != nil {
		return true, err
	} else if ws, err := r.root.workspaceService.GetByID(ctx, snapshot.WorkspaceID); errors.Is(err, sql.ErrNoRows) {
		return true, nil
	} else if isStale, err := r.root.workspaceStatusesService.StatusIsStaleForWorkspace(ctx, ws, r.status); err != nil {
		return false, err
	} else {
		return isStale, nil
	}
}

func (r *workspaceResolver) Workspace(ctx context.Context) (resolvers.WorkspaceResolver, error) {
	if snapshot, err := r.getSnapshot(ctx); err != nil {
		return nil, gqlerrors.Error(err)
	} else {
		t := true
		return (*r.root.workspaceRootResolver).Workspace(ctx, resolvers.WorkspaceArgs{
			ID:            graphql.ID(snapshot.WorkspaceID),
			AllowArchived: &t,
		})
	}
}
