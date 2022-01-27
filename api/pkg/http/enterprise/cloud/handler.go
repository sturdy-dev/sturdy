package cloud

import (
	"getsturdy.com/api/pkg/analytics"
	"getsturdy.com/api/pkg/analytics/cloud/routes"
	"getsturdy.com/api/pkg/http/enterprise/selfhosted"
	routes_v3_licenses "getsturdy.com/api/pkg/licenses/enterprise/cloud/routes"
	service_licenses "getsturdy.com/api/pkg/licenses/enterprise/cloud/service"
	service_validations "getsturdy.com/api/pkg/licenses/enterprise/cloud/validations/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ProvideHandler(
	logger *zap.Logger,
	analyticsClient analytics.Client,
	enterpriseEngine *selfhosted.Engine,
	serviceLicenses *service_licenses.Service,
	serviceValidations *service_validations.Service,
) *gin.Engine {
	publ := enterpriseEngine.Group("")
	publ.POST("/v3/analytics/batch/", routes.Batch(logger, analyticsClient))
	publ.GET("/v3/licenses/:key", routes_v3_licenses.Validate(logger, serviceLicenses, serviceValidations))
	return (*gin.Engine)(enterpriseEngine)
}
