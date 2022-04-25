package routes

import (
	routes_buildkite "getsturdy.com/api/pkg/buildkite/enterprise/routes"
	service_buildkite_enterprise "getsturdy.com/api/pkg/buildkite/enterprise/service"
	service_ci "getsturdy.com/api/pkg/ci/service"
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
	enterpriseBuildkiteService *service_buildkite_enterprise.Service,
) func(c *gin.Context) {
	isBuildkite := func(c *gin.Context) bool {
		return c.GetHeader("X-Buildkite-Event") != ""
	}
	return func(c *gin.Context) {
		switch {
		case isBuildkite(c):
			routes_buildkite.WebhookHandler(logger, statusesService, ciService, serviceTokensService, enterpriseBuildkiteService)(c)
		default:
			c.AbortWithStatus(404)
			return
		}
	}
}
