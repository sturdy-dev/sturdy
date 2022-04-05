package workers

import (
	"context"
	"fmt"

	"getsturdy.com/api/pkg/changes"
	service_ci "getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/queue"
	"getsturdy.com/api/pkg/queue/names"

	"go.uber.org/zap"
)

// BuildQueue is a background queue that triggers builds for enqueued changes.
type BuildQueue struct {
	logger *zap.Logger

	queue queue.Queue
	name  names.IncompleteQueueName

	ciService *service_ci.Service
}

func New(logger *zap.Logger, queue queue.Queue, ciService *service_ci.Service) *BuildQueue {
	return &BuildQueue{
		logger:    logger.Named("ciRunnerQueue"),
		queue:     queue,
		name:      names.CITriggerQueue,
		ciService: ciService,
	}
}

func (r *BuildQueue) EnqueueChange(ctx context.Context, ch *changes.Change) error {
	if err := r.queue.Publish(ctx, r.name, ch); err != nil {
		return fmt.Errorf("failed to publish to queue: %w", err)
	}
	return nil
}

// Start starts the runner.
func (r *BuildQueue) Start(ctx context.Context) error {
	messages := make(chan queue.Message)
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				r.logger.Error("panic in runner", zap.String("panic", fmt.Sprintf("%v", rec)))
			}
		}()
		for msg := range messages {
			ch := &changes.Change{}
			if err := msg.As(ch); err != nil {
				r.logger.Error("failed to decode message", zap.Error(err), zap.Any("message", msg))
				continue
			}

			if err := r.trigger(ctx, ch); err != nil {
				r.logger.Error("failed to build", zap.Error(err), zap.Any("message", msg))
				continue
			}

			if err := msg.Ack(); err != nil {
				r.logger.Error("failed to ack message", zap.Error(err), zap.Any("message", msg))
				continue
			}
		}
	}()

	r.logger.Info("starting queue", zap.Stringer("queue_name", r.name))
	if err := r.queue.Subscribe(ctx, r.name, messages); err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}
	r.logger.Info("queue stoped", zap.Stringer("queue_name", r.name))

	return nil
}

func (r *BuildQueue) trigger(ctx context.Context, ch *changes.Change) error {
	r.logger.Info(
		"trigger ci build",
		zap.String("change_id", string(ch.ID)),
		zap.Stringer("codebase_id", ch.CodebaseID),
	)

	if _, err := r.ciService.TriggerChange(ctx, ch); err != nil {
		return fmt.Errorf("failed to trigger change: %w", err)
	}

	return nil
}
