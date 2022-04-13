package service

import (
	"context"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/codebases"
	"getsturdy.com/api/pkg/queue"
	"getsturdy.com/api/pkg/queue/names"
	"getsturdy.com/api/pkg/users"

	"go.uber.org/zap"
)

type PullRequestImportEvent struct {
	CodebaseID codebases.ID `json:"codebase_id"`
	UserID     users.ID     `json:"user_id"`
}

type ImporterQueue struct {
	logger        *zap.Logger
	queue         queue.Queue
	name          names.IncompleteQueueName
	gitHubService *Service
}

func NewImporterQueue(
	logger *zap.Logger,
	queue queue.Queue,
) *ImporterQueue {
	return &ImporterQueue{
		logger: logger.Named("GitHubPullRequestImporterWorker"),
		queue:  queue,
		name:   names.CodebaseGitHubPullRequestImporter,
	}
}

func (q *ImporterQueue) setService(svc *Service) {
	q.gitHubService = svc
}

func (q *ImporterQueue) Enqueue(ctx context.Context, codebaseID codebases.ID, userID users.ID) error {
	if err := q.queue.Publish(ctx, q.name, &PullRequestImportEvent{
		CodebaseID: codebaseID,
		UserID:     userID,
	}); err != nil {
		return fmt.Errorf("failed to publish event to queue: %w", err)
	}
	return nil
}

func (q *ImporterQueue) Start(ctx context.Context) error {
	messages := make(chan queue.Message)
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				q.logger.Error("panic in runner", zap.String("panic", fmt.Sprintf("%v", rec)))
			}
		}()

		for msg := range messages {
			t0 := time.Now()

			event := &PullRequestImportEvent{}
			if err := msg.As(event); err != nil {
				q.logger.Error("failed to decode message", zap.Error(err))
				continue
			}

			logger := q.logger.With(zap.Stringer("codebase_id", event.CodebaseID), zap.Stringer("user_id", event.UserID))

			if err := q.gitHubService.ImportOpenPullRequestsByUser(ctx, event.CodebaseID, event.UserID); err != nil {
				logger.Error("failed to import pull request", zap.Error(err))
				// No return, ack message
			}

			if err := msg.Ack(); err != nil {
				logger.Error("failed to ack message", zap.Error(err))
				continue
			}

			logger.Info("finished", zap.Duration("duration", time.Since(t0)))
		}
	}()

	q.logger.Info("starting queue", zap.Stringer("queue_name", q.name))
	if err := q.queue.Subscribe(ctx, q.name, messages); err != nil {
		return fmt.Errorf("could not subscribe to queue: %w", err)
	}
	q.logger.Info("queue stoped", zap.Stringer("queue_name", q.name))

	return nil
}

type nopImporter struct{}

func (*nopImporter) Enqueue(context.Context, string, string) error {
	return nil
}

func NopImporter() *nopImporter {
	return &nopImporter{}
}

func (*nopImporter) Start(context.Context) error {
	return nil
}
