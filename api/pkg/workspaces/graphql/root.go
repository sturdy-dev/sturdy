package graphql

import (
	"context"
	"fmt"

	service_auth "getsturdy.com/api/pkg/auth/service"
	service_change "getsturdy.com/api/pkg/changes/service"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	db_comments "getsturdy.com/api/pkg/comments/db"
	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	service_suggestions "getsturdy.com/api/pkg/suggestions/service"
	db_view "getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/view/open"
	"getsturdy.com/api/pkg/workspaces"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"go.uber.org/zap"
)

type WorkspaceRootResolver struct {
	workspaceReader db_workspaces.WorkspaceReader
	codebaseRepo    db_codebase.CodebaseRepository
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

	suggestionsService *service_suggestions.Service
	workspaceService   service_workspace.Service
	authService        *service_auth.Service
	changeService      *service_change.Service

	logger           *zap.Logger
	viewEvents       events.EventReadWriter
	workspaceWriter  db_workspaces.WorkspaceWriter
	executorProvider executor.Provider
	eventsSender     events.EventSender
	gitSnapshotter   snapshotter.Snapshotter
}

func NewResolver(
	workspaceReader db_workspaces.WorkspaceReader,
	codebaseRepo db_codebase.CodebaseRepository,
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

	suggestionsService *service_suggestions.Service,
	workspaceService service_workspace.Service,
	authService *service_auth.Service,
	changeService *service_change.Service,

	logger *zap.Logger,
	viewEventsWriter events.EventReadWriter,
	workspaceWriter db_workspaces.WorkspaceWriter,
	executorProvider executor.Provider,
	eventsSender events.EventSender,
	gitSnapshotter snapshotter.Snapshotter,
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

		suggestionsService: suggestionsService,
		workspaceService:   workspaceService,
		authService:        authService,
		changeService:      changeService,

		logger:           logger.Named("workspaceRootResolver"),
		viewEvents:       viewEventsWriter,
		workspaceWriter:  workspaceWriter,
		executorProvider: executorProvider,
		eventsSender:     eventsSender,

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
	cb, err := r.codebaseRepo.Get(string(args.CodebaseID))
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

	workspaces, err := r.workspaceReader.ListByCodebaseIDs([]string{string(args.CodebaseID)}, includeDeleted)
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

	if ws.ViewID == nil {
		return &WorkspaceResolver{w: ws, root: r}, nil
	}

	view, err := r.viewRepo.Get(*ws.ViewID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	// if the workspace had a view, move it to a new workspace to make sure it is always connected to a workspace

	newWorkspace, err := r.workspaceService.Create(ctx, service_workspace.CreateWorkspaceRequest{
		CodebaseID: ws.CodebaseID,
		UserID:     ws.UserID,
	})
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	if err := open.OpenWorkspaceOnView(
		ctx,
		r.logger,
		view,
		newWorkspace,
		r.viewRepo,
		r.workspaceReader,
		r.gitSnapshotter,
		r.snapshotsRepo,
		r.workspaceWriter,
		r.executorProvider,
		r.eventsSender,
	); err != nil {
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

func (r *WorkspaceRootResolver) LandWorkspaceChange(ctx context.Context, args resolvers.LandWorkspaceArgs) (resolvers.WorkspaceResolver, error) {
	ws, err := r.workspaceReader.Get(string(args.Input.WorkspaceID))
	if err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to get workspace: %w", err))
	}

	if err := r.authService.CanWrite(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	var diffOpts []vcs.DiffOption
	if args.Input.DiffMaxSize > 0 {
		diffOpts = append(diffOpts, vcs.WithGitMaxSize(args.Input.DiffMaxSize))
	}

	if _, err := r.workspaceService.LandChange(ctx, ws, args.Input.PatchIDs, diffOpts...); err != nil {
		return nil, gqlerrors.Error(fmt.Errorf("failed to land change: %w", err))
	}

	return r.Workspace(ctx, resolvers.WorkspaceArgs{ID: args.Input.WorkspaceID})
}
