package pr

import (
	service_auth "getsturdy.com/api/pkg/auth/service"
	db_codebases "getsturdy.com/api/pkg/codebases/db"
	configuration "getsturdy.com/api/pkg/configuration/module"
	"getsturdy.com/api/pkg/di"
	eventsv2 "getsturdy.com/api/pkg/events/v2"
	"getsturdy.com/api/pkg/github/enterprise/client"
	"getsturdy.com/api/pkg/github/enterprise/db"
	service_github "getsturdy.com/api/pkg/github/enterprise/service"
	"getsturdy.com/api/pkg/graphql/resolvers"
	"getsturdy.com/api/pkg/logger"
	db_user "getsturdy.com/api/pkg/users/db"
	db_view "getsturdy.com/api/pkg/view/db"
	db_workspaces "getsturdy.com/api/pkg/workspaces/db"
)

func Module(c *di.Container) {
	c.Import(logger.Module)
	c.Import(resolvers.Module)
	c.Import(db_user.Module)
	c.Import(db_codebases.Module)
	c.Import(db_workspaces.Module)
	c.Import(db_view.Module)
	c.Import(configuration.Module)
	c.Import(db.Module)
	c.Import(service_auth.Module)
	c.Import(service_github.Module)
	c.Import(client.Module)
	c.Import(eventsv2.Module)
	c.Register(NewResolver)
}
