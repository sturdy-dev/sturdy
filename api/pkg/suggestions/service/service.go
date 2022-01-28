package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/notification"
	sender_notification "getsturdy.com/api/pkg/notification/sender"
	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/pkg/snapshots/snapshotter"
	"getsturdy.com/api/pkg/suggestions"
	db_suggestions "getsturdy.com/api/pkg/suggestions/db"
	"getsturdy.com/api/pkg/unidiff"
	"getsturdy.com/api/pkg/events"
	vcs_view "getsturdy.com/api/pkg/view/vcs"
	"getsturdy.com/api/pkg/workspace"
	service_workspace "getsturdy.com/api/pkg/workspace/service"
	vcs_workspace "getsturdy.com/api/pkg/workspace/vcs"
	"getsturdy.com/api/vcs"
	"getsturdy.com/api/vcs/executor"
	"getsturdy.com/api/vcs/provider"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger

	suggestionRepo db_suggestions.Repository

	workspaceService service_workspace.Service

	executorProvider   executor.Provider
	snapshotter        snapshotter.Snapshotter
	analyticsClient    analytics.Client
	notificationSender sender_notification.NotificationSender
	eventSender        events.EventSender
}

func New(
	logger *zap.Logger,
	suggestionRepo db_suggestions.Repository,
	workspaceService service_workspace.Service,
	executorProvider executor.Provider,
	snapshotter snapshotter.Snapshotter,
	analyticsClient analytics.Client,
	notificationSender sender_notification.NotificationSender,
	eventSender events.EventSender,
) *Service {
	return &Service{
		logger: logger,

		suggestionRepo: suggestionRepo,

		workspaceService: workspaceService,

		executorProvider:   executorProvider,
		snapshotter:        snapshotter,
		analyticsClient:    analyticsClient,
		notificationSender: notificationSender,
		eventSender:        eventSender,
	}
}

// todo: move this to the workspace service
func (s *Service) copyWorkspace(ctx context.Context, userID string, ws *workspace.Workspace) (*workspace.Workspace, error) {
	changeID := ""
	if ws.HeadCommitID != nil {
		changeID = *ws.HeadCommitID
	}

	name := ""
	if ws.Name != nil {
		name = fmt.Sprintf("Suggestions: %s", *ws.Name)
	}

	createRequest := service_workspace.CreateWorkspaceRequest{
		UserID:     userID,
		CodebaseID: ws.CodebaseID,
		Name:       name,
		ChangeID:   changeID,
	}

	newWorkspace, err := s.workspaceService.Create(createRequest)
	if err != nil {
		return nil, fmt.Errorf("faliled to create a workspace: %w", err)
	}

	return newWorkspace, nil
}

// RecordActivity sends notifications and resurrects existing suggestions.
func (s *Service) RecordActivity(ctx context.Context, workspaceID string) error {
	suggestion, err := s.GetByWorkspaceID(ctx, workspaceID)
	switch {
	case err == nil:
	case errors.Is(err, sql.ErrNoRows):
		return nil
	default:
		return fmt.Errorf("failed to get suggestion: %w", err)
	}

	forWorkspace, err := s.workspaceService.GetByID(ctx, suggestion.ForWorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	// if the user hasn't been notified yet, notify them
	shouldNotify := suggestion.NotifiedAt == nil
	if shouldNotify {
		if err := s.notificationSender.User(ctx, forWorkspace.UserID, forWorkspace.CodebaseID, notification.NewSuggestionNotificationType, string(suggestion.ID)); err != nil {
			s.logger.Error("failed to send notification", zap.Error(err))
		}
		now := time.Now()
		suggestion.NotifiedAt = &now
	}

	// resurrect the suggestion if dismissed
	shouldResurrect := suggestion.DismissedAt != nil
	if shouldResurrect {
		suggestion.DismissedAt = nil
	}

	shouldUpdate := shouldNotify || shouldResurrect
	if shouldUpdate {
		if err := s.suggestionRepo.Update(ctx, suggestion); err != nil {
			return fmt.Errorf("failed to update suggestion: %w", err)
		}

		if err := s.eventSender.Workspace(suggestion.ForWorkspaceID, events.WorkspaceUpdatedSuggestion, string(suggestion.ID)); err != nil {
			s.logger.Error("failed to send event", zap.Error(err))
		}
	}

	return nil
}

func (s *Service) Create(ctx context.Context, userID string, forWorkspace *workspace.Workspace) (*suggestions.Suggestion, error) {
	if forWorkspace.LatestSnapshotID == nil {
		return nil, fmt.Errorf("workspace has no snapshot")
	}

	ws, err := s.copyWorkspace(ctx, userID, forWorkspace)
	if err != nil {
		return nil, fmt.Errorf("failed to copy workspace: %w", err)
	}

	if err := s.workspaceService.CopyPatches(ctx, ws, forWorkspace); err != nil {
		return nil, fmt.Errorf("failed to copy patches: %w", err)
	}

	suggestion := &suggestions.Suggestion{
		ID:             suggestions.ID(uuid.NewString()),
		CodebaseID:     ws.CodebaseID,
		WorkspaceID:    ws.ID,
		ForSnapshotID:  *forWorkspace.LatestSnapshotID,
		ForWorkspaceID: forWorkspace.ID,
		UserID:         userID,
		CreatedAt:      time.Now(),
	}
	if err := s.suggestionRepo.Create(ctx, suggestion); err != nil {
		return nil, fmt.Errorf("failed to create: %w", err)
	}

	if err := s.analyticsClient.Enqueue(analytics.Capture{
		DistinctId: suggestion.UserID,
		Event:      "suggestions-create",
		Properties: analytics.NewProperties().
			Set("workspace_id", forWorkspace.ID).
			Set("suggestion_id", suggestion.ID),
	}); err != nil {
		s.logger.Error("failed to send analytics event", zap.Error(err))
	}

	return suggestion, nil
}

// GetByID returns a suggestion by id.
func (s *Service) GetByID(ctx context.Context, id suggestions.ID) (*suggestions.Suggestion, error) {
	return s.suggestionRepo.GetByID(ctx, id)
}

// GetByWorkspaceID returns a suggestion that is made from the workspaceID.
func (s *Service) GetByWorkspaceID(ctx context.Context, workspaceID string) (*suggestions.Suggestion, error) {
	return s.suggestionRepo.GetByWorkspaceID(ctx, workspaceID)
}

func (s *Service) ListBySnapshotID(ctx context.Context, snapshotID string) ([]*suggestions.Suggestion, error) {
	return s.suggestionRepo.ListBySnapshotID(ctx, snapshotID)
}

// ListForWorkspaceID return a list of currently opened suggestions for the workspace.
func (s *Service) ListForWorkspaceID(ctx context.Context, forWorkspaceID string) ([]*suggestions.Suggestion, error) {
	ss, err := s.suggestionRepo.ListForWorkspaceID(ctx, forWorkspaceID)
	if err != nil {
		return nil, err
	}

	activeSuggestions := make([]*suggestions.Suggestion, 0, len(ss))
	for _, suggestion := range ss {
		ws, err := s.workspaceService.GetByID(ctx, suggestion.WorkspaceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get workspace: %w", err)
		}
		if !ws.IsArchived() {
			activeSuggestions = append(activeSuggestions, suggestion)
		}
	}

	return activeSuggestions, nil
}

// Dismiss marks the suggestion as dismissed.
func (s *Service) Dismiss(ctx context.Context, suggestion *suggestions.Suggestion) error {
	now := time.Now()
	suggestion.DismissedAt = &now
	if err := s.suggestionRepo.Update(ctx, suggestion); err != nil {
		return fmt.Errorf("failed to update suggestion: %w", err)
	}
	return nil
}

// ApplyHunks applies the suggested hunks to the workspace.
func (s *Service) ApplyHunks(ctx context.Context, suggestion *suggestions.Suggestion, hunkIDs ...string) error {
	if len(hunkIDs) == 0 {
		return nil
	}

	originalWorkspace, err := s.workspaceService.GetByID(ctx, suggestion.ForWorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to get original workspace: %w", err)
	}

	fileDiffs, err := s.diffs(ctx, suggestion, originalWorkspace)
	if err != nil {
		return fmt.Errorf("failed to get diffs: %w", err)
	}

	toApply := make(map[string]bool, len(hunkIDs))
	for _, id := range hunkIDs {
		toApply[id] = true
	}

	patches := [][]byte{}
	appliedHunks := []string{}
	for _, fd := range fileDiffs {
		for hunkIndex, hunk := range fd.Hunks {
			if !toApply[hunk.ID] {
				continue
			}

			patches = append(patches, []byte(hunk.Patch))
			appliedHunks = append(appliedHunks, (&suggestions.Hunk{
				FileName: fd.PreferredName,
				Index:    hunkIndex,
			}).String())
		}
	}

	if originalWorkspace.ViewID == nil { // apply patches to the snapshot
		if err := s.executorProvider.New().Schedule(func(repoProvider provider.RepoProvider) error {
			repo, cancel, err := vcs_view.TemporaryViewFromSnapshot(repoProvider, originalWorkspace.CodebaseID, originalWorkspace.ID, *originalWorkspace.LatestSnapshotID)
			if err != nil {
				return fmt.Errorf("failed to create temporary view: %w", err)
			}
			defer func() {
				if err := cancel(); err != nil {
					s.logger.Error("failed to cleanup temporary view", zap.Error(err))
				}
			}()

			if err := repo.ApplyPatchesToWorkdir(patches); err != nil {
				return fmt.Errorf("failed to apply patches: %w", err)
			}

			if _, err := s.snapshotter.Snapshot(
				originalWorkspace.CodebaseID,
				originalWorkspace.ID,
				snapshots.ActionSuggestionApply,
				snapshotter.WithOnView(*repo.ViewID()),
				snapshotter.WithOnRepo(repo),
				snapshotter.WithMarkAsLatestInWorkspace(),
			); err != nil {
				return fmt.Errorf("failed to snapshot: %w", err)
			}

			return nil
		}).ExecTrunk(originalWorkspace.CodebaseID, "applySuggestionDiffs"); err != nil {
			return fmt.Errorf("failed to apply patches: %w", err)
		}
	} else { // apply to the view
		if err := s.executorProvider.New().Write(func(repo vcs.RepoWriter) error {
			return repo.ApplyPatchesToWorkdir(patches)
		}).ExecView(originalWorkspace.CodebaseID, *originalWorkspace.ViewID, "applySuggestionDiffs"); err != nil {
			return fmt.Errorf("failed to apply patches: %w", err)
		}
	}

	suggestion.AppliedHunks = append(suggestion.AppliedHunks, appliedHunks...)
	if err := s.suggestionRepo.Update(ctx, suggestion); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	if err := s.analyticsClient.Enqueue(analytics.Capture{
		DistinctId: originalWorkspace.UserID,
		Event:      "suggestions-apply",
		Properties: analytics.NewProperties().
			Set("workspace_id", originalWorkspace.ID).
			Set("suggestion_id", suggestion.ID),
	}); err != nil {
		s.logger.Error("failed to send analytics event", zap.Error(err))
	}

	return nil
}

// DismissHunks marks suggeted hunks as dismissed.
func (s *Service) DismissHunks(ctx context.Context, suggestion *suggestions.Suggestion, hunkIDs ...string) error {
	if len(hunkIDs) == 0 {
		return nil
	}

	originalWorkspace, err := s.workspaceService.GetByID(ctx, suggestion.ForWorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to get original workspace: %w", err)
	}

	fileDiffs, err := s.diffs(ctx, suggestion, originalWorkspace)
	if err != nil {
		return fmt.Errorf("failed to get diffs: %w", err)
	}

	toDismiss := make(map[string]bool, len(hunkIDs))
	for _, id := range hunkIDs {
		toDismiss[id] = true
	}
	dismissedHunks := []string{}
	for _, fd := range fileDiffs {
		for hunkIndex, hunk := range fd.Hunks {
			if !toDismiss[hunk.ID] {
				continue
			}
			dismissedHunks = append(dismissedHunks, (&suggestions.Hunk{
				FileName: fd.PreferredName,
				Index:    hunkIndex,
			}).String())
		}
	}

	suggestion.DismissedHunks = append(suggestion.DismissedHunks, dismissedHunks...)
	if err := s.suggestionRepo.Update(ctx, suggestion); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	if err := s.analyticsClient.Enqueue(analytics.Capture{
		DistinctId: originalWorkspace.UserID,
		Event:      "suggestions-dismiss",
		Properties: analytics.NewProperties().
			Set("workspace_id", originalWorkspace.ID).
			Set("suggestion_id", suggestion.ID),
	}); err != nil {
		s.logger.Error("failed to send analytics event", zap.Error(err))
	}

	return nil
}

func (s *Service) RemovePatches(ctx context.Context, suggestion *suggestions.Suggestion, patchIDs ...string) error {
	if len(patchIDs) == 0 {
		return nil
	}

	workspace, err := s.workspaceService.GetByID(ctx, suggestion.WorkspaceID)
	if err != nil {
		return fmt.Errorf("failed to get workspace: %w", err)
	}

	diffs, err := s.Diffs(ctx, suggestion)
	if err != nil {
		return fmt.Errorf("failed to get diffs: %w", err)
	}
	patches := [][]byte{}
	for _, diff := range diffs {
		for _, hunk := range diff.Hunks {
			patches = append(patches, []byte(hunk.Patch))
		}
	}
	removePatches := vcs_workspace.RemoveWithPatches(s.logger, patches, patchIDs...)

	if workspace.ViewID != nil {
		if err := s.executorProvider.New().Write(func(repo vcs.RepoWriter) error {
			if err := removePatches(repo); err != nil {
				return err
			}
			if _, err := s.snapshotter.Snapshot(
				workspace.CodebaseID,
				workspace.ID,
				snapshots.ActionFileUndoPatch,
				snapshotter.WithOnRepo(repo),
				snapshotter.WithOnView(*workspace.ViewID),
				snapshotter.WithMarkAsLatestInWorkspace(),
			); err != nil {
				return fmt.Errorf("failed to snapshot: %w", err)
			}
			return nil
		}).ExecView(workspace.CodebaseID, *workspace.ViewID, "removeSuggestionPatches"); err != nil {
			return fmt.Errorf("failed to apply patches: %w", err)
		}
		return nil
	}

	if workspace.LatestSnapshotID != nil {
		if err := s.executorProvider.New().Schedule(func(repoProvider provider.RepoProvider) error {
			repo, cancel, err := vcs_view.TemporaryViewFromSnapshot(repoProvider, workspace.CodebaseID, workspace.ID, *workspace.LatestSnapshotID)
			if err != nil {
				return fmt.Errorf("failed to create temporary view: %w", err)
			}
			defer func() {
				if err := cancel(); err != nil {
					s.logger.Error("failed to cleanup temporary view", zap.Error(err))
				}
			}()

			if err := removePatches(repo); err != nil {
				return err
			}
			if _, err := s.snapshotter.Snapshot(
				workspace.CodebaseID,
				workspace.ID,
				snapshots.ActionFileUndoPatch,
				snapshotter.WithOnRepo(repo),
				snapshotter.WithOnView(*workspace.ViewID),
				snapshotter.WithMarkAsLatestInWorkspace(),
			); err != nil {
				return fmt.Errorf("failed to snapshot: %w", err)
			}

			return nil
		}).ExecTrunk(workspace.CodebaseID, "removeSuggestionPatches"); err != nil {
			return fmt.Errorf("failed to apply patches: %w", err)
		}
		return nil
	}

	return fmt.Errorf("workspace has no view nor latest snapshot")
}

// Diffs returns all the diffs of the suggestion as viewed by the suggestion.ForWorkspace.
func (s *Service) Diffs(ctx context.Context, suggestion *suggestions.Suggestion, oo ...unidiff.Option) ([]unidiff.FileDiff, error) {
	originalWorkspace, err := s.workspaceService.GetByID(ctx, suggestion.ForWorkspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get original workspace: %w", err)
	}
	return s.diffs(ctx, suggestion, originalWorkspace, oo...)
}

func (s *Service) diffs(
	ctx context.Context,
	suggestion *suggestions.Suggestion,
	originalWorkspace *workspace.Workspace,
	oo ...unidiff.Option,
) ([]unidiff.FileDiff, error) {
	suggestingWorkspace, err := s.workspaceService.GetByID(ctx, suggestion.WorkspaceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggesting workspace: %w", err)
	}

	if suggestingWorkspace.LatestSnapshotID == nil {
		return nil, nil
	}

	suggestingSnapshot, err := s.snapshotter.GetByID(ctx, *suggestingWorkspace.LatestSnapshotID)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggesting snapshot: %w", err)
	}

	baseSnapshot, err := s.snapshotter.GetByID(ctx, suggestion.ForSnapshotID)
	if err != nil {
		return nil, fmt.Errorf("failed to get base snapshot: %w", err)
	}

	var diffs []unidiff.FileDiff
	if err := s.executorProvider.New().Git(func(repo vcs.Repo) error {
		gitDiffs, err := repo.DiffCommits(baseSnapshot.CommitID, suggestingSnapshot.CommitID)
		if err != nil {
			return fmt.Errorf("failed to get diffs: %w", err)
		}

		differ := unidiff.NewUnidiff(unidiff.NewGitPatchReader(gitDiffs), s.logger, oo...).
			WithExpandedHunks().
			WithIgnoreBinary()

		hunkifiedDiffs, err := differ.Decorate()
		if err != nil {
			return fmt.Errorf("failed to decorate diffs: %w", err)
		}
		diffs = hunkifiedDiffs
		return nil
	}).ExecTrunk(suggestingWorkspace.CodebaseID, "snapshotDiffs"); err != nil {
		return nil, fmt.Errorf("failed to schedule repo on trunk: %w", err)
	}

	appliedHunks := make([]*suggestions.Hunk, 0, len(suggestion.AppliedHunks))
	for _, ah := range suggestion.AppliedHunks {
		if a, err := suggestions.ParseAppliedHunkID(ah); err == nil {
			appliedHunks = append(appliedHunks, a)
		} else {
			return nil, fmt.Errorf("couldn't parse applied hunk id: %w", err)
		}
	}

	dismissedHunkIDs := make([]*suggestions.Hunk, 0, len(suggestion.DismissedHunks))
	for _, ah := range suggestion.DismissedHunks {
		if a, err := suggestions.ParseAppliedHunkID(ah); err == nil {
			dismissedHunkIDs = append(dismissedHunkIDs, a)
		} else {
			return nil, fmt.Errorf("couldn't parse applied hunk id: %w", err)
		}
	}

	// todo: decrease complexity
	// mark applied and dismissed hunks
	for _, fd := range diffs {
		for _, appliedHunk := range appliedHunks {
			if appliedHunk.FileName == fd.PreferredName && len(fd.Hunks) > appliedHunk.Index {
				fd.Hunks[appliedHunk.Index].IsApplied = true
			}
		}
		for _, dismissedHunk := range dismissedHunkIDs {
			if dismissedHunk.FileName == fd.PreferredName && len(fd.Hunks) > dismissedHunk.Index {
				fd.Hunks[dismissedHunk.Index].IsDismissed = true
			}
		}
	}

	if originalWorkspace.ViewID == nil {
		return diffs, nil
	}

	// mark outdated hunks
	if err := s.executorProvider.New().Read(func(repo vcs.RepoReader) error {
		for _, fd := range diffs {
			for hunkIndex, hunk := range fd.Hunks {
				if hunk.IsApplied || hunk.IsDismissed {
					continue
				}

				canApply, err := repo.CanApplyPatch([]byte(hunk.Patch))
				if err != nil {
					return fmt.Errorf("can not check if patch can be applied: %w", err)
				}

				if !canApply {
					fd.Hunks[hunkIndex].IsOutdated = true
				}
			}
		}
		return nil
	}).ExecView(originalWorkspace.CodebaseID, *originalWorkspace.ViewID, "calculateOutdatedDiffs"); err != nil {
		return nil, fmt.Errorf("failed to calculate outdated diffs: %w", err)
	}

	return diffs, nil
}
