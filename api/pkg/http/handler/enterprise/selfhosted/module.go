package selfhosted

import (
	"getsturdy.com/api/pkg/di"
	db_github "getsturdy.com/api/pkg/github/enterprise/db"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	webhooks_github "getsturdy.com/api/pkg/github/enterprise/webhooks"
	"getsturdy.com/api/pkg/http/handler"
	service_buildkite "getsturdy.com/api/pkg/integrations/providers/buildkite/enterprise/service"
	service_jwt "getsturdy.com/api/pkg/jwt/service"
	"getsturdy.com/api/pkg/logger"
	routes_remote "getsturdy.com/api/pkg/remote/enterprise/routes"
	service_servicetokens "getsturdy.com/api/pkg/servicetokens/service"
	service_statuses "getsturdy.com/api/pkg/statuses/service"
	db_user "getsturdy.com/api/pkg/users/db"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(db_user.Module)
	c.Import(db_github.Module)
	c.Import(service_statuses.Module)
	c.Import(service_jwt.Module)
	c.Import(service_github.Module)
	c.Import(service_servicetokens.Module)
	c.Import(service_buildkite.Module)
	c.Import(handler.Module)
	c.Import(webhooks_github.Module)
	c.Import(routes_remote.Module)
	c.Register(ProvideHandler)
}
