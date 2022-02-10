package worker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"getsturdy.com/api/vcs"

	"getsturdy.com/api/pkg/gc"
	"getsturdy.com/api/pkg/gc/db"
	"getsturdy.com/api/pkg/queue"
	"getsturdy.com/api/pkg/queue/names"
	"getsturdy.com/api/pkg/snapshots"
	db_snapshots "getsturdy.com/api/pkg/snapshots/db"
	service_suggestion "getsturdy.com/api/pkg/suggestions/service"
	db_view "getsturdy.com/api/pkg/view/db"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	"getsturdy.com/api/vcs/executor"

	"go.uber.org/zap"
)

type CodebaseGarbageCollectionQueueEntry struct {
	CodebaseID string `json:"codebase_id"`
}

type Queue struct {
	logger *zap.Logger
	queue  queue.Queue
	name   names.IncompleteQueueName

	gcRepo            db.Repository
	viewRepo          db_view.Repository
	snapshotsRepo     db_snapshots.Repository
	workspaceReader   db_workspace.WorkspaceReader
	suggestionService *service_suggestion.Service
	executorProvider  executor.Provider
}

func New(
	logger *zap.Logger,
	queue queue.Queue,
	gcRepo db.Repository,
	viewRepo db_view.Repository,
	snapshotsRepo db_snapshots.Repository,
	workspaceReader db_workspace.WorkspaceReader,
	suggestionService *service_suggestion.Service,
	executorProvider executor.Provider,
) *Queue {
	return &Queue{
		logger:            logger.Named("gcRunnerQueue"),
		queue:             queue,
		name:              names.CodebaseGarbageCollection,
		gcRepo:            gcRepo,
		viewRepo:          viewRepo,
		snapshotsRepo:     snapshotsRepo,
		workspaceReader:   workspaceReader,
		suggestionService: suggestionService,
		executorProvider:  executorProvider,
	}
}

func (q *Queue) Enqueue(ctx context.Context, codebaseID string) error {
	if err := q.queue.Publish(ctx, q.name, &CodebaseGarbageCollectionQueueEntry{
		CodebaseID: codebaseID,
	}); err != nil {
		return fmt.Errorf("could not publish to queue: %w", err)
	}
	return nil
}

func (q *Queue) Start(ctx context.Context) error {
	messages := make(chan queue.Message)
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				q.logger.Error("panic in runner", zap.String("panic", fmt.Sprintf("%v", rec)))
			}
		}()

		for msg := range messages {
			t0 := time.Now()

			m := &CodebaseGarbageCollectionQueueEntry{}
			if err := msg.As(m); err != nil {
				q.logger.Error("failed to decode message", zap.Error(err))
				continue
			}
			logger := q.logger.With(zap.String("codebase_id", m.CodebaseID))

			if err := q.work(context.Background(), *m, logger); err != nil {
				logger.Error("failed to gc codebase", zap.Error(err))
				continue
			}

			if err := msg.Ack(); err != nil {
				logger.Error("failed to ack message", zap.Error(err))
				continue
			}

			logger.Info("gc ran", zap.Duration("duration", time.Since(t0)))
		}
	}()

	q.logger.Info("starting queue", zap.Stringer("queue_name", q.name))
	if err := q.queue.Subscribe(ctx, q.name, messages); err != nil {
		return fmt.Errorf("could not subscribe to queue: %w", err)
	}
	q.logger.Info("queue stoped", zap.Stringer("queue_name", q.name))

	return nil
}

func (q *Queue) gcSnapshots(ctx context.Context, m CodebaseGarbageCollectionQueueEntry) error {
	// GC unused snapshots
	snapshots, err := q.snapshotsRepo.ListUndeletedInCodebase(m.CodebaseID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("could not get snapshots: %w", err)
	}

	q.logger.Info("cleaning up snapshots", zap.Int("total_snapshots", len(snapshots)))

	// Delete snapshots older than 3 hours
	threshold := time.Now().Add(-3 * time.Hour)

	for _, snapshot := range snapshots {
		logger := q.logger.With(zap.String("snapshot_id", snapshot.ID))

		if err := q.gcSnapshot(
			ctx,
			snapshot,
			threshold,
			m,
			logger,
		); err != nil {
			logger.Error("failed to gc snapshot", zap.Error(err))
			// do not fail
		}
	}

	return nil
}

func (q *Queue) isSnapshotUsedAsSuggestion(ctx context.Context, snapshot *snapshots.Snapshot) (bool, error) {
	if snapshot.WorkspaceID == nil {
		return false, nil
	}

	// if there is a suggestion for this snapshot, it is used
	ss, err := q.suggestionService.ListBySnapshotID(ctx, snapshot.ID)
	if err != nil {
		return false, fmt.Errorf("could not get suggestions: %w", err)
	}
	if len(ss) > 0 {
		return true, nil
	}

	s, err := q.suggestionService.GetByWorkspaceID(ctx, *snapshot.WorkspaceID)
	switch {
	case err == nil:
		return s.ForSnapshotID == snapshot.ID, nil
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	default:
		return false, fmt.Errorf("failed to list suggestions: %w", err)
	}
}

func (q *Queue) gcSnapshot(
	ctx context.Context,
	snapshot *snapshots.Snapshot,
	threshold time.Time,
	m CodebaseGarbageCollectionQueueEntry,
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

	partOfSuggestion, err := q.isSnapshotUsedAsSuggestion(ctx, snapshot)
	if err != nil {
		return fmt.Errorf("failed to calculate if snapshot is a part of suggestion: %w", err)
	}

	if partOfSuggestion {
		logger.Info("snapshot is a part of an open suggestion, skipping")
		return nil
	}

	// Throttle heavy operations
	time.Sleep(time.Second / 2)

	if err := q.deleteSnapshotBranch(logger, snapshot); err != nil {
		return fmt.Errorf("failed to delete snapshot id=%s: %w", snapshot.ID, err)
	}

	t := time.Now()
	snapshot.DeletedAt = &t
	if err := q.snapshotsRepo.Update(snapshot); err != nil {
		return fmt.Errorf("failed to mark snapshot as deleted: %w", err)
	}

	return nil
}

func (q *Queue) deleteSnapshotBranch(logger *zap.Logger, snapshot *snapshots.Snapshot) error {
	logger.Info("deleting snapshot")

	if ws, err := q.workspaceReader.GetBySnapshotID(snapshot.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("could not get workspace by snapshot: %w", err)
	} else if err == nil && !ws.IsArchived() {
		logger.Info("snapshot is in use by non-archived workspace, skipping", zap.String("workspace_id", ws.ID))
		return nil
	}

	snapshotBranchName := "snapshot-" + snapshot.ID

	// Delete branch on trunk
	if err := q.executorProvider.New().GitWrite(func(trunkRepo vcs.RepoGitWriter) error {
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

	// Delete branch on the view that created the snapshot
	if snapshot.ViewID != "" && !strings.HasPrefix(snapshot.ViewID, "tmp-") {
		if err := q.executorProvider.New().
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
	}

	return nil
}

func getGCInterval() time.Duration {
	return time.Hour
}

func (q *Queue) work(
	ctx context.Context,
	m CodebaseGarbageCollectionQueueEntry,
	logger *zap.Logger,
) error {
	t0 := time.Now()

	gcInterval := getGCInterval()
	// Skip if recently run
	entries, err := q.gcRepo.ListSince(ctx, m.CodebaseID, t0.Add(-1*gcInterval))
	if err != nil {
		return fmt.Errorf("failed to get last runs: %w", err)

	}
	if len(entries) > 0 {
		logger.Sugar().Infof("skipping gc ran in the last %s", gcInterval)
		return nil
	}

	logger.Info("starting gc")

	if err := q.gcSnapshots(ctx, m); err != nil {
		logger.Error("failed to gc snapshots", zap.Error(err))
		// do not fail
	}

	if err := q.executorProvider.New().GitWrite(func(trunkRepo vcs.RepoGitWriter) error {
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
	}).ExecTrunk(m.CodebaseID, "gcTrunk"); err != nil {
		logger.Error("failed to git gc trunk", zap.Error(err))
		// don't exit
	}

	// gc all views
	views, err := q.viewRepo.ListByCodebase(m.CodebaseID)
	if err != nil {
		return err
	}

	for _, view := range views {
		logger := logger.With(zap.String("view_id", view.ID))

		if err := q.executorProvider.New().GitWrite(func(viewGitRepo vcs.RepoGitWriter) error {
			if err := viewGitRepo.GitReflogExpire(); err != nil {
				logger.Error("failed to run git-reflog expire on trunk", zap.Error(err))
				// don't exit
			}

			if err := viewGitRepo.GitGC(); err != nil {
				logger.Error("failed to run git-gc on view", zap.Error(err))
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
	if err := q.gcRepo.Create(ctx, &gc.CodebaseGarbageStatus{
		CodebaseID:     m.CodebaseID,
		CompletedAt:    now,
		DurationMillis: now.Sub(t0).Milliseconds(),
	}); err != nil {
		return fmt.Errorf("failed to record gc run stats: %w", err)
	}

	return nil
}
