package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	service_auth "mash/pkg/auth/service"
	"mash/pkg/change"
	service_changes "mash/pkg/change/service"
	db_github "mash/pkg/github/db"
	gqlerrors "mash/pkg/graphql/errors"
	"mash/pkg/graphql/resolvers"
	"mash/pkg/statuses"
	service_statuses "mash/pkg/statuses/service"
	"mash/pkg/view/events"
	service_workspace "mash/pkg/workspace/service"
	"time"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

type rootResolver struct {
	logger *zap.Logger

	svc              *service_statuses.Service
	changeService    *service_changes.Service
	workspaceService service_workspace.Service
	authService      *service_auth.Service

	gitHubPrRepo db_github.GitHubPRRepo

	changeRootResolver *resolvers.ChangeRootResolver
	gitHubPrResovler   *resolvers.GitHubPullRequestRootResolver

	eventsReader events.EventReader
}

func New(
	logger *zap.Logger,
	svc *service_statuses.Service,
	changeService *service_changes.Service,
	workspaceService service_workspace.Service,
	authService *service_auth.Service,
	gitHubPrRepo db_github.GitHubPRRepo,
	changeRootResolver *resolvers.ChangeRootResolver,
	gitHubPrResovler *resolvers.GitHubPullRequestRootResolver,
	eventsReader events.EventReader,
) resolvers.StatusesRootResolver {
	return &rootResolver{
		logger:             logger,
		svc:                svc,
		changeService:      changeService,
		workspaceService:   workspaceService,
		authService:        authService,
		gitHubPrRepo:       gitHubPrRepo,
		changeRootResolver: changeRootResolver,
		gitHubPrResovler:   gitHubPrResovler,
		eventsReader:       eventsReader,
	}
}

func (r *rootResolver) InternalStatus(status *statuses.Status) resolvers.StatusResolver {
	return &resolver{status: status, root: r}
}

func (r *rootResolver) InteralStatusesByCodebaseIDAndCommitID(ctx context.Context, codebaseID, commitID string) ([]resolvers.StatusResolver, error) {
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

func (r *rootResolver) UpdateStatus(ctx context.Context, args resolvers.UpdateStatusArgs) (resolvers.StatusResolver, error) {
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
	root   *rootResolver
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
	change, err := (*r.root.changeRootResolver).Change(ctx, resolvers.ChangeArgs{
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
	pr, err := r.root.gitHubPrRepo.GetByCodebaseIDaAndHeadSHA(ctx, r.status.CodebaseID, r.status.CommitID)
	switch {
	case err == nil:
		return (*r.root.gitHubPrResovler).InternalGitHubPullRequest(pr)
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	default:
		return nil, gqlerrors.Error(fmt.Errorf("failed to get github pr: %w", err))
	}
}
