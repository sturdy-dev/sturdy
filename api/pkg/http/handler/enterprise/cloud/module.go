package cloud

import (
	"net/http"

	analytics "getsturdy.com/api/pkg/analytics/enterprise/cloud/posthog"
	"getsturdy.com/api/pkg/di"
	"getsturdy.com/api/pkg/http/handler/enterprise/selfhosted"
	service_statistics "getsturdy.com/api/pkg/installations/statistics/enterprise/cloud/service"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	service_licenses "getsturdy.com/api/pkg/licenses/enterprise/cloud/service"
	service_validations "getsturdy.com/api/pkg/licenses/enterprise/cloud/validations/service"
	"getsturdy.com/api/pkg/logger"
	service_user "getsturdy.com/api/pkg/users/service/module"
)

func Module(c *di.Container) {
	c.Import(selfhosted.Module)
	c.Import(logger.Module)
	c.Import(analytics.Module)
	c.Import(service_licenses.Module)
	c.Import(service_validations.Module)
	c.Import(service_statistics.Module)
	c.Import(service_jwt.Module)
	c.Import(service_user.Module)
	c.Register(ProvideHandler, new(http.Handler))
}
