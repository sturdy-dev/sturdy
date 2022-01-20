package routes

import (
	service_ci "getsturdy.com/api/pkg/ci/service"
	routes_buildkite "getsturdy.com/api/pkg/integrations/buildkite/enterprise/routes"
	service_buildkite "getsturdy.com/api/pkg/integrations/buildkite/enterprise/service"
	service_servicetokens "getsturdy.com/api/pkg/servicetokens/service"
	service_statuses "getsturdy.com/api/pkg/statuses/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func WebhookHandler(
	logger *zap.Logger,
	statusesService *service_statuses.Service,
	ciService *service_ci.Service,
	serviceTokensService *service_servicetokens.Service,
	buildkiteService *service_buildkite.Service,
) func(c *gin.Context) {
	isBuildkite := func(c *gin.Context) bool {
		return c.GetHeader("X-Buildkite-Event") != ""
	}
	return func(c *gin.Context) {
		switch {
		case isBuildkite(c):
			routes_buildkite.WebhookHandler(logger, statusesService, ciService, serviceTokensService, buildkiteService)(c)
		default:
			c.AbortWithStatus(404)
			return
		}
	}
}
