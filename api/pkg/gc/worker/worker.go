package worker

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/gc/service"
	"getsturdy.com/api/pkg/queue"
	"getsturdy.com/api/pkg/queue/names"
)

type CodebaseGarbageCollectionQueueEntry struct {
	CodebaseID codebases.ID `json:"codebase_id"`
}

type Queue struct {
	logger *zap.Logger
	queue  queue.Queue
	name   names.IncompleteQueueName

	service *service.Service
}

func New(
	logger *zap.Logger,
	queue queue.Queue,
	service *service.Service,
) *Queue {
	return &Queue{
		logger:  logger.Named("gcRunnerQueue"),
		queue:   queue,
		name:    names.CodebaseGarbageCollection,
		service: service,
	}
}

func (q *Queue) Enqueue(ctx context.Context, codebaseID codebases.ID) error {
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
			logger := q.logger.With(zap.Stringer("codebase_id", m.CodebaseID))

			if err := q.service.Work(context.Background(), logger, m.CodebaseID); err != nil {
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
