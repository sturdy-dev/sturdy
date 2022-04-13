package cloud

import (
	routes_v3_analytics "getsturdy.com/api/pkg/analytics/enterprise/cloud/routes"
	authz "getsturdy.com/api/pkg/auth"
	"getsturdy.com/api/pkg/http/handler/enterprise/selfhosted"
	routes_v3_statistics "getsturdy.com/api/pkg/installations/statistics/enterprise/cloud/routes"
	service_statistics "getsturdy.com/api/pkg/installations/statistics/enterprise/cloud/service"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	routes_v3_licenses "getsturdy.com/api/pkg/licenses/enterprise/cloud/routes"
	service_licenses "getsturdy.com/api/pkg/licenses/enterprise/cloud/service"
	service_validations "getsturdy.com/api/pkg/licenses/enterprise/cloud/validations/service"
	routes_v3_logger "getsturdy.com/api/pkg/logger/enterprise/cloud/routes"
	routes_v3_user "getsturdy.com/api/pkg/users/enterprise/cloud/routes"
	service_user "getsturdy.com/api/pkg/users/enterprise/cloud/service"

	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
	"github.com/posthog/posthog-go"
	"go.uber.org/zap"
)

func ProvideHandler(
	logger *zap.Logger,
	posthogClient posthog.Client,
	enterpriseEngine *selfhosted.Engine,
	serviceLicenses *service_licenses.Service,
	serviceValidations *service_validations.Service,
	serviceStatistics *service_statistics.Service,
	sentryClient *raven.Client,
	jwtService *service_jwt.Service,
	userService *service_user.Service,
) *gin.Engine {
	auth := enterpriseEngine.Group("")
	auth.Use(authz.GinMiddleware(logger, jwtService))
	auth.POST("/v3/users/verify-email", routes_v3_user.SendEmailVerification(logger, userService)) // Used by the web (2021-11-14)

	publ := enterpriseEngine.Group("")
	publ.POST("/v3/analytics/batch/", routes_v3_analytics.Batch(logger, posthogClient))
	publ.GET("/v3/licenses/:key", routes_v3_licenses.Validate(logger, serviceLicenses, serviceValidations))
	publ.POST("/v3/statistics", gin.WrapF(routes_v3_statistics.Create(logger, serviceStatistics)))
	publ.POST("v3/sentry/store/", gin.WrapF(routes_v3_logger.Store(logger, sentryClient)))
	publ.POST("/v3/auth/magic-link/send", routes_v3_user.SendMagicLink(logger, userService))
	publ.POST("/v3/auth/magic-link/verify", routes_v3_user.VerifyMagicLink(logger, userService, jwtService))
	return (*gin.Engine)(enterpriseEngine)
}
