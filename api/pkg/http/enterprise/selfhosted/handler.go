package selfhosted

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	authz "getsturdy.com/api/pkg/auth"
	service_ci "getsturdy.com/api/pkg/ci/service"
	"getsturdy.com/api/pkg/github/enterprise/config"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	routes_v3_ghapp "getsturdy.com/api/pkg/github/enterprise/routes"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	workers_github "getsturdy.com/api/pkg/github/enterprise/workers"
	"getsturdy.com/api/pkg/http"
	service_buildkite "getsturdy.com/api/pkg/integrations/providers/buildkite/enterprise/service"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	service_servicetokens "getsturdy.com/api/pkg/servicetokens/service"
	routes_ci "getsturdy.com/api/pkg/statuses/enterprise/routes"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	db_user "getsturdy.com/api/pkg/users/db"
)

type Engine gin.Engine

func ProvideHandler(
	logger *zap.Logger,
	userRepo db_user.Repository,
	gitHubUserRepo db_github.GitHubUserRepo,
	gitHubAppConfig *config.GitHubAppConfig,
	statusesService *service_statuses.Service,
	jwtService *service_jwt.Service,
	gitHubService *service_github.Service,
	ciService *service_ci.Service,
	serviceTokensService *service_servicetokens.Service,
	buildkiteService *service_buildkite.Service,
	ossEngine *http.Engine,
	gitHubWebhooksQueue *workers_github.WebhooksQueue,
) *Engine {
	auth := ossEngine.Group("")
	auth.Use(authz.GinMiddleware(logger, jwtService))
	auth.POST("/v3/github/oauth", routes_v3_ghapp.Oauth(logger, gitHubAppConfig, userRepo, gitHubUserRepo, gitHubService))

	publ := ossEngine.Group("")
	publ.POST("/v3/github/webhook", routes_v3_ghapp.Webhook(logger, gitHubWebhooksQueue))
	publ.POST("/v3/statuses/webhook", routes_ci.WebhookHandler(logger, statusesService, ciService, serviceTokensService, buildkiteService))
	return (*Engine)(ossEngine)
}
