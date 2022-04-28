package graphql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/graph-gophers/graphql-go"

	gqlerrors "getsturdy.com/api/pkg/graphql/errors"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/workspaces"
	"getsturdy.com/api/pkg/workspaces/db"
	service_workspace "getsturdy.com/api/pkg/workspaces/service"
	"getsturdy.com/api/pkg/workspaces/vcs"
	vcsvcs "getsturdy.com/api/vcs"
)

type WorkspaceResolver struct {
	w    *workspaces.Workspace
	root *WorkspaceRootResolver

	hasConflicts     bool
	hasConflictsErr  error
	hasConflictsOnce sync.Once

	latestSnapshot     *snapshots.Snapshot
	latestSnapshotErr  error
	latestSnapshotOnce sync.Once
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
	if r.w.CreatedAt == nil {
		return 0
	}
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

func (r *WorkspaceResolver) CommentsCount(ctx context.Context) (int32, error) {
	c, err := r.root.commentResolver.InternalCountByWorkspaceID(ctx, r.w.ID)
	if err != nil {
		return 0, gqlerrors.Error(err)
	}
	return c, nil
}

func (r *WorkspaceResolver) DiffsCount(_ context.Context) *int32 {
	return r.w.DiffsCount
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

func (r *WorkspaceResolver) UpToDateWithTrunk(ctx context.Context) (bool, error) {
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

	// Save updated cache
	if err := r.root.workspaceWriter.UpdateFields(ctx, r.w.ID, db.SetUpToDateWithTrunk(&upToDate)); err != nil {
		return err
	}

	// Also update the cached version of the workspace that we have in memory
	r.w.UpToDateWithTrunk = &upToDate

	return nil
}

func (r *WorkspaceResolver) Conflicts(ctx context.Context) (bool, error) {
	r.hasConflictsOnce.Do(func() {
		r.hasConflicts, r.hasConflictsErr = r.root.workspaceService.HasConflicts(ctx, r.w)
	})
	return r.hasConflicts, gqlerrors.Error(r.hasConflictsErr)
}

func (r *WorkspaceResolver) HeadChange(ctx context.Context) (resolvers.ChangeResolver, error) {
	headChange, err := r.root.workspaceService.HeadChange(ctx, r.w)
	switch {
	case errors.Is(err, service_workspace.ErrNotFound):
		return nil, nil
	case err != nil:
		return nil, gqlerrors.Error(err)
	}

	cid := graphql.ID(r.w.CodebaseID)
	changeID := graphql.ID(headChange.ID)

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

func (r *WorkspaceResolver) Activity(ctx context.Context, args resolvers.ActivityArgs) ([]resolvers.ActivityResolver, error) {
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

func (r *WorkspaceResolver) getLatestSnapshot() (*snapshots.Snapshot, error) {
	r.latestSnapshotOnce.Do(func() {
		if r.w.LatestSnapshotID == nil {
			return
		}
		snaphot, err := r.root.snapshotsRepo.Get(*r.w.LatestSnapshotID)
		r.latestSnapshot = snaphot
		r.latestSnapshotErr = err
	})
	return r.latestSnapshot, r.latestSnapshotErr
}

func (r *WorkspaceResolver) Statuses(ctx context.Context) ([]resolvers.WorkspaceStatusResolver, error) {
	return r.root.statusRootResolver.InternalWorkspaceStatuses(ctx, r.w.ID)
}

func (r *WorkspaceResolver) Watchers(ctx context.Context) ([]resolvers.WorkspaceWatcherResolver, error) {
	return r.root.workspaceWatcherRootResolver.InternalWorkspaceWatchers(ctx, r.w)
}

func (r *WorkspaceResolver) SuggestingViews() []resolvers.ViewResolver {
	return nil
}

func (r *WorkspaceResolver) Diffs(ctx context.Context) ([]resolvers.FileDiffResolver, error) {
	diffs, err := r.diffs(ctx)
	if err != nil {
		return nil, gqlerrors.Error(err)
	}
	res := make([]resolvers.FileDiffResolver, len(diffs))
	for k, diff := range diffs {
		res[k] = r.root.fileDiffRootResolver.InternalFileDiffWithWorkspace(r.w.ID, &diff, r.w)
	}
	return res, nil
}

func (r *WorkspaceResolver) diffs(ctx context.Context) ([]unidiff.FileDiff, error) {
	allower, err := r.root.authService.GetAllower(ctx, r.w)
	if err != nil {
		return nil, fmt.Errorf("failed to get allowed patterns: %w", err)
	}

	suggestion, err := r.root.suggestionsService.GetByWorkspaceID(ctx, r.w.ID)
	switch {
	case err == nil:
		diffs, err := r.root.suggestionsService.Diffs(ctx, suggestion, unidiff.WithAllower(allower))
		if err != nil {
			return nil, fmt.Errorf("failed to get diffs from suggestion: %w", err)
		}
		return diffs, nil
	case errors.Is(err, sql.ErrNoRows):
		diffs, isConflicting, err := r.root.workspaceService.Diffs(ctx, r.w.ID, service_workspace.WithAllower(allower))
		if err != nil {
			return nil, fmt.Errorf("failed to get diffs from workspace: %w", err)
		}
		if isConflicting {
			// TODO
			return diffs, nil
		}
		return diffs, nil
	default:
		return nil, fmt.Errorf("failed to get suggestion: %w", err)
	}
}

func (r *WorkspaceResolver) Change(ctx context.Context) (resolvers.ChangeResolver, error) {
	if r.w.ChangeID == nil {
		return nil, nil
	}
	chID := graphql.ID(*r.w.ChangeID)
	return r.root.changeResolver.Change(ctx, resolvers.ChangeArgs{
		ID: &chID,
	})
}

func (r *WorkspaceResolver) RebaseStatus(ctx context.Context) (resolvers.RebaseStatusResolver, error) {
	return r.root.rebaseStatusRootResolver.InternalWorkspaceRebaseStatus(ctx, r.w.ID)
}

func (r *WorkspaceResolver) DownloadTarGz(ctx context.Context) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.root.downloadsResolver.InternalWorkspaceDownloadTarGzUrl(ctx, r.w)
}

func (r *WorkspaceResolver) DownloadZip(ctx context.Context) (resolvers.ContentsDownloadUrlResolver, error) {
	return r.root.downloadsResolver.InternalWorkspaceDownloadZipUrl(ctx, r.w)
}

func (r *WorkspaceResolver) CanUndo(ctx context.Context) (bool, error) {

}

func (r *WorkspaceResolver) CanRedo(ctx context.Context) (bool, error) {

}
