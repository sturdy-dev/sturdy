package routes

import (
	"fmt"
	"net/http"

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

		event, err := gh.ParseWebHook(gh.WebHookType(c.Request), payload)
		if err != nil {
			logger.Warn("failed to parse webhook", zap.Error(err))
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		logger = logger.With(
			zap.String("github_delivery", c.Request.Header.Get("X-GitHub-Delivery")),
			zap.String("github_event", c.Request.Header.Get("X-GitHub-Event")),
			zap.String("hook_id", c.Request.Header.Get("X-GitHub-Hook-Id")),
		)

		logger.Info("github webhook", zap.String("type", fmt.Sprintf("%T", event)))

		switch event := event.(type) {
		case *gh.InstallationEvent:
			if err := queue.Enqueue(c.Request.Context(), &workers_github.WebhookEvent{
				Installation: event,
			}); err != nil {
				logger.Error("failed to enqueue webhook", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		case *gh.InstallationRepositoriesEvent:
			if err := queue.Enqueue(c.Request.Context(), &workers_github.WebhookEvent{
				InstallationRepositories: event,
			}); err != nil {
				logger.Error("failed to enqueue webhook", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		case *gh.PushEvent:
			if err := queue.Enqueue(c.Request.Context(), &workers_github.WebhookEvent{
				Push: event,
			}); err != nil {
				logger.Error("failed to enqueue webhook", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		case *gh.PullRequestEvent:
			if err := queue.Enqueue(c.Request.Context(), &workers_github.WebhookEvent{
				PullRequest: event,
			}); err != nil {
				logger.Error("failed to enqueue webhook", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		case *gh.StatusEvent:
			if err := queue.Enqueue(c.Request.Context(), &workers_github.WebhookEvent{
				Status: event,
			}); err != nil {
				logger.Error("failed to enqueue webhook", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		case *gh.WorkflowJobEvent:
			if err := queue.Enqueue(c.Request.Context(), &workers_github.WebhookEvent{
				WorkflowJob: event,
			}); err != nil {
				logger.Error("failed to enqueue webhook", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		default:
			c.Status(http.StatusNotFound)
			return
		}

		c.Status(http.StatusOK)
	}
}
