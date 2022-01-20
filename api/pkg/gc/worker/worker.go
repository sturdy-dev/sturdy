package worker

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"mash/vcs"

	"mash/pkg/gc"
	"mash/pkg/gc/db"
	"mash/pkg/queue"
	"mash/pkg/queue/names"
	"mash/pkg/snapshots"
	db_snapshots "mash/pkg/snapshots/db"
	service_suggestion "mash/pkg/suggestions/service"
	db_view "mash/pkg/view/db"
	db_workspace "mash/pkg/workspace/db"
	"mash/vcs/executor"

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

			if err := work(
				context.Background(),
				q.gcRepo,
				*m,
				q.viewRepo,
				logger,
				q.snapshotsRepo,
				q.workspaceReader,
				q.suggestionService,
				q.executorProvider,
			); err != nil {
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

func gcSnapshots(
	ctx context.Context,
	m CodebaseGarbageCollectionQueueEntry,
	logger *zap.Logger,
	snapshotsRepo db_snapshots.Repository,
	workspaceReader db_workspace.WorkspaceReader,
	suggestionService *service_suggestion.Service,
	executorProvider executor.Provider,
) error {
	// Activate only for the Sturdy codebase
	if m.CodebaseID != "31596772-e9d6-445e-8144-856a3022744b" {
		return nil
	}

	// GC unused snapshots
	snapshots, err := snapshotsRepo.ListUndeletedInCodebase(m.CodebaseID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("could not get snapshots: %w", err)
	}

	logger.Info("cleaning up snapshots", zap.Int("total_snapshots", len(snapshots)))

	// Delete snapshots older than 3 hours
	threshold := time.Now().Add(-3 * time.Hour)

	for _, snapshot := range snapshots {
		logger := logger.With(zap.String("snapshot_id", snapshot.ID))

		if err := gcSnapshot(
			ctx,
			snapshot,
			threshold,
			m,
			logger,
			snapshotsRepo,
			workspaceReader,
			suggestionService,
			executorProvider,
		); err != nil {
			logger.Error("failed to gc snapshot", zap.Error(err))
			// do not fail
		}
	}

	return nil
}

func isSnapshotUsedAsSuggestion(
	ctx context.Context,
	snapshot *snapshots.Snapshot,
	suggestionService *service_suggestion.Service,
) (bool, error) {
	if snapshot.WorkspaceID == nil {
		return false, nil
	}

	// if there is a suggestion for this snapshot, it is used
	ss, err := suggestionService.ListBySnapshotID(ctx, snapshot.ID)
	if err != nil {
		return false, fmt.Errorf("could not get suggestions: %w", err)
	}
	if len(ss) > 0 {
		return true, nil
	}

	s, err := suggestionService.GetByWorkspaceID(ctx, *snapshot.WorkspaceID)
	switch {
	case err == nil:
		return s.ForSnapshotID == snapshot.ID, nil
	case errors.Is(err, sql.ErrNoRows):
		return false, nil
	default:
		return false, fmt.Errorf("failed to list suggestions: %w", err)
	}
}

func gcSnapshot(
	ctx context.Context,
	snapshot *snapshots.Snapshot,
	threshold time.Time,
	m CodebaseGarbageCollectionQueueEntry,
	logger *zap.Logger,
	snapshotsRepo db_snapshots.Repository,
	workspaceReader db_workspace.WorkspaceReader,
	suggestionService *service_suggestion.Service,
	executorProvider executor.Provider,
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

	partOfSuggestion, err := isSnapshotUsedAsSuggestion(ctx, snapshot, suggestionService)
	if err != nil {
		return fmt.Errorf("failed to calculate if snapshot is a part of suggestion: %w", err)
	}

	if partOfSuggestion {
		logger.Info("snapshot is a part of an open suggestion, skipping")
		return nil
	}

	snapshotBranchName := "snapshot-" + snapshot.ID

	if ws, err := workspaceReader.GetBySnapshotID(snapshot.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("could not get workspace by snapshot: %w", err)
	} else if err == nil && !ws.IsArchived() {
		logger.Info("snapshot is in use by non-archived workspace, skipping", zap.String("workspace_id", ws.ID))
		return nil
	}

	// Throttle heavy operations
	time.Sleep(time.Second / 2)

	logger.Info("deleting snapshot")

	// Delete branch on trunk
	if err := executorProvider.New().Git(func(trunkRepo vcs.Repo) error {
		if err := trunkRepo.DeleteBranch(snapshotBranchName); err != nil {
			return fmt.Errorf("failed to delete snapshot branch on trunk: %w", err)
		}

		logger.Info("trunk branch deleted", zap.String("branch_name", snapshotBranchName))

		return nil
	}).ExecTrunk(m.CodebaseID, "deleteTrunkSnapshot"); err != nil {
		logger.Error("failed to delete snapshot on trunk", zap.Error(err))
		// do not fail
		return nil
	}

	// Delete branch on the view that created the snapshot
	if snapshot.ViewID != "" {
		if err := executorProvider.New().
			AllowRebasingState(). // allowed to enable branch deletion even if the view is currently rebasing
			Git(func(viewGitRepo vcs.Repo) error {
				if err := viewGitRepo.DeleteBranch(snapshotBranchName); err != nil {
					return fmt.Errorf("failed to delete snapshot branch from view: %w", err)
				}

				logger.Info("view branch deleted", zap.String("branch_name", snapshotBranchName), zap.String("view_id", snapshot.ViewID))

				return nil
			}).ExecView(m.CodebaseID, snapshot.ViewID, "deleteViewSnapshot"); err != nil {
			logger.Error("failed to open view", zap.Error(err))
			return nil
		}
	}

	t := time.Now()
	snapshot.DeletedAt = &t
	if err := snapshotsRepo.Update(snapshot); err != nil {
		return fmt.Errorf("failed to mark snapshot as deleted: %w", err)
	}

	return nil
}

func getGCInterval(m CodebaseGarbageCollectionQueueEntry) time.Duration {
	return time.Hour
}

func work(
	ctx context.Context,
	gcRepo db.Repository,
	m CodebaseGarbageCollectionQueueEntry,
	viewRepo db_view.Repository,
	logger *zap.Logger,
	snapshotsRepo db_snapshots.Repository,
	workspaceReader db_workspace.WorkspaceReader,
	suggestionService *service_suggestion.Service,
	executorProvider executor.Provider,
) error {
	t0 := time.Now()

	gcInterval := getGCInterval(m)
	// Skip if recently run
	entries, err := gcRepo.ListSince(ctx, m.CodebaseID, t0.Add(-1*gcInterval))
	if err != nil {
		return fmt.Errorf("failed to get last runs: %w", err)

	}
	if len(entries) > 0 {
		logger.Sugar().Infof("skipping gc ran in the last %s", gcInterval)
		return nil
	}

	logger.Info("starting gc")

	if err := gcSnapshots(
		ctx,
		m,
		logger,
		snapshotsRepo,
		workspaceReader,
		suggestionService,
		executorProvider,
	); err != nil {
		logger.Error("failed to gc snapshots", zap.Error(err))
		// do not fail
	}

	if err := executorProvider.New().Git(func(trunkRepo vcs.Repo) error {
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
	views, err := viewRepo.ListByCodebase(m.CodebaseID)
	if err != nil {
		return err
	}

	for _, view := range views {
		logger := logger.With(zap.String("view_id", view.ID))

		if err := executorProvider.New().Git(func(viewGitRepo vcs.Repo) error {
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
	if err := gcRepo.Create(ctx, &gc.CodebaseGarbageStatus{
		CodebaseID:     m.CodebaseID,
		CompletedAt:    now,
		DurationMillis: now.Sub(t0).Milliseconds(),
	}); err != nil {
		return fmt.Errorf("failed to record gc run stats: %w", err)
	}

	return nil
}
