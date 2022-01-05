package routes

import (
	service_ci "mash/pkg/ci/service"
	db_codebase "mash/pkg/codebase/db"
	routes_buildkite "mash/pkg/integrations/buildkite/routes"
	service_buildkite "mash/pkg/integrations/buildkite/service"
	service_servicetokens "mash/pkg/servicetokens/service"
	service_statuses "mash/pkg/statuses/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func WebhookHandler(
	logger *zap.Logger,
	codebaseRepo db_codebase.CodebaseRepository,
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
