package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"getsturdy.com/api/vcs"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/gc"
	"getsturdy.com/api/pkg/gc/db"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	service_suggestion "getsturdy.com/api/pkg/suggestions/service"
	"getsturdy.com/api/pkg/view"
	db_view "getsturdy.com/api/pkg/view/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
	"getsturdy.com/api/vcs/executor"

	"go.uber.org/zap"
)

type Service struct {
	logger            *zap.Logger
	gcRepo            db.Repository
	viewRepo          db_view.Repository
	snapshotsRepo     db_snapshots.Repository
	workspaceReader   db_workspaces.WorkspaceReader
	suggestionService *service_suggestion.Service
	executorProvider  executor.Provider
}

func New(
	logger *zap.Logger,
	gcRepo db.Repository,
	viewRepo db_view.Repository,
	snapshotsRepo db_snapshots.Repository,
	workspaceReader db_workspaces.WorkspaceReader,
	suggestionService *service_suggestion.Service,
	executorProvider executor.Provider,
) *Service {
	return &Service{
		logger:            logger.Named("gcService"),
		gcRepo:            gcRepo,
		viewRepo:          viewRepo,
		snapshotsRepo:     snapshotsRepo,
		workspaceReader:   workspaceReader,
		suggestionService: suggestionService,
		executorProvider:  executorProvider,
	}
}

func mapSlice[In any, Out any](slice []In, mapper func(In) Out) []Out {
	out := make([]Out, len(slice))
	for i, v := range slice {
		out[i] = mapper(v)
	}
	return out
}

func filterSlice[T any](slice []T, filter func(T) bool) []T {
	out := make([]T, 0)
	for _, v := range slice {
		if filter(v) {
			out = append(out, v)
		}
	}
	return out
}

func (svc *Service) gcSnapshotsInView(ctx context.Context, view *view.View, snapshotThreshold time.Duration) error {
	allBranches := []string{}
	if err := svc.executorProvider.New().GitRead(func(repo vcs.RepoGitReader) error {
		branches, err := repo.Branches()
		if err != nil {
			return fmt.Errorf("failed to list branches: %w", err)
		}
		allBranches = branches
		return nil
	}).ExecView(view.CodebaseID, view.ID, "listBranches"); err != nil {
		return fmt.Errorf("failed to execute view: %w", err)
	}

	isSnapshotBranch := func(branch string) bool {
		return strings.HasPrefix(branch, "snapshot-")
	}

	snapshotBranches := filterSlice(allBranches, isSnapshotBranch)
	if len(snapshotBranches) == 0 {
		return nil
	}

	branchSnapshotID := func(branch string) string {
		return strings.TrimPrefix(branch, "snapshot-")
	}

	snapshotIDs := mapSlice(snapshotBranches, branchSnapshotID)
	snapshots, err := svc.snapshotsRepo.ListByIDs(ctx, snapshotIDs)
	if err != nil {
		return fmt.Errorf("failed to list snapshots: %w", err)
	}

	// Delete snapshots older than
	threshold := time.Now().Add(snapshotThreshold)
	for _, snapshot := range snapshots {
		logger := svc.logger.With(zap.String("snapshot_id", snapshot.ID), zap.String("view_id", view.ID))

		if err := svc.gcSnapshot(
			ctx,
			snapshot,
			view,
			threshold,
			logger,
		); err != nil {
			logger.Error("failed to gc snapshot", zap.Error(err))
			// do not fail
		}
	}

	return nil
}

func (svc *Service) gcSnapshots(ctx context.Context, codebaseID codebases.ID, snapshotThreshold time.Duration) error {
	views, err := svc.viewRepo.ListByCodebase(codebaseID)
	if err != nil {
		return fmt.Errorf("failed to list views: %w", err)
	}

	for _, view := range views {
		logger := svc.logger.With(zap.String("view_id", view.ID), zap.Stringer("codebase_id", view.CodebaseID))
		if err != svc.gcSnapshotsInView(ctx, view, snapshotThreshold) {
			logger.Error("failed to gc snapshots in view", zap.Error(err))
		}
	}

	return nil
}

func (svc *Service) isSnapshotUsedAsSuggestion(ctx context.Context, snapshot *snapshots.Snapshot) (bool, error) {
	// if there is a suggestion for this snapshot, it is used
	ss, err := svc.suggestionService.ListBySnapshotID(ctx, snapshot.ID)
	if err != nil {
		return false, fmt.Errorf("could not get suggestions: %w", err)
	}
	if len(ss) > 0 {
		return true, nil
	}

	s, err := svc.suggestionService.GetByWorkspaceID(ctx, snapshot.WorkspaceID)
	switch {
	case err == nil:
		return s.ForSnapshotID == snapshot.ID, nil
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	default:
		return false, fmt.Errorf("failed to list suggestions: %w", err)
	}
}

func (svc *Service) gcSnapshot(
	ctx context.Context,
	snapshot *snapshots.Snapshot,
	view *view.View,
	threshold time.Time,
	logger *zap.Logger,
) error {
	if snapshot.CreatedAt.After(threshold) {
		logger.Info(
			"snapshot too new, skipping",
			zap.Time("threshold", threshold),
			zap.Time("created_at", snapshot.CreatedAt),
		)
		return nil
	}

	if snapshot.DeletedAt != nil {
		logger.Info("snapshot is deleted, skipping")
		return nil
	}

	if ws, err := svc.workspaceReader.GetBySnapshotID(snapshot.ID); errors.Is(err, sql.ErrNoRows) {
		// continue
	} else if err != nil {
		return fmt.Errorf("could not get workspace: %w", err)
	} else if !ws.IsArchived() {
		logger.Info("snapshot is used by a workspace, skipping", zap.String("workspace_id", ws.ID))
		return nil
	}

	partOfSuggestion, err := svc.isSnapshotUsedAsSuggestion(ctx, snapshot)
	if err != nil {
		return fmt.Errorf("failed to calculate if snapshot is a part of suggestion: %w", err)
	}

	if partOfSuggestion {
		logger.Info("snapshot is a part of an open suggestion, skipping")
		return nil
	}

	// Throttle heavy operations
	time.Sleep(time.Second / 2)

	if err := svc.deleteSnapshotBranchInView(logger, view, snapshot); err != nil {
		return fmt.Errorf("failed to delete snapshot id=%s: %w", snapshot.ID, err)
	}

	t := time.Now()
	snapshot.DeletedAt = &t
	if err := svc.snapshotsRepo.Update(snapshot); err != nil {
		return fmt.Errorf("failed to mark snapshot as deleted: %w", err)
	}

	return nil
}

func (svc *Service) deleteSnapshotBranchInView(logger *zap.Logger, view *view.View, snapshot *snapshots.Snapshot) error {
	logger.Info("deleting snapshot")

	if ws, err := svc.workspaceReader.GetBySnapshotID(snapshot.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("could not get workspace by snapshot: %w", err)
	} else if err == nil && !ws.IsArchived() {
		logger.Info("snapshot is in use by non-archived workspace, skipping", zap.String("workspace_id", ws.ID))
		return nil
	}

	snapshotBranchName := snapshot.BranchName()

	// Delete branch on trunk
	if err := svc.executorProvider.New().GitWrite(func(trunkRepo vcs.RepoGitWriter) error {
		if err := trunkRepo.DeleteBranch(snapshotBranchName); err != nil {
			return fmt.Errorf("failed to delete snapshot branch on trunk: %w", err)
		}

		logger.Info("trunk branch deleted", zap.String("branch_name", snapshotBranchName))

		return nil
	}).ExecTrunk(snapshot.CodebaseID, "deleteTrunkSnapshot"); err != nil {
		logger.Error("failed to delete snapshot on trunk", zap.Error(err))
		// do not fail
		return nil
	}

	if err := svc.executorProvider.New().
		AllowRebasingState(). // allowed to enable branch deletion even if the view is currently rebasing
		GitWrite(func(viewGitRepo vcs.RepoGitWriter) error {
			if err := viewGitRepo.DeleteBranch(snapshotBranchName); err != nil {
				return fmt.Errorf("failed to delete snapshot branch from view: %w", err)
			}

			logger.Info("view branch deleted", zap.String("branch_name", snapshotBranchName), zap.String("view_id", snapshot.ViewID))

			return nil
		}).ExecView(snapshot.CodebaseID, snapshot.ViewID, "deleteViewSnapshot"); err != nil {
		logger.Error("failed to open view", zap.Error(err))
		return nil
	}

	return nil
}

func getGCInterval() time.Duration {
	return time.Hour
}

func getSnapshotThreshold() time.Duration {
	return -3 * time.Hour
}

func (svc *Service) Work(
	ctx context.Context,
	logger *zap.Logger,
	codebaseID codebases.ID,
) error {
	return svc.WorkWithOptions(ctx, logger, codebaseID, getGCInterval(), getSnapshotThreshold())
}

func (svc *Service) WorkWithOptions(
	ctx context.Context,
	logger *zap.Logger,
	codebaseID codebases.ID,
	gcInterval time.Duration,
	gcSnapshotsThreshold time.Duration,
) error {
	t0 := time.Now()

	// Skip if recently run
	entries, err := svc.gcRepo.ListSince(ctx, codebaseID, t0.Add(-1*gcInterval))
	if err != nil {
		return fmt.Errorf("failed to get last runs: %w", err)

	}
	if len(entries) > 0 {
		logger.Sugar().Infof("skipping gc ran in the last %s", gcInterval)
		return nil
	}

	logger.Info("starting gc")

	if err := svc.gcSnapshots(ctx, codebaseID, gcSnapshotsThreshold); err != nil {
		logger.Error("failed to gc snapshots", zap.Error(err))
		// do not fail
	}

	if err := svc.executorProvider.New().GitWrite(func(trunkRepo vcs.RepoGitWriter) error {
		if err := trunkRepo.GitReflogExpire(); err != nil {
			logger.Error("failed to run git-reflog expire on trunk", zap.Error(err))
			// don't exit
		}

		if err := trunkRepo.GitGC(); err != nil {
			logger.Error("failed to run git-gc on trunk", zap.Error(err))
			// don't exit
		}

		logger.Info("trunk cleaned up")

		return nil
	}).ExecTrunk(codebaseID, "gcTrunk"); err != nil {
		logger.Error("failed to git gc trunk", zap.Error(err))
		// don't exit
	}

	// gc all views
	views, err := svc.viewRepo.ListByCodebase(codebaseID)
	if err != nil {
		return err
	}

	for _, view := range views {
		logger := logger.With(zap.String("view_id", view.ID))

		if err := svc.executorProvider.New().GitWrite(func(viewGitRepo vcs.RepoGitWriter) error {
			if err := viewGitRepo.GitReflogExpire(); err != nil {
				logger.Error("failed to run git-reflog expire on trunk", zap.Error(err))
				// don't exit
			}

			if err := viewGitRepo.GitGC(); err != nil {
				logger.Error("failed to run git-gc on view", zap.Error(err))
				// don't exit
			}

			if err := viewGitRepo.GitRemotePrune("origin"); err != nil {
				logger.Error("failed to run git remote prune on view", zap.Error(err))
				// don't exit
			}

			logger.Info("view cleaned up")
			return nil
		}).ExecView(view.CodebaseID, view.ID, "gcView"); err != nil {
			// If the view is rebasing, it will be GC'd on the next run, no big deal.
			if errors.Is(err, executor.ErrIsRebasing) {
				logger.Warn("failed to run git gc on view", zap.Error(err))
			} else {
				logger.Error("failed to run git gc on view", zap.Error(err))
			}
		}
	}

	now := time.Now()
	if err := svc.gcRepo.Create(ctx, &gc.CodebaseGarbageStatus{
		CodebaseID:     codebaseID,
		CompletedAt:    now,
		DurationMillis: now.Sub(t0).Milliseconds(),
	}); err != nil {
		return fmt.Errorf("failed to record gc run stats: %w", err)
	}

	return nil
}
