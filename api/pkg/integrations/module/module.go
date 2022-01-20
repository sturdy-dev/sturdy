package module

import (
	"getsturdy.com/api/pkg/di"
	module_buildkite "getsturdy.com/api/pkg/integrations/buildkite/module"
	"getsturdy.com/api/pkg/integrations/db"
	"getsturdy.com/api/pkg/integrations/graphql"
)

func Module(c *di.Container) {
	c.Import(db.Module)
	c.Import(graphql.Module)
	c.Import(module_buildkite.Module)
}
