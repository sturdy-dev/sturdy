package enterprise

import (
	"getsturdy.com/api/pkg/analytics"
	authz "getsturdy.com/api/pkg/auth"
	db_change "getsturdy.com/api/pkg/change/db"
	service_ci "getsturdy.com/api/pkg/ci/service"
	workers_ci "getsturdy.com/api/pkg/ci/workers"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	service_comments "getsturdy.com/api/pkg/comments/service"
	"getsturdy.com/api/pkg/github/config"
	ghappclient "getsturdy.com/api/pkg/github/enterprise/client"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	routes_v3_ghapp "getsturdy.com/api/pkg/github/enterprise/routes"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	workers_github "getsturdy.com/api/pkg/github/enterprise/workers"
	"getsturdy.com/api/pkg/http/oss"
	service_buildkite "getsturdy.com/api/pkg/integrations/buildkite/enterprise/service"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	db_review "getsturdy.com/api/pkg/review/db"
	service_servicetokens "getsturdy.com/api/pkg/servicetokens/service"
	routes_ci "getsturdy.com/api/pkg/statuses/enterprise/routes"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	service_sync "getsturdy.com/api/pkg/sync/service"
	db_user "getsturdy.com/api/pkg/user/db"
	"getsturdy.com/api/pkg/view/events"
	activity_sender "getsturdy.com/api/pkg/workspace/activity/sender"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	service_workspace "getsturdy.com/api/pkg/workspace/service"
	"getsturdy.com/api/vcs/executor"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DevelopmentAllowExtraCorsOrigin string

type Engine gin.Engine

func ProvideHandler(
	logger *zap.Logger,
	userRepo db_user.Repository,
	analyticsClient analytics.Client,
	codebaseRepo db_codebase.CodebaseRepository,
	codebaseUserRepo db_codebase.CodebaseUserRepository,
	workspaceReader db_workspace.WorkspaceReader,
	changeRepo db_change.Repository,
	changeCommitRepo db_change.CommitRepository,
	gitHubInstallationsRepo db_github.GitHubInstallationRepo,
	gitHubRepositoryRepo db_github.GitHubRepositoryRepo,
	gitHubUserRepo db_github.GitHubUserRepo,
	gitHubPRRepo db_github.GitHubPRRepo,
	gitHubAppConfig config.GitHubAppConfig,
	gitHubClientProvider ghappclient.ClientProvider,
	gitHubClonerPublisher *workers_github.ClonerQueue,
	workspaceWriter db_workspace.WorkspaceWriter,
	executorProvider executor.Provider,
	reviewRepo db_review.ReviewRepository,
	activitySender activity_sender.ActivitySender,
	eventSender events.EventSender,
	workspaceService service_workspace.Service,
	statusesService *service_statuses.Service,
	syncService *service_sync.Service,
	jwtService *service_jwt.Service,
	commentsService *service_comments.Service,
	gitHubService *service_github.Service,
	ciBuildQueue *workers_ci.BuildQueue,
	ciService *service_ci.Service,
	serviceTokensService *service_servicetokens.Service,
	buildkiteService *service_buildkite.Service,
	ossEngine *oss.Engine,
) *Engine {
	auth := ossEngine.Group("")
	auth.Use(authz.GinMiddleware(logger, jwtService))
	auth.POST("/v3/github/oauth", routes_v3_ghapp.Oauth(logger, gitHubAppConfig, userRepo, gitHubUserRepo, gitHubService))

	publ := ossEngine.Group("")
	publ.POST("/v3/github/webhook", routes_v3_ghapp.Webhook(logger, gitHubAppConfig, analyticsClient, gitHubInstallationsRepo, gitHubRepositoryRepo, codebaseRepo, executorProvider, gitHubClientProvider, gitHubUserRepo, codebaseUserRepo, gitHubClonerPublisher, gitHubPRRepo, workspaceReader, workspaceWriter, workspaceService, syncService, changeRepo, changeCommitRepo, reviewRepo, eventSender, activitySender, statusesService, commentsService, gitHubService, ciBuildQueue))
	publ.POST("/v3/statuses/webhook", routes_ci.WebhookHandler(logger, statusesService, ciService, serviceTokensService, buildkiteService))
	return (*Engine)(ossEngine)
}
