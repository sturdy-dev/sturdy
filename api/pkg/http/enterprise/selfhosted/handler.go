package selfhosted

import (
	service_analytics "getsturdy.com/api/pkg/analytics/service"
	authz "getsturdy.com/api/pkg/auth"
	db_change "getsturdy.com/api/pkg/change/db"
	service_ci "getsturdy.com/api/pkg/ci/service"
	db_codebase "getsturdy.com/api/pkg/codebase/db"
	"getsturdy.com/api/pkg/events"
	ghappclient "getsturdy.com/api/pkg/github/enterprise/client"
	"getsturdy.com/api/pkg/github/enterprise/config"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	routes_v3_ghapp "getsturdy.com/api/pkg/github/enterprise/routes"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	service_github_webhooks "getsturdy.com/api/pkg/github/enterprise/webhooks"
	"getsturdy.com/api/pkg/http"
	service_buildkite "getsturdy.com/api/pkg/integrations/buildkite/enterprise/service"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	service_servicetokens "getsturdy.com/api/pkg/servicetokens/service"
	routes_ci "getsturdy.com/api/pkg/statuses/enterprise/routes"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	db_user "getsturdy.com/api/pkg/users/db"
	db_workspace "getsturdy.com/api/pkg/workspace/db"
	"getsturdy.com/api/vcs/executor"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DevelopmentAllowExtraCorsOrigin string

type Engine gin.Engine

func ProvideHandler(
	logger *zap.Logger,
	userRepo db_user.Repository,
	analyticsSerivce *service_analytics.Service,
	codebaseRepo db_codebase.CodebaseRepository,
	workspaceReader db_workspace.WorkspaceReader,
	changeRepo db_change.Repository,
	gitHubInstallationRepo db_github.GitHubInstallationRepo,
	gitHubRepositoryRepo db_github.GitHubRepositoryRepo,
	gitHubUserRepo db_github.GitHubUserRepo,
	gitHubPRRepo db_github.GitHubPRRepo,
	gitHubAppConfig *config.GitHubAppConfig,
	githubClientProvider ghappclient.InstallationClientProvider,
	workspaceWriter db_workspace.WorkspaceWriter,
	executorProvider executor.Provider,
	eventsSender events.EventSender,
	statusesService *service_statuses.Service,
	jwtService *service_jwt.Service,
	gitHubService *service_github.Service,
	ciService *service_ci.Service,
	serviceTokensService *service_servicetokens.Service,
	buildkiteService *service_buildkite.Service,
	ossEngine *http.Engine,
	gitHubWebhooksService *service_github_webhooks.Service,
) *Engine {
	auth := ossEngine.Group("")
	auth.Use(authz.GinMiddleware(logger, jwtService))
	auth.POST("/v3/github/oauth", routes_v3_ghapp.Oauth(logger, gitHubAppConfig, userRepo, gitHubUserRepo, gitHubService))

	publ := ossEngine.Group("")
	publ.POST("/v3/github/webhook", routes_v3_ghapp.Webhook(logger, analyticsSerivce, gitHubInstallationRepo, gitHubRepositoryRepo, codebaseRepo, statusesService, gitHubService, gitHubWebhooksService))
	publ.POST("/v3/statuses/webhook", routes_ci.WebhookHandler(logger, statusesService, ciService, serviceTokensService, buildkiteService))
	return (*Engine)(ossEngine)
}
