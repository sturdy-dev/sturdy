package routes

import (
	"fmt"
	"net/http"

	db_change "mash/pkg/change/db"
	workers_ci "mash/pkg/ci/workers"
	db_codebase "mash/pkg/codebase/db"
	service_comments "mash/pkg/comments/service"
	"mash/pkg/github/config"
	"mash/pkg/github/enterprise/client"
	"mash/pkg/github/enterprise/db"
	"mash/pkg/github/enterprise/routes/installation"
	"mash/pkg/github/enterprise/routes/pr"
	"mash/pkg/github/enterprise/routes/push"
	"mash/pkg/github/enterprise/routes/statuses"
	"mash/pkg/github/enterprise/routes/workflows"
	service_github "mash/pkg/github/enterprise/service"
	workers_github "mash/pkg/github/enterprise/workers"
	db_review "mash/pkg/review/db"
	service_statuses "mash/pkg/statuses/service"
	service_sync "mash/pkg/sync/service"
	"mash/pkg/view/events"
	activity_sender "mash/pkg/workspace/activity/sender"
	db_workspace "mash/pkg/workspace/db"
	service_workspace "mash/pkg/workspace/service"
	"mash/vcs/executor"

	"github.com/gin-gonic/gin"
	gh "github.com/google/go-github/v39/github"
	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

func Webhook(
	logger *zap.Logger,
	config config.GitHubAppConfig,
	postHogClient posthog.Client,
	gitHubInstallationRepo db.GitHubInstallationRepo,
	gitHubRepositoryRepo db.GitHubRepositoryRepo,
	codebaseRepo db_codebase.CodebaseRepository,
	executorProvider executor.Provider,
	githubClientProvider client.ClientProvider,
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
			if err := installation.HandleInstallationEvent(c, logger.Named("githubHandleInstallationEvent"), event, gitHubInstallationRepo, gitHubRepositoryRepo, postHogClient, codebaseRepo, gitHubService); err != nil {
				logger.Error("failed to handle github installation webhook event", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

		case *gh.InstallationRepositoriesEvent:
			if err := installation.HandleInstallationRepositoriesEvent(c, logger.Named("githubHandleInstallationRepositoriesEvent"), event, gitHubInstallationRepo, gitHubRepositoryRepo, postHogClient, codebaseRepo, gitHubService); err != nil {
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

			if err := push.HandlePushEvent(c, logger, event, gitHubRepositoryRepo, gitHubInstallationRepo, workspaceWriter, workspaceReader, workspaceService, syncService, gitHubPRRepo, changeRepo, changeCommitRepo, executorProvider, config, githubClientProvider, eventsSender, postHogClient, reviewRepo, activitySender, commentsService, buildQueue); err != nil {
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
