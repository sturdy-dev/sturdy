package routes

import (
	"fmt"
	"net/http"

	"getsturdy.com/api/pkg/github/enterprise/webhooks"
	workers_github "getsturdy.com/api/pkg/github/enterprise/workers"

	"github.com/gin-gonic/gin"
	gh "github.com/google/go-github/v39/github"
	"go.uber.org/zap"
)

func Webhook(logger *zap.Logger, queue *workers_github.WebhooksQueue) func(c *gin.Context) {
	return func(c *gin.Context) {
		payload, err := gh.ValidatePayload(c.Request, nil)
		if err != nil {
			logger.Warn("failed to validate payload", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		event, err := webhooks.ParseWebHook(gh.WebHookType(c.Request), payload)
		if err != nil {
			logger.Warn("failed to parse webhook", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		logger := logger.With(
			zap.String("X-GitHub-Delivery", c.Request.Header.Get("X-GitHub-Delivery")),
			zap.String("X-GitHub-Event", c.Request.Header.Get("X-GitHub-Event")),
			zap.String("X-GitHub-Hook-Id", c.Request.Header.Get("X-GitHub-Hook-Id")),
			zap.Int64("Content-Length", c.Request.ContentLength),
		)

		logger.Info("github webhook", zap.String("type", fmt.Sprintf("%T", event)))

		switch event := event.(type) {
		case *webhooks.InstallationEvent:
			if err := queue.Enqueue(c.Request.Context(), &workers_github.WebhookEvent{
				Installation: event,
			}); err != nil {
				logger.Error("failed to enqueue webhook", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		case *webhooks.InstallationRepositoriesEvent:
			if err := queue.Enqueue(c.Request.Context(), &workers_github.WebhookEvent{
				InstallationRepositories: event,
			}); err != nil {
				logger.Error("failed to enqueue webhook", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		case *webhooks.PushEvent:
			if err := queue.Enqueue(c.Request.Context(), &workers_github.WebhookEvent{
				Push: event,
			}); err != nil {
				logger.Error("failed to enqueue webhook", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		case *webhooks.PullRequestEvent:
			if err := queue.Enqueue(c.Request.Context(), &workers_github.WebhookEvent{
				PullRequest: event,
			}); err != nil {
				logger.Error("failed to enqueue webhook", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		case *webhooks.StatusEvent:
			if err := queue.Enqueue(c.Request.Context(), &workers_github.WebhookEvent{
				Status: event,
			}); err != nil {
				logger.Error("failed to enqueue webhook", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		case *webhooks.WorkflowJobEvent:
			if err := queue.Enqueue(c.Request.Context(), &workers_github.WebhookEvent{
				WorkflowJob: event,
			}); err != nil {
				logger.Error("failed to enqueue webhook", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		default:
			logger.Warn("unsupported webhook type")
			c.Status(http.StatusNotFound)
			return
		}

		c.Status(http.StatusOK)
	}
}
