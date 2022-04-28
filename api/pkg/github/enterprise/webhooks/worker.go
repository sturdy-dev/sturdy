package webhooks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"getsturdy.com/api/pkg/queue"
	"getsturdy.com/api/pkg/queue/names"

	"go.uber.org/zap"
)

type Queue struct {
	logger          *zap.Logger
	queue           queue.Queue
	name            names.IncompleteQueueName
	webhooksService *Service
}

func NewWebhooksQueue(logger *zap.Logger, queue queue.Queue, webhooksService *Service) *Queue {
	return &Queue{
		logger:          logger.Named("github webhooks queue"),
		queue:           queue,
		name:            names.GithubWebhooks,
		webhooksService: webhooksService,
	}
}

type WebhookEvent struct {
	Installation             *InstallationEvent
	InstallationRepositories *InstallationRepositoriesEvent
	Push                     *PushEvent
	PullRequest              *PullRequestEvent
	Status                   *StatusEvent
	WorkflowJob              *WorkflowJobEvent
}

func (q *Queue) Enqueue(ctx context.Context, event *WebhookEvent) error {
	if err := q.queue.Publish(ctx, q.name, event); err != nil {
		return fmt.Errorf("failed to publish to queue: %w", err)
	}
	return nil
}

func (q *Queue) Start(ctx context.Context) error {
	messages := make(chan queue.Message)

	work := func(threadNumber int) {
		logger := q.logger.With(zap.Int("thread", threadNumber))

		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("panic in runner", zap.String("panic", fmt.Sprintf("%v", rec)), zap.Stack("stack"))
			}
		}()

		for msg := range messages {
			t0 := time.Now()
			event := &WebhookEvent{}

			if err := msg.As(event); err != nil {
				q.logger.Error("failed to convert message to event", zap.Error(err))
				continue
			}
			
			logger := getLogger(logger, event)
			logger.Info("processing")

			ctx := context.Background()

			if err := q.work(ctx, event); err != nil {

				retryAllowed := event.Status != nil || event.WorkflowJob != nil || event.Installation != nil || event.InstallationRepositories != nil
				willRetry := retryAllowed && !errors.Is(err, errUnknownType)
				if willRetry {
					logger.Error("failed to process event", zap.Error(err), zap.Bool("will retry", true))
					// retry
					continue
				}

				logger.Error("failed to process event", zap.Error(err), zap.Bool("will retry", false))
				// No return, ack message
			}

			if err := msg.Ack(); err != nil {
				logger.Error("failed to ack message", zap.Error(err))
				continue
			}

			logger.Info("done", zap.Duration("duration", time.Since(t0)))
		}
	}

	// start 5 worker threads
	for i := 0; i < 5; i++ {
		i := i
		go work(i)
	}

	q.logger.Info("starting queue", zap.Stringer("queue_name", q.name))
	if err := q.queue.Subscribe(ctx, q.name, messages); err != nil {
		return fmt.Errorf("could not subscribe to queue: %w", err)
	}
	q.logger.Info("queue stopped", zap.Stringer("queue_name", q.name))

	return nil
}

func getLogger(logger *zap.Logger, event *WebhookEvent) *zap.Logger {
	if event.Installation != nil {
		return logger.With(
			zap.String("event_type", "installation"),
			zap.Int64("installation_id", event.Installation.GetInstallation().GetID()),
		)
	} else if event.InstallationRepositories != nil {
		return logger.With(
			zap.String("event_type", "installation repositories"),
			zap.Int64("installation_id", event.InstallationRepositories.GetInstallation().GetID()),
		)
	} else if event.Push != nil {
		return logger.With(
			zap.String("event_type", "push"),
			zap.Int64("installation_id", event.Push.GetInstallation().GetID()),
			zap.String("repo", event.Push.GetRepo().GetFullName()),
		)
	} else if event.PullRequest != nil {
		return logger.With(
			zap.String("event_type", "pull request"),
			zap.Int64("installation_id", event.PullRequest.GetInstallation().GetID()),
			zap.String("repo", event.PullRequest.GetRepo().GetFullName()),
		)
	} else if event.Status != nil {
		return logger.With(
			zap.String("event_type", "status"),
			zap.Int64("installation_id", event.Status.GetInstallation().GetID()),
			zap.String("repo", event.Status.GetRepo().GetFullName()),
		)
	} else if event.WorkflowJob != nil {
		return logger.With(
			zap.String("event_type", "workflow job"),
			zap.Int64("installation_id", event.WorkflowJob.GetInstallation().GetID()),
			zap.String("repo", event.WorkflowJob.GetRepo().GetFullName()),
		)
	} else {
		return logger.With(
			zap.String("event_type", "unknown"),
		)
	}
}

var errUnknownType = fmt.Errorf("unknown event type")

func (q *Queue) work(ctx context.Context, event *WebhookEvent) error {
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
