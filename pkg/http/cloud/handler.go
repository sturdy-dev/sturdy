package cloud

import (
	"mash/pkg/analytics"
	"mash/pkg/analytics/proxy/routes"
	"mash/pkg/http/enterprise"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ProvideHandler(
	logger *zap.Logger,
	analyticsClient analytics.Client,
	enterpriseEngine *enterprise.Engine,
) *gin.Engine {
	publ := enterpriseEngine.Group("")
	publ.POST("/v3/analytics/batch/", routes.Batch(logger, analyticsClient))
	return (*gin.Engine)(enterpriseEngine)
}
