package worker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/queue"
	"getsturdy.com/api/pkg/queue/names"
	"getsturdy.com/api/pkg/snapshots"
	"getsturdy.com/api/pkg/snapshots/snapshotter"

	"go.uber.org/zap"
)

type Queue interface {
	Enqueue(ctx context.Context, codebaseID, viewID, workspaceID string, paths []string, action snapshots.Action) error
	Start(ctx context.Context) error
}

type q struct {
	logger *zap.Logger
	queue  queue.Queue
	name   names.IncompleteQueueName

	snapshotter snapshotter.Snapshotter
}

func New(
	logger *zap.Logger,
	queue queue.Queue,
	snapshotter snapshotter.Snapshotter,
) Queue {
	return &q{
		logger:      logger.Named("snapshotterQueue"),
		queue:       queue,
		name:        names.ViewSnapshot,
		snapshotter: snapshotter,
	}
}

func (q *q) Enqueue(ctx context.Context, codebaseID, viewID, workspaceID string, paths []string, action snapshots.Action) error {
	if err := q.queue.Publish(ctx, q.name, &SnapshotQueueEntry{
		CodebaseID:   codebaseID,
		ViewID:       viewID,
		WorkspaceID:  workspaceID,
		ChangedFiles: paths,
		Action:       action,
	}); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

type SnapshotQueueEntry struct {
	CodebaseID   string           `json:"codebase_id"`
	ViewID       string           `json:"view_id"`
	WorkspaceID  string           `json:"workspace_id"`
	ChangedFiles []string         `json:"changed_files"`
	Action       snapshots.Action `json:"action"`
}

func (q *q) Start(ctx context.Context) error {
	messages := make(chan queue.Message)
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				q.logger.Error("panic in runner", zap.String("panic", fmt.Sprintf("%v", rec)))
			}
		}()

		for msg := range messages {
			t0 := time.Now()

			m := &SnapshotQueueEntry{}
			if err := msg.As(m); err != nil {
				q.logger.Error("failed to read message", zap.Error(err))
				continue
			}

			logger := q.logger.With(
				zap.String("view_id", m.ViewID),
				zap.String("codebase_id", m.CodebaseID),
				zap.String("workspace_id", m.WorkspaceID),
				zap.Stringer("action", m.Action),
			)

			if _, err := q.snapshotter.Snapshot(
				m.CodebaseID,
				m.WorkspaceID,
				m.Action,
				snapshotter.WithPaths(m.ChangedFiles),
				snapshotter.WithOnView(m.ViewID),
			); errors.Is(err, snapshotter.ErrCantSnapshotRebasing) {
				logger.Warn("failed to make snapshot", zap.Error(err))
				continue
			} else if errors.Is(err, snapshotter.ErrCantSnapshotWrongBranch) {
				logger.Warn("failed to make snapshot", zap.Error(err))
				continue
			} else if err != nil {
				logger.Error("failed to make snapshot", zap.Error(err))
				continue
			}

			if err := msg.Ack(); err != nil {
				logger.Error("failed to ack message", zap.Error(err))
				continue
			}

			logger.Info("created snapshot", zap.Duration("duration", time.Since(t0)))
		}
	}()

	q.logger.Info("starting queue", zap.Stringer("queue_name", q.name))
	if err := q.queue.Subscribe(ctx, q.name, messages); err != nil {
		return fmt.Errorf("could not subscribe to queue: %w", err)
	}
	q.logger.Info("queue stoped", zap.Stringer("queue_name", q.name))

	return nil
}
