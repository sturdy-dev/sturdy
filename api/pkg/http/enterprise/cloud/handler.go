package cloud

import (
	"getsturdy.com/api/pkg/analytics"
	routes_v3_analytics "getsturdy.com/api/pkg/analytics/enterprise/cloud/routes"
	"getsturdy.com/api/pkg/http/enterprise/selfhosted"
	routes_v3_statistics "getsturdy.com/api/pkg/installations/statistics/enterprise/cloud/routes"
	service_statistics "getsturdy.com/api/pkg/installations/statistics/enterprise/cloud/service"
	routes_v3_licenses "getsturdy.com/api/pkg/licenses/enterprise/cloud/routes"
	service_licenses "getsturdy.com/api/pkg/licenses/enterprise/cloud/service"
	service_validations "getsturdy.com/api/pkg/licenses/enterprise/cloud/validations/service"
	routes_v3_logger "getsturdy.com/api/pkg/logger/enterprise/cloud/routes"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ProvideHandler(
	logger *zap.Logger,
	analyticsClient analytics.Client,
	enterpriseEngine *selfhosted.Engine,
	serviceLicenses *service_licenses.Service,
	serviceValidations *service_validations.Service,
	serviceStatistics *service_statistics.Service,
	sentryClient *raven.Client,
) *gin.Engine {
	publ := enterpriseEngine.Group("")
	publ.POST("/v3/analytics/batch/", routes_v3_analytics.Batch(logger, analyticsClient))
	publ.GET("/v3/licenses/:key", routes_v3_licenses.Validate(logger, serviceLicenses, serviceValidations))
	publ.POST("/v3/statistics", gin.WrapF(routes_v3_statistics.Create(logger, serviceStatistics)))
	publ.POST("v3/sentry/store/", gin.WrapF(routes_v3_logger.Store(logger, sentryClient)))
	return (*gin.Engine)(enterpriseEngine)
}
