package graphql

import (
	"context"
	"errors"
	"fmt"
	"time"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/change"
	service_changes "getsturdy.com/api/pkg/change/service"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/statuses"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	"getsturdy.com/api/pkg/view/events"
	service_workspace "getsturdy.com/api/pkg/workspace/service"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

type RootResolver struct {
	logger *zap.Logger

	svc              *service_statuses.Service
	changeService    *service_changes.Service
	workspaceService service_workspace.Service
	authService      *service_auth.Service

	changeRootResolver resolvers.ChangeRootResolver
	gitHubPrResovler   resolvers.GitHubPullRequestRootResolver

	eventsReader events.EventReader
}

func New(
	logger *zap.Logger,
	svc *service_statuses.Service,
	changeService *service_changes.Service,
	workspaceService service_workspace.Service,
	authService *service_auth.Service,
	changeRootResolver resolvers.ChangeRootResolver,
	gitHubPrResovler resolvers.GitHubPullRequestRootResolver,
	eventsReader events.EventReader,
) *RootResolver {
	return &RootResolver{
		logger:             logger,
		svc:                svc,
		changeService:      changeService,
		workspaceService:   workspaceService,
		authService:        authService,
		changeRootResolver: changeRootResolver,
		gitHubPrResovler:   gitHubPrResovler,
		eventsReader:       eventsReader,
	}
}

func (r *RootResolver) InternalStatus(status *statuses.Status) resolvers.StatusResolver {
	return &resolver{status: status, root: r}
}

func (r *RootResolver) InteralStatusesByCodebaseIDAndCommitID(ctx context.Context, codebaseID, commitID string) ([]resolvers.StatusResolver, error) {
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

func (r *RootResolver) UpdateStatus(ctx context.Context, args resolvers.UpdateStatusArgs) (resolvers.StatusResolver, error) {
	ch, err := r.changeService.GetChangeByID(ctx, change.ID(args.Input.ChangeID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	chCommit, err := r.changeService.GetChangeCommitOnTrunkByChangeID(ctx, change.ID(args.Input.ChangeID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, ch); err != nil {
		return nil, gqlerrors.Error(err)
	}

	var tp statuses.Type
	switch args.Input.Type {
	case resolvers.StatusTypePending:
		tp = statuses.TypePending
	case resolvers.StatusTypeFailing:
		tp = statuses.TypeFailing
	case resolvers.StatusTypeHealthy:
		tp = statuses.TypeHealty
	default:
		return nil, gqlerrors.Error(gqlerrors.ErrBadRequest, "unsupported type")
	}

	status := &statuses.Status{
		ID:          uuid.NewString(),
		CommitID:    chCommit.CommitID,
		CodebaseID:  ch.CodebaseID,
		Type:        tp,
		Title:       args.Input.Title,
		Description: args.Input.Description,
		DetailsURL:  args.Input.DetailsUrl,
		Timestamp:   time.Now(),
	}

	if err := r.svc.Set(ctx, status); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &resolver{root: r, status: status}, nil
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
		CommitID:   (*graphql.ID)(&r.status.CommitID),
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
	if pullRequest, err := r.root.gitHubPrResovler.InternalByCodebaseIDAndHeadSHA(ctx, r.status.CodebaseID, r.status.CommitID); errors.Is(err, gqlerrors.ErrNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		return pullRequest, nil
	}
}
