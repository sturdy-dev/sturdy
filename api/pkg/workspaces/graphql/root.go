package graphql

import (
	"context"
	"fmt"

	service_auth "getsturdy.com/api/pkg/auth/service"
	service_change "getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/codebases"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	db_comments "getsturdy.com/api/pkg/comments/db"
	"getsturdy.com/api/pkg/events"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	service_suggestions "getsturdy.com/api/pkg/suggestions/service"
	service_user "getsturdy.com/api/pkg/users/service"
	db_view "getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs/executor"

	"go.uber.org/zap"
)

type WorkspaceRootResolver struct {
	workspaceReader db_workspaces.WorkspaceReader
	codebaseRepo    db_codebases.CodebaseRepository
	viewRepo        db_view.Repository
	commentRepo     db_comments.Repository
	snapshotsRepo   db_snapshots.Repository

	codebaseResolver              resolvers.CodebaseRootResolver
	authorResolver                resolvers.AuthorRootResolver
	viewResolver                  resolvers.ViewRootResolver
	commentResolver               resolvers.CommentRootResolver
	prResolver                    resolvers.GitHubPullRequestRootResolver
	changeResolver                resolvers.ChangeRootResolver
	workspaceActivityRootResolver resolvers.ActivityRootResolver
	reviewRootResolver            resolvers.ReviewRootResolver
	presenceRootResolver          resolvers.PresenceRootResolver
	suggestionRootResolver        resolvers.SuggestionRootResolver
	statusRootResolver            resolvers.StatusesRootResolver
	workspaceWatcherRootResolver  resolvers.WorkspaceWatcherRootResolver
	fileDiffRootResolver          resolvers.FileDiffRootResolver
	rebaseStatusRootResolver      resolvers.RebaseStatusRootResolver
	downloadsResolver             resolvers.ContentsDownloadUrlRootResolver
	snapshotsResolver             resolvers.SnapshotsRootResolver

	suggestionsService *service_suggestions.Service
	workspaceService   *service_workspace.Service
	authService        *service_auth.Service
	changeService      *service_change.Service
	userService        service_user.Service

	logger           *zap.Logger
	viewEvents       events.EventReadWriter
	workspaceWriter  db_workspaces.WorkspaceWriter
	executorProvider executor.Provider
	eventsSender     events.EventSender
	eventsSubscriber *eventsv2.Subscriber
	gitSnapshotter   *service_snapshots.Service
}

func NewResolver(
	workspaceReader db_workspaces.WorkspaceReader,
	codebaseRepo db_codebases.CodebaseRepository,
	viewRepo db_view.Repository,
	commentRepo db_comments.Repository,
	snapshotRepo db_snapshots.Repository,

	codebaseResolver resolvers.CodebaseRootResolver,
	authorResolver resolvers.AuthorRootResolver,
	viewResolver resolvers.ViewRootResolver,
	commentResolver resolvers.CommentRootResolver,
	prResolver resolvers.GitHubPullRequestRootResolver,
	changeResolver resolvers.ChangeRootResolver,
	workspaceActivityRootResolver resolvers.ActivityRootResolver,
	reviewRootResolver resolvers.ReviewRootResolver,
	presenceRootResolver resolvers.PresenceRootResolver,
	suggestionRootResolver resolvers.SuggestionRootResolver,
	statusRootResolver resolvers.StatusesRootResolver,
	workspaceWatcherRootResolver resolvers.WorkspaceWatcherRootResolver,
	fileDiffRootResolver resolvers.FileDiffRootResolver,
	rebaseStatusRootResolver resolvers.RebaseStatusRootResolver,
	downloadsResolver resolvers.ContentsDownloadUrlRootResolver,
	snapshotsResolver resolvers.SnapshotsRootResolver,

	suggestionsService *service_suggestions.Service,
	workspaceService *service_workspace.Service,
	authService *service_auth.Service,
	changeService *service_change.Service,
	userService service_user.Service,

	logger *zap.Logger,
	viewEventsWriter events.EventReadWriter,
	workspaceWriter db_workspaces.WorkspaceWriter,
	executorProvider executor.Provider,
	eventsSender events.EventSender,
	eventsSubscriber *eventsv2.Subscriber,
	gitSnapshotter *service_snapshots.Service,
) resolvers.WorkspaceRootResolver {
	return &WorkspaceRootResolver{
		workspaceReader: workspaceReader,
		codebaseRepo:    codebaseRepo,
		viewRepo:        viewRepo,
		commentRepo:     commentRepo,
		snapshotsRepo:   snapshotRepo,

		codebaseResolver:              codebaseResolver,
		authorResolver:                authorResolver,
		viewResolver:                  viewResolver,
		commentResolver:               commentResolver,
		prResolver:                    prResolver,
		changeResolver:                changeResolver,
		workspaceActivityRootResolver: workspaceActivityRootResolver,
		reviewRootResolver:            reviewRootResolver,
		presenceRootResolver:          presenceRootResolver,
		suggestionRootResolver:        suggestionRootResolver,
		statusRootResolver:            statusRootResolver,
		workspaceWatcherRootResolver:  workspaceWatcherRootResolver,
		fileDiffRootResolver:          fileDiffRootResolver,
		rebaseStatusRootResolver:      rebaseStatusRootResolver,
		downloadsResolver:             downloadsResolver,
		snapshotsResolver:             snapshotsResolver,

		suggestionsService: suggestionsService,
		workspaceService:   workspaceService,
		authService:        authService,
		changeService:      changeService,
		userService:        userService,

		logger:           logger.Named("workspaceRootResolver"),
		viewEvents:       viewEventsWriter,
		workspaceWriter:  workspaceWriter,
		executorProvider: executorProvider,
		eventsSender:     eventsSender,
		eventsSubscriber: eventsSubscriber,

		gitSnapshotter: gitSnapshotter,
	}
}

func (r *WorkspaceRootResolver) Workspace(ctx context.Context, args resolvers.WorkspaceArgs) (resolvers.WorkspaceResolver, error) {
	// Return single
	ws, err := r.workspaceReader.Get(string(args.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if ws.ArchivedAt != nil && (args.AllowArchived == nil || !*args.AllowArchived) {
		return nil, gqlerrors.Error(gqlerrors.ErrNotFound)
	}

	if err := r.authService.CanRead(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &WorkspaceResolver{w: ws, root: r}, nil
}

func (r *WorkspaceRootResolver) InternalWorkspace(ws *workspaces.Workspace) resolvers.WorkspaceResolver {
	return &WorkspaceResolver{w: ws, root: r}
}

func (r *WorkspaceRootResolver) Workspaces(ctx context.Context, args resolvers.WorkspacesArgs) ([]resolvers.WorkspaceResolver, error) {
	codebaseID := codebases.ID(args.CodebaseID)
	cb, err := r.codebaseRepo.Get(codebaseID)
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("codebase not found: %w", err))
	}
	if err := r.authService.CanRead(ctx, cb); err != nil {
		return nil, gqlerrors.Error(err)
	}

	var includeDeleted bool
	if args.IncludeArchived != nil && *args.IncludeArchived {
		includeDeleted = true
	}

	workspaces, err := r.workspaceReader.ListByCodebaseIDs([]codebases.ID{codebaseID}, includeDeleted)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	var res []resolvers.WorkspaceResolver
	for _, ws := range workspaces {
		res = append(res, &WorkspaceResolver{w: ws, root: r})
	}

	return res, nil
}

func (r *WorkspaceRootResolver) ArchiveWorkspace(ctx context.Context, args resolvers.ArchiveWorkspaceArgs) (resolvers.WorkspaceResolver, error) {
	ws, err := r.workspaceReader.Get(string(args.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.workspaceService.Archive(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &WorkspaceResolver{w: ws, root: r}, nil
}

func (r *WorkspaceRootResolver) UnarchiveWorkspace(ctx context.Context, args resolvers.UnarchiveWorkspaceArgs) (resolvers.WorkspaceResolver, error) {
	ws, err := r.workspaceReader.Get(string(args.ID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.workspaceService.Unarchive(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &WorkspaceResolver{w: ws, root: r}, nil
}

func (r *WorkspaceRootResolver) SetWorkspaceSnapshot(ctx context.Context, args resolvers.SetWorkspaceSnapshotArgs) (resolvers.WorkspaceResolver, error) {
	ws, err := r.workspaceService.GetByID(ctx, string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	snap, err := r.snapshotsRepo.Get(snapshots.ID(args.Input.SnapshotID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := r.workspaceService.SetSnapshot(ctx, ws, snap); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &WorkspaceResolver{w: ws, root: r}, nil
}
