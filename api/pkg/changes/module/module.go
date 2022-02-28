package module

import (
	"getsturdy.com/api/pkg/changes/db"
	module_downloads "getsturdy.com/api/pkg/changes/downloads/module"
	"getsturdy.com/api/pkg/changes/graphql"
	"getsturdy.com/api/pkg/changes/service"
	"getsturdy.com/api/pkg/di"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(service.Module)
	c.Import(graphql.Module)
	c.Import(module_downloads.Module)
}
