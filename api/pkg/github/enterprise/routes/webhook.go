package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	gh "github.com/google/go-github/v39/github"
	"go.uber.org/zap"

	service_github_webhooks "getsturdy.com/api/pkg/github/enterprise/webhooks"
)

func Webhook(logger *zap.Logger, gitHubWebhooksService *service_github_webhooks.Service) func(c *gin.Context) {
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

		logger.Info("github webhook", zap.String("type", fmt.Sprintf("%T", event)))

		switch event := event.(type) {
		case *gh.InstallationEvent:
			if err := gitHubWebhooksService.HandleInstallationEvent(event); err != nil {
				logger.Error("failed to handle github installation webhook event", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

		case *gh.InstallationRepositoriesEvent:
			if err := gitHubWebhooksService.HandleInstallationRepositoriesEvent(event); err != nil {
				logger.Error("failed to handle github installation repository webhook event", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

		case *gh.PushEvent:
			logger := logger.Named("githubHandlePushEvent").With(
				zap.String("repo", event.GetRepo().GetFullName()),
				zap.Int64("installation_id", event.GetInstallation().GetID()),
			)

			logger.Info("about to handle push event")

			if err := gitHubWebhooksService.HandlePushEvent(event); err != nil {
				logger.Error("failed to handle github push event", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			logger.Info("successfully handled push event")

		case *gh.PullRequestEvent:
			logger := logger.Named("githubHandlePullRequestEvent").With(
				zap.String("repo", event.GetRepo().GetFullName()),
				zap.Int64("installation_id", event.GetInstallation().GetID()),
			)

			if err := gitHubWebhooksService.HandlePullRequestEvent(event); err != nil {
				logger.Error("failed to handle github pull request event", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

		case *gh.StatusEvent:
			logger := logger.Named("githubHandleStatusEvent").With(
				zap.String("repo", event.GetRepo().GetFullName()),
				zap.Int64("installation_id", event.GetInstallation().GetID()),
			)

			if err := gitHubWebhooksService.HandleStatusEvent(event); err != nil {
				logger.Error("failed to handle status event", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		case *gh.WorkflowJobEvent:
			logger := logger.Named("githubHandleWorkflowJobEvent").With(
				zap.String("repo", event.GetRepo().GetFullName()),
				zap.Int64("installation_id", event.GetInstallation().GetID()),
			)

			if err := gitHubWebhooksService.HandleWorkflowJobEvent(event); err != nil {
				logger.Error("failed to handle workflow job event", zap.Error(err))
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
