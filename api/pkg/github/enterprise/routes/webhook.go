package routes

import (
	"fmt"
	"net/http"

	"getsturdy.com/api/pkg/analytics"
	db_change "getsturdy.com/api/pkg/change/db"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/events"
	"getsturdy.com/api/pkg/github/enterprise/client"
	"getsturdy.com/api/pkg/github/enterprise/config"
	"getsturdy.com/api/pkg/github/enterprise/db"
	"getsturdy.com/api/pkg/github/enterprise/routes/installation"
	"getsturdy.com/api/pkg/github/enterprise/routes/pr"
	"getsturdy.com/api/pkg/github/enterprise/routes/push"
	"getsturdy.com/api/pkg/github/enterprise/routes/statuses"
	"getsturdy.com/api/pkg/github/enterprise/routes/workflows"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	workers_github "getsturdy.com/api/pkg/github/enterprise/workers"
	db_review "getsturdy.com/api/pkg/review/db"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	service_sync "getsturdy.com/api/pkg/sync/service"
	activity_sender "getsturdy.com/api/pkg/workspace/activity/sender"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	service_workspace "getsturdy.com/api/pkg/workspace/service"
	"getsturdy.com/api/vcs/executor"

	"github.com/gin-gonic/gin"
	gh "github.com/google/go-github/v39/github"
	"go.uber.org/zap"
)

func Webhook(
	logger *zap.Logger,
	config *config.GitHubAppConfig,
	analyticsClient analytics.Client,
	gitHubInstallationRepo db.GitHubInstallationRepo,
	gitHubRepositoryRepo db.GitHubRepositoryRepo,
	codebaseRepo db_codebase.CodebaseRepository,
	executorProvider executor.Provider,
	githubClientProvider client.InstallationClientProvider,
	gitHubUserRepo db.GitHubUserRepo,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	gitHubClonerPublisher *workers_github.ClonerQueue,
	gitHubPRRepo db.GitHubPRRepo,
	workspaceReader db_workspace.WorkspaceReader,
	workspaceWriter db_workspace.WorkspaceWriter,
	workspaceService service_workspace.Service,
	syncService *service_sync.Service,
	changeRepo db_change.Repository,
	changeCommitRepo db_change.CommitRepository,
	reviewRepo db_review.ReviewRepository,
	eventsSender events.EventSender,
	activitySender activity_sender.ActivitySender,
	statusesService *service_statuses.Service,
	commentsService *service_comments.Service,
	gitHubService *service_github.Service,
	buildQueue *workers_ci.BuildQueue,
) func(c *gin.Context) {
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
			if err := installation.HandleInstallationEvent(c, logger.Named("githubHandleInstallationEvent"), event, gitHubInstallationRepo, gitHubRepositoryRepo, analyticsClient, codebaseRepo, gitHubService); err != nil {
				logger.Error("failed to handle github installation webhook event", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

		case *gh.InstallationRepositoriesEvent:
			if err := installation.HandleInstallationRepositoriesEvent(c, logger.Named("githubHandleInstallationRepositoriesEvent"), event, gitHubInstallationRepo, gitHubRepositoryRepo, analyticsClient, codebaseRepo, gitHubService); err != nil {
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

			if err := push.HandlePushEvent(c, logger, event, gitHubRepositoryRepo, gitHubInstallationRepo, workspaceWriter, workspaceReader, workspaceService, syncService, gitHubPRRepo, changeRepo, changeCommitRepo, executorProvider, config, githubClientProvider, eventsSender, analyticsClient, reviewRepo, activitySender, commentsService, buildQueue); err != nil {
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

			if err := pr.HandlePullRequestEvent(logger, event, workspaceReader, gitHubPRRepo, eventsSender, workspaceWriter); err != nil {
				logger.Error("failed to handle github pull request event", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

		case *gh.StatusEvent:
			logger := logger.Named("githubHandleStatusEvent").With(
				zap.String("repo", event.GetRepo().GetFullName()),
				zap.Int64("installation_id", event.GetInstallation().GetID()),
			)

			if err := statuses.HandleStatusEvent(
				c.Request.Context(),
				logger,
				event,
				gitHubRepositoryRepo,
				statusesService,
			); err != nil {
				logger.Error("failed to handle status event", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}
		case *gh.WorkflowJobEvent:
			logger := logger.Named("githubHandleWorkflowJobEvent").With(
				zap.String("repo", event.GetRepo().GetFullName()),
				zap.Int64("installation_id", event.GetInstallation().GetID()),
			)

			if err := workflows.HandleWorkflowJobEvent(
				c.Request.Context(),
				logger,
				event,
				gitHubRepositoryRepo,
				statusesService,
			); err != nil {
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
