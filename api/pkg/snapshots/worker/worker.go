package worker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/queue"
	"getsturdy.com/api/pkg/queue/names"
	"getsturdy.com/api/pkg/snapshots"
	service_snapshots "getsturdy.com/api/pkg/snapshots/service"
	"getsturdy.com/api/pkg/users"
	service_users "getsturdy.com/api/pkg/users/service"

	"go.uber.org/zap"
)

type Queue interface {
	Enqueue(ctx context.Context, codebaseID codebases.ID, viewID, workspaceID string, userID users.ID, action snapshots.Action) error
	Start(ctx context.Context) error
}

type q struct {
	logger *zap.Logger
	queue  queue.Queue
	name   names.IncompleteQueueName

	snapshotter *service_snapshots.Service
	userService service_users.Service
}

func New(
	logger *zap.Logger,
	queue queue.Queue,
	snapshotter *service_snapshots.Service,
	userService service_users.Service,
) Queue {
	return &q{
		logger:      logger.Named("snapshotterQueue"),
		queue:       queue,
		name:        names.ViewSnapshot,
		snapshotter: snapshotter,
		userService: userService,
	}
}

func (q *q) Enqueue(ctx context.Context, codebaseID codebases.ID, viewID, workspaceID string, userID users.ID, action snapshots.Action) error {
	if err := q.queue.Publish(ctx, q.name, &SnapshotQueueEntry{
		CodebaseID:  codebaseID,
		ViewID:      viewID,
		WorkspaceID: workspaceID,
		UserID:      &userID,
		Action:      action,
	}); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

type SnapshotQueueEntry struct {
	CodebaseID  codebases.ID     `json:"codebase_id"`
	ViewID      string           `json:"view_id"`
	UserID      *users.ID        `json:"user_id"` // nullable for backwards compatability
	WorkspaceID string           `json:"workspace_id"`
	Action      snapshots.Action `json:"action"`
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
				zap.Stringer("codebase_id", m.CodebaseID),
				zap.String("workspace_id", m.WorkspaceID),
				zap.Stringer("action", m.Action),
			)

			var options []service_snapshots.SnapshotOption
			options = append(options, service_snapshots.WithOnView(m.ViewID))

			if m.UserID != nil {
				if user, err := q.userService.GetByID(ctx, *m.UserID); err != nil {
					logger.Error("failed to get user", zap.Error(err))
					// don't fail
				} else {
					options = append(options, service_snapshots.WithUser(user))
				}
			}

			if _, err := q.snapshotter.Snapshot(
				m.CodebaseID,
				m.WorkspaceID,
				m.Action,
				options...,
			); errors.Is(err, service_snapshots.ErrCantSnapshotRebasing) {
				logger.Warn("failed to make snapshot", zap.Error(err))
				continue
			} else if errors.Is(err, service_snapshots.ErrCantSnapshotWrongBranch) {
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
