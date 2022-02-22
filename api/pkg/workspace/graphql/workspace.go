package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	service_auth "getsturdy.com/api/pkg/auth/service"
	"getsturdy.com/api/pkg/change"
	service_change "getsturdy.com/api/pkg/change/service"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	db_comments "getsturdy.com/api/pkg/comments/db"
	"getsturdy.com/api/pkg/events"
	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	service_suggestions "getsturdy.com/api/pkg/suggestions/service"
	db_view "getsturdy.com/api/pkg/view/db"
	"getsturdy.com/api/pkg/workspace"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	service_workspace "getsturdy.com/api/pkg/workspace/service"
	"getsturdy.com/api/pkg/workspace/vcs"
	vcsvcs "getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"

	"github.com/graph-gophers/graphql-go"
	"go.uber.org/zap"
)

type WorkspaceRootResolver struct {
	workspaceReader db_workspace.WorkspaceReader
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
	workspaceActivityRootResolver resolvers.WorkspaceActivityRootResolver
	reviewRootResolver            resolvers.ReviewRootResolver
	presenceRootResolver          resolvers.PresenceRootResolver
	suggestionRootResolver        resolvers.SuggestionRootResolver
	statusRootResolver            resolvers.StatusesRootResolver
	workspaceWatcherRootResolver  resolvers.WorkspaceWatcherRootResolver

	suggestionsService *service_suggestions.Service
	workspaceService   service_workspace.Service
	authService        *service_auth.Service
	changeService      *service_change.Service

	logger           *zap.Logger
	viewEvents       events.EventReadWriter
	workspaceWriter  db_workspace.WorkspaceWriter
	executorProvider executor.Provider
	eventsSender     events.EventSender
	gitSnapshotter   snapshotter.Snapshotter
}

func NewResolver(
	workspaceReader db_workspace.WorkspaceReader,
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
	workspaceActivityRootResolver resolvers.WorkspaceActivityRootResolver,
	reviewRootResolver resolvers.ReviewRootResolver,
	presenceRootResolver resolvers.PresenceRootResolver,
	suggestionRootResolver resolvers.SuggestionRootResolver,
	statusRootResolver resolvers.StatusesRootResolver,
	workspaceWatcherRootResolver resolvers.WorkspaceWatcherRootResolver,

	suggestionsService *service_suggestions.Service,
	workspaceService service_workspace.Service,
	authService *service_auth.Service,
	changeService *service_change.Service,

	logger *zap.Logger,
	viewEventsWriter events.EventReadWriter,
	workspaceWriter db_workspace.WorkspaceWriter,
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
		return nil, gqlerrors.ErrNotFound
	}

	if err := r.authService.CanRead(ctx, ws); err != nil {
		return nil, gqlerrors.Error(err)
	}

	return &WorkspaceResolver{w: ws, root: r}, nil
}

func (r *WorkspaceRootResolver) InternalWorkspace(ws *workspace.Workspace) resolvers.WorkspaceResolver {
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

type WorkspaceResolver struct {
	w    *workspace.Workspace
	root *WorkspaceRootResolver

	hasConflicts     bool
	hasConflictsErr  error
	hasConflictsOnce sync.Once
}

func (r *WorkspaceResolver) ID() graphql.ID {
	return graphql.ID(r.w.ID)
}

func (r *WorkspaceResolver) Name() string {
	return r.w.NameOrFallback()
}

func (r *WorkspaceResolver) Codebase(ctx context.Context) (resolvers.CodebaseResolver, error) {
	id := graphql.ID(r.w.CodebaseID)
	cb, err := r.root.codebaseResolver.Codebase(ctx, resolvers.CodebaseArgs{ID: &id})
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return cb, nil
}

func (r *WorkspaceResolver) Author(ctx context.Context) (resolvers.AuthorResolver, error) {
	author, err := r.root.authorResolver.Author(ctx, graphql.ID(r.w.UserID))
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	return author, nil
}

func (r *WorkspaceResolver) CreatedAt() int32 {
	return int32(r.w.CreatedAt.Unix())
}

func (r *WorkspaceResolver) LastLandedAt() *int32 {
	if r.w.LastLandedAt == nil {
		return nil
	}
	t := int32(r.w.LastLandedAt.Unix())
	return &t
}

func (r *WorkspaceResolver) ArchivedAt() *int32 {
	if r.w.ArchivedAt == nil {
		return nil
	}
	t := int32(r.w.ArchivedAt.Unix())
	return &t
}

func (r *WorkspaceResolver) UnarchivedAt() *int32 {
	if r.w.UnarchivedAt == nil {
		return nil
	}
	t := int32(r.w.UnarchivedAt.Unix())
	return &t
}

func (r *WorkspaceResolver) UpdatedAt() *int32 {
	if r.w.UpdatedAt == nil {
		return nil
	}
	t := int32(r.w.UpdatedAt.Unix())
	return &t
}

func (r *WorkspaceResolver) LastActivityAt() int32 {
	var largestTime int32

	maybeTime := []*time.Time{
		r.w.CreatedAt,
		r.w.LastLandedAt,
		r.w.ArchivedAt,
		r.w.UnarchivedAt,
		r.w.UpdatedAt,
	}

	for _, t := range maybeTime {
		if t == nil {
			continue
		}
		t2 := int32(t.Unix())
		if t2 > largestTime {
			largestTime = t2
		}
	}

	return largestTime
}

func (r *WorkspaceResolver) DraftDescription() string {
	return r.w.DraftDescription
}

func (r *WorkspaceResolver) View(ctx context.Context) (resolvers.ViewResolver, error) {
	if r.w.ViewID == nil {
		return nil, nil
	}
	return r.root.viewResolver.View(ctx, resolvers.ViewArgs{ID: graphql.ID(*r.w.ViewID)})
}

func (r *WorkspaceResolver) Comments() ([]resolvers.TopCommentResolver, error) {
	comments, err := r.root.commentResolver.InternalWorkspaceComments(r.w)
	switch {
	case err == nil:
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	default:
		return nil, gqlerrors.Error(err)
	}

	var res []resolvers.TopCommentResolver
	for _, comment := range comments {
		if topComment, ok := comment.ToTopComment(); ok {
			res = append(res, topComment)
		}
	}
	return res, nil
}

func (r *WorkspaceResolver) GitHubPullRequest(ctx context.Context) (resolvers.GitHubPullRequestResolver, error) {
	id := graphql.ID(r.w.ID)
	pr, err := r.root.prResolver.InternalGitHubPullRequestByWorkspaceID(ctx, resolvers.GitHubPullRequestArgs{WorkspaceID: &id})
	switch {
	case err == nil:
		return pr, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	default:
		return nil, gqlerrors.Error(err)
	}
}

func (r WorkspaceResolver) UpToDateWithTrunk(ctx context.Context) (bool, error) {
	if err := r.updateIsUpToDateWithTrunk(ctx); err != nil {
		return false, gqlerrors.Error(err)
	}
	if r.w.UpToDateWithTrunk == nil {
		return false, gqlerrors.Error(fmt.Errorf("failed to calculate UpToDateWithTrunk"))
	}
	return *r.w.UpToDateWithTrunk, nil
}

func (r *WorkspaceResolver) updateIsUpToDateWithTrunk(ctx context.Context) error {
	// We have a cached result, don't do anything
	if r.w.UpToDateWithTrunk != nil {
		return nil
	}

	var upToDate bool
	err := r.root.executorProvider.New().GitRead(func(repo vcsvcs.RepoGitReader) error {
		// Recalculate
		var err error
		upToDate, err = vcs.UpToDateWithTrunk(repo, r.w.ID)
		if err != nil {
			return fmt.Errorf("failed to check if workspace is up to date with trunk: %w", err)
		}
		return nil
	}).ExecTrunk(r.w.CodebaseID, "updateIsUpToDateWithTrunk")
	if err != nil {
		return err
	}

	// Fetch a new version of the workspace, and perform the update
	// TODO: Wrap all workspace mutations in a lock?
	wsForUpdates, err := r.root.workspaceReader.Get(r.w.ID)
	if err != nil {
		return err
	}

	wsForUpdates.UpToDateWithTrunk = &upToDate

	// Save updated cache
	if err := r.root.workspaceWriter.Update(ctx, wsForUpdates); err != nil {
		return err
	}

	// Also update the cached version of the workspace that we have in memory
	r.w.UpToDateWithTrunk = wsForUpdates.UpToDateWithTrunk

	return nil
}

func (r *WorkspaceResolver) Conflicts(ctx context.Context) (bool, error) {
	r.hasConflictsOnce.Do(func() {
		r.hasConflicts, r.hasConflictsErr = r.root.workspaceService.HasConflicts(ctx, r.w)
	})
	return r.hasConflicts, gqlerrors.Error(r.hasConflictsErr)
}

func (r *WorkspaceResolver) HeadChange(ctx context.Context) (resolvers.ChangeResolver, error) {
	// Recalculate head change
	if !r.w.HeadChangeComputed {
		var headCommitID string

		err := r.root.executorProvider.New().GitRead(func(repo vcsvcs.RepoGitReader) error {
			var err error
			headCommitID, err = repo.BranchCommitID(r.w.ID)
			if err != nil {
				return fmt.Errorf("could not get head commit from git: %w", err)
			}
			return nil
		}).ExecTrunk(r.w.CodebaseID, "workspaceHeadChange")
		if err != nil {
			return nil, gqlerrors.Error(err)
		}
		var newHeadChangeID *change.ID

		ch, err := r.root.changeService.GetByCommitAndCodebase(ctx, headCommitID, r.w.CodebaseID)
		switch {
		case errors.Is(err, sql.ErrNoRows), errors.Is(err, service_change.ErrNotFound):
			// change not found (could be the root commit, etc), hide it
			newHeadChangeID = nil
		case err != nil:
			return nil, gqlerrors.Error(fmt.Errorf("could not get change by commit: %w", err))
		default:
			newHeadChangeID = &ch.ID
		}

		// Fetch a new version of the workspace, and perform the update
		// TODO: Wrap all workspace mutations in a lock?
		wsForUpdates, err := r.root.workspaceReader.Get(r.w.ID)
		if err != nil {
			return nil, gqlerrors.Error(err)
		}

		wsForUpdates.HeadChangeComputed = true
		wsForUpdates.HeadChangeID = newHeadChangeID

		// Save updated cache
		if err := r.root.workspaceWriter.Update(ctx, wsForUpdates); err != nil {
			return nil, gqlerrors.Error(err)
		}

		// Also update the cached version of the workspace that we have in memory
		r.w.HeadChangeComputed = wsForUpdates.HeadChangeComputed
		r.w.HeadChangeID = newHeadChangeID

		r.root.logger.Info("recalculated head change", zap.String("workspace_id", r.w.ID), zap.Stringer("head", r.w.HeadChangeID))
	}

	if r.w.HeadChangeID == nil || !r.w.HeadChangeComputed {
		return nil, nil
	}

	cid := graphql.ID(r.w.CodebaseID)
	changeID := graphql.ID(*r.w.HeadChangeID)

	resolver, err := r.root.changeResolver.Change(ctx, resolvers.ChangeArgs{
		ID:         &changeID,
		CodebaseID: &cid,
	})
	switch {
	case err == nil:
		return resolver, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	default:
		return nil, err
	}
}

func (r *WorkspaceResolver) Activity(ctx context.Context, args resolvers.WorkspaceActivityArgs) ([]resolvers.WorkspaceActivityResolver, error) {
	return r.root.workspaceActivityRootResolver.InternalActivityByWorkspace(ctx, r.w.ID, args)
}

func (r *WorkspaceResolver) Reviews(ctx context.Context) ([]resolvers.ReviewResolver, error) {
	res, err := r.root.reviewRootResolver.InternalReviews(ctx, r.w.ID)
	switch {
	case err == nil:
		return res, nil
	case errors.Is(err, gqlerrors.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}

func (r *WorkspaceResolver) Presence(ctx context.Context) ([]resolvers.PresenceResolver, error) {
	return r.root.presenceRootResolver.InternalWorkspacePresence(ctx, r.w.ID)
}

func (r *WorkspaceResolver) Suggestion(ctx context.Context) (resolvers.SuggestionResolver, error) {
	suggestion, err := r.root.suggestionsService.GetByWorkspaceID(ctx, r.w.ID)
	switch {
	case err == nil:
		return r.root.suggestionRootResolver.InternalSuggestion(ctx, suggestion)
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	default:
		return nil, gqlerrors.Error(err)
	}
}

func (r *WorkspaceResolver) Suggestions(ctx context.Context) ([]resolvers.SuggestionResolver, error) {
	ss, err := r.root.suggestionsService.ListForWorkspaceID(ctx, r.w.ID)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}

	rr := make([]resolvers.SuggestionResolver, 0, len(ss))
	for _, s := range ss {
		r, err := r.root.suggestionRootResolver.InternalSuggestion(ctx, s)
		if err != nil {
			return nil, err
		}
		rr = append(rr, r)
	}

	return rr, nil
}

func (r *WorkspaceResolver) Statuses(ctx context.Context) ([]resolvers.StatusResolver, error) {
	if r.w.LatestSnapshotID == nil {
		return nil, nil
	}

	lastSnapshot, err := r.root.snapshotsRepo.Get(*r.w.LatestSnapshotID)
	switch {
	case err == nil:
		return r.root.statusRootResolver.InteralStatusesByCodebaseIDAndCommitID(ctx, lastSnapshot.CodebaseID, lastSnapshot.CommitID)
	case errors.Is(err, sql.ErrNoRows):
		return nil, nil
	default:
		return nil, gqlerrors.Error(err)
	}
}

func (r *WorkspaceResolver) Watchers(ctx context.Context) ([]resolvers.WorkspaceWatcherResolver, error) {
	return r.root.workspaceWatcherRootResolver.InternalWorkspaceWatchers(ctx, r.w)
}

func (r *WorkspaceResolver) SuggestingViews() []resolvers.ViewResolver {
	return nil
}
