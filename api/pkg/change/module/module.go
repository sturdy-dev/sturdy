package module

import (
	"getsturdy.com/api/pkg/change/db"
	"getsturdy.com/api/pkg/change/decorate"
	module_downloads "getsturdy.com/api/pkg/change/downloads/module"
	"getsturdy.com/api/pkg/change/graphql"
	"getsturdy.com/api/pkg/change/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
	c.Import(module_downloads.Module)
	c.Import(decorate.Module)
}
