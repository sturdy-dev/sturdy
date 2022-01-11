package workers

import (
	"context"
	"fmt"
	"mash/pkg/github"
	"time"

	service_github "mash/pkg/github/service"
	"mash/pkg/queue"
	"mash/pkg/queue/names"

	"go.uber.org/zap"
)

type ClonerQueue struct {
	logger        *zap.Logger
	queue         queue.Queue
	name          names.IncompleteQueueName
	gitHubService *service_github.Service
}

func NewClonerQueue(
	logger *zap.Logger,
	queue queue.Queue,
	githubService *service_github.Service,
) *ClonerQueue {
	return &ClonerQueue{
		logger:        logger.Named("GitHubClonerQueue"),
		queue:         queue,
		gitHubService: githubService,
		name:          names.CodebaseGitHubCloner,
	}
}

func (q *ClonerQueue) Enqueue(ctx context.Context, event *github.CloneRepositoryEvent) error {
	if err := q.queue.Publish(ctx, q.name, event); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}
	return nil
}

func (q *ClonerQueue) Start(ctx context.Context) error {
	messages := make(chan queue.Message, 10)
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				q.logger.Error("panic in runner", zap.String("panic", fmt.Sprintf("%v", rec)), zap.Stack("recovered"))
			}
		}()

		for msg := range messages {
			t0 := time.Now()

			event := &github.CloneRepositoryEvent{}
			if err := msg.As(event); err != nil {
				q.logger.Error("failed to parse codebase event in worker", zap.Error(err))
				continue
			}

			q.logger.Info("cloning", zap.String("codebase_id", event.CodebaseID))

			if err := q.gitHubService.Clone(
				event.CodebaseID,
				event.InstallationID,
				event.GitHubRepositoryID,
				event.SenderUserID,
			); err != nil {
				q.logger.Error("failed to clone", zap.Error(err))
				continue
			}

			if err := msg.Ack(); err != nil {
				q.logger.Error("failed to ack", zap.Error(err))
				continue
			}

			q.logger.Info("cloned", zap.String("codebase_id", event.CodebaseID), zap.Duration("duration", time.Since(t0)))
		}
	}()

	q.logger.Info("starting queue", zap.Stringer("queue_name", q.name))
	if err := q.queue.Subscribe(ctx, q.name, messages); err != nil {
		return fmt.Errorf("could not subscribe to queue: %w", err)
	}
	q.logger.Info("queue stoped", zap.Stringer("queue_name", q.name))

	return nil
}

type nopCLoner struct{}

func NopCloner() *nopCLoner {
	return &nopCLoner{}
}

func (n *nopCLoner) Enqueue(ctx context.Context, event *github.CloneRepositoryEvent) error {
	return nil
}
