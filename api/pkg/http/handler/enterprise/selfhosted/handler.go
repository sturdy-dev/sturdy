package selfhosted

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	authz "getsturdy.com/api/pkg/auth"
	service_buildkite_enterprise "getsturdy.com/api/pkg/buildkite/enterprise/service"
	service_ci "getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/github/enterprise/config"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	routes_v3_ghapp "getsturdy.com/api/pkg/github/enterprise/routes"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	webhooks_github "getsturdy.com/api/pkg/github/enterprise/webhooks"
	"getsturdy.com/api/pkg/http/handler"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	routes_remote "getsturdy.com/api/pkg/remote/enterprise/routes"
	service_servicetokens "getsturdy.com/api/pkg/servicetokens/service"
	routes_ci "getsturdy.com/api/pkg/statuses/enterprise/routes"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	db_user "getsturdy.com/api/pkg/users/db"
)

type Engine gin.Engine

func ProvideHandler(
	logger *zap.Logger,
	userRepo db_user.Repository,
	gitHubUserRepo db_github.GitHubUserRepository,
	gitHubAppConfig *config.GitHubAppConfig,
	statusesService *service_statuses.Service,
	jwtService *service_jwt.Service,
	gitHubService *service_github.Service,
	ciService *service_ci.Service,
	serviceTokensService *service_servicetokens.Service,
	enterpriseBuildkiteService *service_buildkite_enterprise.Service,
	ossEngine *handler.Engine,
	gitHubWebhooksQueue *webhooks_github.Queue,
	triggerSyncCodebaseWebhookHandler routes_remote.TriggerSyncCodebaseWebhookHandler,
) *Engine {
	auth := ossEngine.Group("")
	auth.Use(authz.GinMiddleware(logger, jwtService))
	auth.POST("/v3/github/oauth", routes_v3_ghapp.Oauth(logger, gitHubAppConfig, userRepo, gitHubUserRepo, gitHubService))

	publ := ossEngine.Group("")
	publ.POST("/v3/github/webhook", routes_v3_ghapp.Webhook(logger, gitHubWebhooksQueue))
	publ.POST("/v3/statuses/webhook", routes_ci.WebhookHandler(logger, statusesService, ciService, serviceTokensService, enterpriseBuildkiteService))

	// Using Any to give friendly error messages if sent a non-POST request
	publ.Any("/v3/remotes/webhook/sync-codebase/:id", gin.HandlerFunc(triggerSyncCodebaseWebhookHandler))
	return (*Engine)(ossEngine)
}
