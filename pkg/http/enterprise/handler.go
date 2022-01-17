package enterprise

import (
	"net/http"

	"mash/pkg/analytics"
	authz "mash/pkg/auth"
	db_change "mash/pkg/change/db"
	service_ci "mash/pkg/ci/service"
	workers_ci "mash/pkg/ci/workers"
	db_codebase "mash/pkg/codebase/db"
	service_comments "mash/pkg/comments/service"
	"mash/pkg/github/config"
	ghappclient "mash/pkg/github/enterprise/client"
	db_github "mash/pkg/github/enterprise/db"
	routes_v3_ghapp "mash/pkg/github/enterprise/routes"
	service_github "mash/pkg/github/enterprise/service"
	workers_github "mash/pkg/github/enterprise/workers"
	service_buildkite "mash/pkg/integrations/buildkite/enterprise/service"
	service_jwt "mash/pkg/jwt/service"
	db_review "mash/pkg/review/db"
	service_servicetokens "mash/pkg/servicetokens/service"
	routes_ci "mash/pkg/statuses/enterprise/routes"
	service_statuses "mash/pkg/statuses/service"
	service_sync "mash/pkg/sync/service"
	db_user "mash/pkg/user/db"
	"mash/pkg/view/events"
	activity_sender "mash/pkg/workspace/activity/sender"
	db_workspace "mash/pkg/workspace/db"
	service_workspace "mash/pkg/workspace/service"
	"mash/vcs/executor"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DevelopmentAllowExtraCorsOrigin string

type Engine = gin.Engine

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
	ossEngine *gin.Engine,
) http.Handler {
	auth := ossEngine.Group("")
	auth.Use(authz.GinMiddleware(logger, jwtService))
	auth.POST("/v3/github/oauth", routes_v3_ghapp.Oauth(logger, gitHubAppConfig, userRepo, gitHubUserRepo, gitHubService))

	publ := ossEngine.Group("")
	publ.POST("/v3/github/webhook", routes_v3_ghapp.Webhook(logger, gitHubAppConfig, analyticsClient, gitHubInstallationsRepo, gitHubRepositoryRepo, codebaseRepo, executorProvider, gitHubClientProvider, gitHubUserRepo, codebaseUserRepo, gitHubClonerPublisher, gitHubPRRepo, workspaceReader, workspaceWriter, workspaceService, syncService, changeRepo, changeCommitRepo, reviewRepo, eventSender, activitySender, statusesService, commentsService, gitHubService, ciBuildQueue))
	publ.POST("/v3/statuses/webhook", routes_ci.WebhookHandler(logger, statusesService, ciService, serviceTokensService, buildkiteService))
	return ossEngine
}
