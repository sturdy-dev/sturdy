package workers

import (
	"context"
	"errors"
	"fmt"
	"time"

	service_github_webhooks "getsturdy.com/api/pkg/github/enterprise/webhooks"
	"getsturdy.com/api/pkg/queue"
	"getsturdy.com/api/pkg/queue/names"

	"go.uber.org/zap"
)

type WebhooksQueue struct {
	logger          *zap.Logger
	queue           queue.Queue
	name            names.IncompleteQueueName
	webhooksService *service_github_webhooks.Service
}

func NewWebhooksQueue(logger *zap.Logger, queue queue.Queue, webhooksService *service_github_webhooks.Service) *WebhooksQueue {
	return &WebhooksQueue{
		logger:          logger.Named("github webhooks queue"),
		queue:           queue,
		name:            names.GithubWebhooks,
		webhooksService: webhooksService,
	}
}

type WebhookEvent struct {
	Installation             *service_github_webhooks.InstallationEvent
	InstallationRepositories *service_github_webhooks.InstallationRepositoriesEvent
	Push                     *service_github_webhooks.PushEvent
	PullRequest              *service_github_webhooks.PullRequestEvent
	Status                   *service_github_webhooks.StatusEvent
	WorkflowJob              *service_github_webhooks.WorkflowJobEvent
}

func (q *WebhooksQueue) Enqueue(ctx context.Context, event *WebhookEvent) error {
	if err := q.queue.Publish(ctx, q.name, event); err != nil {
		return fmt.Errorf("failed to publish to queue: %w", err)
	}
	return nil
}

func (q *WebhooksQueue) Start(ctx context.Context) error {
	messages := make(chan queue.Message)
	go func() {
		defer func() {
			if rec := recover(); rec != nil {
				q.logger.Error("panic in runner", zap.String("panic", fmt.Sprintf("%v", rec)), zap.Stack("stack"))
			}
		}()
		for msg := range messages {
			t0 := time.Now()
			event := &WebhookEvent{}

			if err := msg.As(event); err != nil {
				q.logger.Error("failed to convert message to event", zap.Error(err))
				continue
			}
			logger := q.getLogger(event)
			logger.Info("processing")

			ctx := context.Background()
			if err := q.work(ctx, event); errors.Is(err, errUnknownType) {
				logger.Error("failed to process event", zap.Error(err), zap.Bool("will retry", false))
				// No return, ack message
			} else if err != nil {
				logger.Error("failed to process event", zap.Error(err), zap.Bool("will retry", true))
				// retry
				continue
			}

			if err := msg.Ack(); err != nil {
				logger.Error("failed to ack message", zap.Error(err))
				continue
			}

			logger.Info("done", zap.Duration("duration", time.Since(t0)))
		}
	}()

	q.logger.Info("starting queue", zap.Stringer("queue_name", q.name))
	if err := q.queue.Subscribe(ctx, q.name, messages); err != nil {
		return fmt.Errorf("could not subscribe to queue: %w", err)
	}
	q.logger.Info("queue stoped", zap.Stringer("queue_name", q.name))

	return nil
}

func (q *WebhooksQueue) getLogger(event *WebhookEvent) *zap.Logger {
	if event.Installation != nil {
		return q.logger.With(
			zap.String("event_type", "installation"),
			zap.Int64("installation_id", event.Installation.GetInstallation().GetID()),
		)
	} else if event.InstallationRepositories != nil {
		return q.logger.With(
			zap.String("event_type", "installation repositories"),
			zap.Int64("installation_id", event.InstallationRepositories.GetInstallation().GetID()),
		)
	} else if event.Push != nil {
		return q.logger.With(
			zap.String("event_type", "push"),
			zap.Int64("installation_id", event.Push.GetInstallation().GetID()),
			zap.String("repo", event.Push.GetRepo().GetFullName()),
		)
	} else if event.PullRequest != nil {
		return q.logger.With(
			zap.String("event_type", "pull request"),
			zap.Int64("installation_id", event.PullRequest.GetInstallation().GetID()),
			zap.String("repo", event.PullRequest.GetRepo().GetFullName()),
		)
	} else if event.Status != nil {
		return q.logger.With(
			zap.String("event_type", "status"),
			zap.Int64("installation_id", event.Status.GetInstallation().GetID()),
			zap.String("repo", event.Status.GetRepo().GetFullName()),
		)
	} else if event.WorkflowJob != nil {
		return q.logger.With(
			zap.String("event_type", "workflow job"),
			zap.Int64("installation_id", event.WorkflowJob.GetInstallation().GetID()),
			zap.String("repo", event.WorkflowJob.GetRepo().GetFullName()),
		)
	} else {
		return q.logger.With(
			zap.String("event_type", "unknown"),
		)
	}
}

var errUnknownType = fmt.Errorf("unknown event type")

func (q *WebhooksQueue) work(ctx context.Context, event *WebhookEvent) error {
	if event.Installation != nil {
		return q.webhooksService.HandleInstallationEvent(ctx, event.Installation)
	} else if event.InstallationRepositories != nil {
		return q.webhooksService.HandleInstallationRepositoriesEvent(ctx, event.InstallationRepositories)
	} else if event.Push != nil {
		return q.webhooksService.HandlePushEvent(ctx, event.Push)
	} else if event.PullRequest != nil {
		return q.webhooksService.HandlePullRequestEvent(ctx, event.PullRequest)
	} else if event.Status != nil {
		return q.webhooksService.HandleStatusEvent(ctx, event.Status)
	} else if event.WorkflowJob != nil {
		return q.webhooksService.HandleWorkflowJobEvent(ctx, event.WorkflowJob)
	} else {
		return errUnknownType
	}
}
